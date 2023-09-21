package main

import (
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/handlers"
	"backend-auth/logger"
	"fmt"
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

	e.GET("/", handlers.Ping)
	e.GET("/users/count", controller.GetUsersCount)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
