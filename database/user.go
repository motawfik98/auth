package database

import (
	"backend-auth/models"
	"strings"
)

func (db *DB) CreateUser(user *models.User) error {
	err := db.connection.Create(user).Error
	if err != nil && strings.HasPrefix(err.Error(), "Error 1062") {
		return &DuplicateEmailError{}
	}
	return err
}

func (db *DB) SaveAccessRefreshTokens(userTokens *models.UserTokens) error {
	// same as upsert
	return db.connection.Save(userTokens).Error
}

func (db *DB) GetUsersCount() int64 {
	var count int64 = 0
	db.connection.Model(&models.User{}).Count(&count)
	return count
}
