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

	e.Use(JWTMiddleware())
	e.GET("/refresh-tokens", controller.RefreshTokens, JWTRefreshMiddleware())
	e.GET("/ping", handlers.Ping)

	e.GET("/users/count", controller.GetUsersCount)
	e.POST("/user", controller.CreateUser)
}

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
