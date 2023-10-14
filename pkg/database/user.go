package database

import (
	"backend-auth/pkg/models"
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
	return db.connection.Save(userTokens).Error
}

func (db *DB) MarkRefreshTokenAsUsed(token *models.UsedRefreshToken) error {
	return db.connection.Create(token).Error
}

func (db *DB) GetUsersCount() int64 {
	var count int64 = 0
	db.connection.Model(&models.User{}).Count(&count)
	return count
}

func (db *DB) IsUsedRefreshToken(refreshToken string) (bool, error) {
	var count int64 = 0
	err := db.connection.Model(&models.UsedRefreshToken{}).Where("refresh_token = ?", refreshToken).Count(&count).Error
	return count > 0, err
}

func (db *DB) GetGeneratedRefreshToken(token string) *models.GeneratedRefreshToken {
	generatedToken := new(models.GeneratedRefreshToken)
	db.connection.Where("refresh_token = ?", token).First(generatedToken)
	return generatedToken
}
