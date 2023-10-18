package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"os"
	"slices"
	"strconv"
	"time"
)

func JWTMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_ACCESS_TOKEN")),
		Skipper: func(c echo.Context) bool {
			skippedPaths := []string{"/refresh-tokens", "/signup", "/login", "/ping"}
			return slices.Contains(skippedPaths, c.Path())
		},
		SuccessHandler: decodeJWT,
	}
	return echojwt.WithConfig(config)
}

func JWTRefreshMiddleware() echo.MiddlewareFunc {
	refreshConfig := echojwt.Config{
		SigningKey:     []byte(os.Getenv("JWT_REFRESH_TOKEN")),
		SuccessHandler: decodeJWT,
	}
	return echojwt.WithConfig(refreshConfig)
}

func decodeJWT(c echo.Context) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := int(claims["id"].(float64))
	deviceID := claims["device_id"].(string)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0).Unix()
	c.Request().Header.Set("x-user-id", strconv.Itoa(id))
	c.Request().Header.Set("x-device-id", deviceID)
	c.Request().Header.Set("x-token-expiry", strconv.Itoa(int(expirationTime)))
	c.Response().Before(func() { // remove the added headers before writing (returning) the response
		c.Request().Header.Del("x-user-id")
		c.Request().Header.Del("x-device-id")
		c.Request().Header.Del("x-token-expiry")
	})
}
