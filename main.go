package main

import (
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/logger"
	"backend-auth/routes"
	"backend-auth/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()

	port, exist := os.LookupEnv("PORT")
	if os.Getenv("ENV") == "dev" {
		port = "1322"
	} else {
		if !exist {
			port = "1323"
		}
	}

	dbConnection, err := database.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize DB connection")
		panic(err)
	}

	db := new(database.DB)
	db.SetDBConnection(dbConnection)

	controller := new(controllers.Controller)
	controller.SetDB(db)

	customValidator := &utils.CustomValidator{Validator: validator.New()}
	customValidator.TranslateErrors()
	e.Validator = customValidator

	routes.InitializeRoutes(e, controller)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
