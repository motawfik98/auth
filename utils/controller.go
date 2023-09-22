package utils

import (
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/logger"
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
