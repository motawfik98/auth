package servers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "Success")
}
