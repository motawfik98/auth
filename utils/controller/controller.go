package controller

import (
	"backend-auth/cache"
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/logger"
	"backend-auth/utils"
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

	cacheObj := new(cache.Cache)
	err = cacheObj.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize cache connection")
		panic(err)
	}

	controller := new(controllers.Controller)
	controller.SetDatasource(db)
	controller.SetCache(cacheObj)
	return controller
}

func InitializeValidator() *utils.CustomValidator {
	customValidator := &utils.CustomValidator{Validator: validator.New()}
	customValidator.TranslateErrors()
	return customValidator
}
