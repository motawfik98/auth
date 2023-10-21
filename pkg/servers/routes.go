package servers

import (
	"backend-auth/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func (s *Server) InitializeRoutes(e *echo.Echo) {
	e.Use(middleware.CORSMiddleware())

	e.Use(middleware.JWTMiddleware())
	e.GET("/refresh-tokens", s.RefreshTokens, middleware.JWTRefreshMiddleware())
	e.GET("/ping", s.Ping)

	e.GET("/users/count", s.GetUsersCount)
	e.POST("/user", s.CreateUser)
	e.POST("/login", s.Login)
}
