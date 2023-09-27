package utils

import (
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/logger"
	"github.com/go-playground/validator/v10"
)

func InitializeController() *controllers.Controller {
	dbConnection, err := database.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize DB connection")
		panic(err)
	}

	db := new(database.DB)
	db.SetDBConnection(dbConnection)

	controller := new(controllers.Controller)
	controller.SetDB(db)
	return controller
}

func InitializeValidator() *CustomValidator {
	customValidator := &CustomValidator{Validator: validator.New()}
	customValidator.TranslateErrors()
	return customValidator
}
