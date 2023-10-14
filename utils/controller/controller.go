package controller

import (
	"backend-auth/controllers"
	"backend-auth/internal/logger"
	"backend-auth/pkg/cache"
	"backend-auth/pkg/database"
	"backend-auth/pkg/messaging"
	"backend-auth/utils"
	"fmt"
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

	messagingObj := new(messaging.Messaging)
	err = messagingObj.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize messaging connection")
		panic(err)
	}
	queueName, err := messagingObj.CreateQueues()
	if err != nil {
		logger.LogFailure(err, fmt.Sprintf("Failed to create queue: %s", queueName))
		panic(err)
	}

	controller := new(controllers.Controller)
	controller.SetDatasource(db)
	controller.SetCache(cacheObj)
	controller.SetMessaging(messagingObj)
	return controller
}

func InitializeValidator() *utils.CustomValidator {
	customValidator := &utils.CustomValidator{Validator: validator.New()}
	customValidator.TranslateErrors()
	return customValidator
}
