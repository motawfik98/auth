package main

import (
	"auth/handlers"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()
	e.GET("/", handlers.Ping)

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
