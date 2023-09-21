package database

import (
	"backend-auth/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DB struct {
	db *gorm.DB
}

func (db *DB) SetDBConnection(gormDB *gorm.DB) {
	db.db = gormDB
}

var (
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Millisecond, // Slow SQL threshold
			LogLevel:      logger.Info,      // Log level
			Colorful:      true,             // Disable color
		},
	)
)

func InitializeConnection() (*gorm.DB, error) {
	connectionStr := os.ExpandEnv("${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/")
	connection, err := gorm.Open(mysql.Open(connectionStr), &gorm.Config{
		Logger: newLogger,
	})
	if err == nil {
		createDB(connection)
		connection.AutoMigrate(&models.User{})
	}

	return connection, err
}

func createDB(connection *gorm.DB) {
	createDBCommand := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", os.Getenv("DB_NAME"))
	connection.Exec(createDBCommand)
	useDBCommand := fmt.Sprintf("USE %s", os.Getenv("DB_NAME"))
	connection.Exec(useDBCommand)
}
