package servers

import (
	"backend-auth/internal/logger"
	"backend-auth/internal/utils"
	"backend-auth/pkg/database"
	"backend-auth/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *Server) CreateUser(ctx echo.Context) error {
	user := new(models.User)
	if err := ctx.Bind(user); err != nil {
		logger.LogFailure(err, "Error binding user for creation")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	hashesPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		logger.LogFailure(err, "Error hashing user's password")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	user.Password = string(hashesPassword)
	err = s.datasource.CreateUser(user)
	if errors.Is(err, &database.DuplicateEmailError{}) {
		return ctx.JSON(http.StatusConflict, echo.Map{})
	}
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}

	deviceID := uuid.New().String()
	accessToken, _, refreshToken, refreshTokenExpiry, err := generateAccessRefreshTokens(user.ID, deviceID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	err = s.createAccessRefreshTokens(user.ID, deviceID, accessToken, refreshToken, refreshTokenExpiry)
	if err != nil {
		logger.LogFailure(err, "Error saving access/refresh tokens")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}

	return ctx.JSON(http.StatusCreated, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *Server) GetUsersCount(ctx echo.Context) error {
	count := s.datasource.GetUsersCount()
	return ctx.String(http.StatusOK, strconv.FormatInt(count, 10))
}

func (s *Server) RefreshTokens(ctx echo.Context) error {
	userID, deviceID, err := utils.FetchUserAndDevice(ctx)
	if err != nil {
		logger.LogFailure(err, "Error fetching the user or device ID from refresh token")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	oldRefreshToken, oldTokenExpiry, err := utils.FetchRefreshTokenAndExpiry(ctx)
	if err != nil {
		logger.LogFailure(err, "Error fetching the old refresh token or old refresh token expiry time")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	usedBefore, _ := s.cache.Connection.IsUsedRefreshToken(oldRefreshToken)
	if !s.cache.Enabled {
		usedBefore, err = s.datasource.IsUsedRefreshToken(oldRefreshToken)
		if err != nil {
			logger.LogFailure(err, fmt.Sprintf("Error checking used refresh token from DB: %s", oldRefreshToken))
			return ctx.JSON(http.StatusInternalServerError, echo.Map{})
		}
	}
	if usedBefore {
		queueName := "auth::invalidate-refresh-token-family"
		err = s.messaging.Connection.SendMessage(queueName, map[string]interface{}{
			"refresh_token": oldRefreshToken,
		})
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{})
		}
		return ctx.JSON(http.StatusBadRequest, echo.Map{})
	}
	compromised, _ := s.cache.Connection.IsCompromisedRefreshToken(oldRefreshToken)
	if compromised {
		return ctx.JSON(http.StatusBadRequest, echo.Map{})
	}
	accessToken, _, refreshToken, refreshTokenExpiry, err := generateAccessRefreshTokens(userID, deviceID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	userTokens := models.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	usedRefreshToken := models.UsedRefreshToken{
		UserID:             userID,
		RefreshToken:       oldRefreshToken,
		RefreshTokenExpiry: oldTokenExpiry,
	}
	generatedRefreshToken := models.GeneratedRefreshToken{
		UserID:               userID,
		RefreshToken:         refreshToken,
		RefreshTokenExpiry:   time.Unix(refreshTokenExpiry, 0),
		ParentRefreshTokenID: sql.NullInt64{Int64: int64(s.datasource.GetGeneratedRefreshToken(oldRefreshToken).ID), Valid: true},
	}
	err = s.datasource.InitializeTransaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("user_id = ? AND device_id = ?", userID, deviceID).
			Updates(&userTokens).Error; err != nil {
			logger.LogFailure(err, "Error creating/updating the user tokens")
			return err // return any error will rollback
		}
		if err := tx.Create(&usedRefreshToken).Error; err != nil {
			logger.LogFailure(err, "Error adding refresh token to used pool")
			return err
		}
		if err := tx.Create(&generatedRefreshToken).Error; err != nil {
			logger.LogFailure(err, "Error creating generated refresh token")
			return err
		}
		if count, err := s.cache.Connection.MarkRefreshTokenAsUsed(oldRefreshToken); err != nil {
			logger.LogFailure(err, "Error adding the used refresh token to Redis")
			return err
		} else if count < 1 {
			err := errors.New("cannot mark refresh token as used (it may have been used before)")
			logger.LogFailure(err, "Error marking the refresh token in redis")
			return err
		}

		return nil // return nil will commit the whole transaction
	})

	if err := s.cache.Connection.SaveAccessRefreshTokens(userID, deviceID, accessToken, refreshToken); err != nil {
		logger.LogFailure(err, "Error adding the refreshed access/refresh token to Redis")
	}
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	return ctx.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func generateAccessRefreshTokens(userID uint, deviceID string) (string, int64, string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["device_id"] = deviceID
	accessTokenExpiry := time.Now().Add(time.Hour * 24).Unix()       // access token to expire in 1 day
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 90).Unix() // refresh token to expire in 90 days
	claims["exp"] = accessTokenExpiry
	claims["r"] = uuid.New().String() // randomize jwt even if two requests came in the same second (useful in testing)
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN")))
	if err != nil {
		logger.LogFailure(err, "Error generating the access token")
		return "", 0, "", 0, err
	}
	claims["exp"] = refreshTokenExpiry
	refreshToken, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN")))
	if err != nil {
		logger.LogFailure(err, "Error generating the refresh token")
		return "", 0, "", 0, err
	}
	return accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, nil
}

func (s *Server) createAccessRefreshTokens(userID uint, deviceID, accessToken, refreshToken string, refreshTokenExpiry int64) error {
	userTokens := &models.UserTokens{
		UserID:       userID,
		DeviceID:     deviceID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	generatedRefreshToken := &models.GeneratedRefreshToken{
		UserID:               userID,
		RefreshToken:         refreshToken,
		RefreshTokenExpiry:   time.Unix(refreshTokenExpiry, 0),
		ParentRefreshTokenID: sql.NullInt64{Int64: 0, Valid: false},
	}
	err := s.datasource.InitializeTransaction(func(tx *gorm.DB) error {
		if err := tx.Create(userTokens).Error; err != nil {
			return err
		}
		if err := tx.Create(generatedRefreshToken).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	_ = s.cache.Connection.SaveAccessRefreshTokens(userID, deviceID, accessToken, refreshToken)
	return nil
}
