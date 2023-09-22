package routes

import (
	"backend-auth/controllers"
	"backend-auth/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeRoutes(e *echo.Echo, controller *controllers.Controller) {
	cORSMiddleware(e)
	e.GET("/", handlers.Ping)
	e.GET("/users/count", controller.GetUsersCount)
	e.POST("/user", controller.CreateUser)
}

func cORSMiddleware(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:1323",
		},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin},
	}))
}
