package main

import (
	"backend-auth/internal/logger"
	"backend-auth/internal/utils/connection"
)

func main() {
	cronjob := connection.InitializeCronJob()
	err := cronjob.CleanupInvalidatedRefreshToken()
	if err != nil {
		logger.LogFailure(err, "error deleting expired (invalidated refresh tokens) from db")
		panic(err)
	}
}
