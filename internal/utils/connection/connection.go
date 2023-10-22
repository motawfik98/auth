package connection

import (
	"backend-auth/internal/logger"
	"backend-auth/internal/utils"
	"backend-auth/pkg/cache"
	"backend-auth/pkg/cronjobs"
	"backend-auth/pkg/database"
	"backend-auth/pkg/messaging"
	"backend-auth/pkg/servers"
	"backend-auth/pkg/workers"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func InitializeServer() *servers.Server {
	db := getDBConnection()

	cacheObj := getCacheConnection()

	messagingObj := getMessagingConnection()

	s := new(servers.Server)
	s.SetDatasource(db)
	s.SetCache(cacheObj)
	s.SetMessaging(messagingObj)
	return s
}

func InitializeWorker() *workers.Worker {
	db := getDBConnection()

	cacheObj := getCacheConnection()

	messagingObj := getMessagingConnection()

	w := new(workers.Worker)
	w.SetDatasource(db)
	w.SetCache(cacheObj)
	w.SetMessaging(messagingObj)
	return w
}

func InitializeCronJob() *cronjobs.CronJob {
	db := getDBConnection()

	cj := new(cronjobs.CronJob)
	cj.SetDatasource(db)
	return cj
}

func getDBConnection() *database.DB {
	dbConnection, err := database.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize DB connection")
		panic(err)
	}

	db := new(database.DB)
	db.SetDBConnection(dbConnection)
	return db
}

func getCacheConnection() *cache.Cache {
	cacheObj := new(cache.Cache)
	err := cacheObj.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize cache connection")
		panic(err)
	}
	return cacheObj
}

func getMessagingConnection() *messaging.Messaging {
	messagingObj := new(messaging.Messaging)
	err := messagingObj.InitializeConnection()
	if err != nil {
		logger.LogFailure(err, "Failed to initialize messaging connection")
		panic(err)
	}
	queueName, err := messagingObj.CreateQueues()
	if err != nil {
		logger.LogFailure(err, fmt.Sprintf("Failed to create queue: %s", queueName))
		panic(err)
	}
	return messagingObj
}

func InitializeValidator() *utils.CustomValidator {
	customValidator := &utils.CustomValidator{Validator: validator.New()}
	customValidator.TranslateErrors()
	return customValidator
}
