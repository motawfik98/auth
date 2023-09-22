package controllers

import (
	"backend-auth/database"
	"backend-auth/logger"
	"backend-auth/models"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
	err = c.db.CreateUser(user)
	if errors.Is(err, &database.DuplicateEmailError{}) {
		return ctx.JSON(http.StatusConflict, echo.Map{})
	}
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}

	accessToken, refreshToken, err := generateAccessRefreshTokens(user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{})
	}

	return ctx.JSON(http.StatusCreated, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *Controller) GetUsersCount(ctx echo.Context) error {
	count := c.db.GetUsersCount()
	return ctx.String(http.StatusOK, strconv.FormatInt(count, 10))
}

func generateAccessRefreshTokens(user *models.User) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // access token to expire in 1 day
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
