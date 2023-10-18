package main

import (
	"backend-auth/configs/dev"
	"backend-auth/internal/utils/connection"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()

	var port string
	if os.Getenv("ENV") == "dev" {
		dev.LoadGlobalEnvFile()
		port = "1322"
	} else {
		port = "1323"
	}
	envPort, exist := os.LookupEnv("PORT")
	if exist {
		port = envPort
	}

	server := connection.InitializeServer()

	e.Validator = connection.InitializeValidator()

	server.InitializeRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
