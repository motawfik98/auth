package main

import (
	"backend-auth/routes"
	"backend-auth/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
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

	customValidator := &utils.CustomValidator{Validator: validator.New()}
	customValidator.TranslateErrors()
	e.Validator = customValidator

	routes.InitializeRoutes(e, controller)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
