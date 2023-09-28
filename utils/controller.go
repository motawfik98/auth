package utils

import (
	"backend-auth/cache"
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

	cacheConnection, err := cache.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize cache connection")
		panic(err)
	}

	db := new(database.DB)
	db.SetDBConnection(dbConnection)
	cacheObj := new(cache.Cache)
	cacheObj.SetCache(cacheConnection)

	controller := new(controllers.Controller)
	controller.SetDatasource(db)
	controller.SetCache(cacheObj)
	return controller
}

func InitializeValidator() *CustomValidator {
	customValidator := &CustomValidator{Validator: validator.New()}
	customValidator.TranslateErrors()
	return customValidator
}
