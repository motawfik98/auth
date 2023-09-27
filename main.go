package main

import (
	"backend-auth/routes"
	"backend-auth/utils"
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

	controller := utils.InitializeController()

	e.Validator = utils.InitializeValidator()

	routes.InitializeRoutes(e, controller)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
