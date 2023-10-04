package controllers

import (
	"backend-auth/database"
	"backend-auth/logger"
	"backend-auth/models"
	"backend-auth/utils"
	"errors"
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

func (c *Controller) CreateUser(ctx echo.Context) error {
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
	err = c.datasource.CreateUser(user)
	if errors.Is(err, &database.DuplicateEmailError{}) {
		return ctx.JSON(http.StatusConflict, echo.Map{})
	}
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}

	deviceID := uuid.New().String()
	accessToken, refreshToken, err := generateAccessRefreshTokens(user.ID, deviceID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	err = c.updateAccessRefreshTokens(user.ID, deviceID, accessToken, refreshToken)
	if err != nil {
		logger.LogFailure(err, "Error saving access/refresh tokens")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}

	return ctx.JSON(http.StatusCreated, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *Controller) GetUsersCount(ctx echo.Context) error {
	count := c.datasource.GetUsersCount()
	return ctx.String(http.StatusOK, strconv.FormatInt(count, 10))
}

func (c *Controller) RefreshTokens(ctx echo.Context) error {
	userID, deviceID, err := utils.FetchUserAndDevice(ctx)
	if err != nil {
		logger.LogFailure(err, "Error fetching the user or device ID from refresh token")
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	accessToken, refreshToken, err := generateAccessRefreshTokens(userID, deviceID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}
	oldRefreshToken, oldTokenExpiry, err := utils.FetchRefreshTokenAndExpiry(ctx)
	if err != nil {
		logger.LogFailure(err, "Error fetching the old refresh token or old refresh token expiry time")
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
	err = c.datasource.InitializeTransaction(func(tx *gorm.DB) error {
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

		if _, err := c.cache.Connection.MarkRefreshTokenAsUsed(oldRefreshToken); err != nil {
			logger.LogFailure(err, "Error adding the used refresh token to Redis")
			return err
		}

		return nil // return nil will commit the whole transaction
	})

	if err := c.cache.Connection.SaveAccessRefreshTokens(userID, deviceID, accessToken, refreshToken); err != nil {
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

func generateAccessRefreshTokens(userID uint, deviceID string) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["device_id"] = deviceID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // access token to expire in 1 day
	claims["r"] = uuid.New().String()                     // randomize jwt even if two requests came in the same second (useful in testing)
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN")))
	if err != nil {
		logger.LogFailure(err, "Error generating the access token")
		return "", "", err
	}
	claims["exp"] = time.Now().Add(time.Hour * 24 * 90).Unix() // refresh token to expire in 90 days
	refreshToken, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN")))
	if err != nil {
		logger.LogFailure(err, "Error generating the refresh token")
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (c *Controller) updateAccessRefreshTokens(userID uint, deviceID, accessToken, refreshToken string) error {
	userTokens := &models.UserTokens{
		UserID:       userID,
		DeviceID:     deviceID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err := c.datasource.SaveAccessRefreshTokens(userTokens)
	if err != nil {
		return err
	}
	_ = c.cache.Connection.SaveAccessRefreshTokens(userID, deviceID, accessToken, refreshToken)
	return nil
}
