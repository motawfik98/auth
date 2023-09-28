package routes

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"strconv"
)

func decodeJWT(c echo.Context) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := int(claims["id"].(float64))
	deviceID := claims["device_id"].(string)
	c.Request().Header.Set("x-user-id", strconv.Itoa(id))
	c.Request().Header.Set("x-device-id", deviceID)
}
