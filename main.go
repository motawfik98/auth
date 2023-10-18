package main

import (
	"backend-auth/routes"
	controllerUtil "backend-auth/utils/controller"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()

	port, exist := os.LookupEnv("PORT")
	if os.Getenv("ENV") == "dev" {
		if err := godotenv.Load(); err != nil {
			fmt.Println(err.Error())
		}
		port = "1322"
	} else {
		if !exist {
			port = "1323"
		}
	}

	controller := controllerUtil.InitializeController()

	e.Validator = controllerUtil.InitializeValidator()

	routes.InitializeRoutes(e, controller)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
