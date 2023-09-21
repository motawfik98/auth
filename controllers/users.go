package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (c *Controller) CreateUser() {

}

func (c *Controller) GetUsersCount(ctx echo.Context) error {
	count := c.db.GetUsersCount()
	return ctx.String(http.StatusOK, strconv.FormatInt(count, 10))
}
