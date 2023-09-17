package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Success")
	})

	port, exist := os.LookupEnv("PORT")
	if os.Getenv("ENV") == "dev" {
		port = "1322"
	} else {
		if !exist {
			port = "1323"
		}
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
