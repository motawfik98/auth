package database

import (
	"backend-auth/models"
	"os/user"
)

func (db *DB) CreateUser(user user.User) error {
	return nil
}

func (db *DB) GetUsersCount() int64 {
	var count int64 = 0
	db.db.Model(&models.User{}).Count(&count)
	return count
}
