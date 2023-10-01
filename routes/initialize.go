package routes

import (
	"backend-auth/controllers"
	"backend-auth/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"os"
	"slices"
)

func InitializeRoutes(e *echo.Echo, controller *controllers.Controller) {
	e.Use(cORSMiddleware())

	config := echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_ACCESS_TOKEN")),
		Skipper: func(c echo.Context) bool {
			skippedPaths := []string{"/refresh-tokens", "/signup", "/login", "/ping"}
			return slices.Contains(skippedPaths, c.Path())
		},
		SuccessHandler: decodeJWT,
	}
	refreshConfig := echojwt.Config{
		SigningKey:     []byte(os.Getenv("JWT_REFRESH_TOKEN")),
		SuccessHandler: decodeJWT,
	}
	e.Use(echojwt.WithConfig(config))
	e.GET("/refresh-tokens", controller.RefreshTokens, echojwt.WithConfig(refreshConfig))
	e.GET("/ping", handlers.Ping)

	e.GET("/users/count", controller.GetUsersCount)
	e.POST("/user", controller.CreateUser)
}
