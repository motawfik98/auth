package main

import (
	"backend-auth/internal/logger"
	"backend-auth/internal/utils/connection"
)

func main() {
	cronjob := connection.InitializeCronJob()
	err := cronjob.CleanupUserTokens()
	if err != nil {
		logger.LogFailure(err, "error deleting expired user tokens from db")
		panic(err)
	}
}
