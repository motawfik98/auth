package database

import (
	"backend-auth/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type DB struct {
	connection *gorm.DB
}

func (db *DB) SetDBConnection(gormDB *gorm.DB) {
	db.connection = gormDB
}

func InitializeConnection() (*gorm.DB, error) {
	connectionStr := os.ExpandEnv("${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/?parseTime=true")
	connection, err := gorm.Open(mysql.Open(connectionStr))
	if err == nil {
		createDB(connection)
		connection.AutoMigrate(&models.User{}, &models.UserTokens{}, &models.UsedRefreshToken{})
	}

	return connection, err
}

func (db *DB) InitializeTransaction(fc func(tx *gorm.DB) error) error {
	return db.connection.Transaction(fc)
}

func createDB(connection *gorm.DB) {
	createDBCommand := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", os.Getenv("DB_NAME"))
	connection.Exec(createDBCommand)
	useDBCommand := fmt.Sprintf("USE %s", os.Getenv("DB_NAME"))
	connection.Exec(useDBCommand)
}
