package utils

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

func FetchUserAndDevice(ctx echo.Context) (uint, string, error) {
	deviceID := ctx.Request().Header.Get("x-device-id")
	if deviceID == "" {
		return 0, "", errors.New("cannot get device ID from refresh token")
	}
	strUserID := ctx.Request().Header.Get("x-user-id")
	userID, err := strconv.ParseUint(strUserID, 10, 64)
	if err != nil {
		return 0, "", err
	}
	return uint(userID), deviceID, nil
}

func FetchRefreshTokenAndExpiry(ctx echo.Context) (string, time.Time, error) {
	oldRefreshToken := ctx.Request().Header.Get("Authorization")[7:] // remove `bearer ` prefix
	oldTokenExpiry := ctx.Request().Header.Get("x-token-expiry")
	exp, err := strconv.ParseInt(oldTokenExpiry, 10, 64)
	if err != nil {
		return "", time.Time{}, err
	}
	tm := time.Unix(exp, 0)
	return oldRefreshToken, tm, nil
}
