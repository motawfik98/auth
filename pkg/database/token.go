package database

import (
	"backend-auth/pkg/models"
	"time"
)

func (db *DB) CleanupUserTokens(exp time.Time) error {
	return db.connection.Where("refresh_token_expiry < ?", exp).Delete(&models.UserTokens{}).Error
}

func (db *DB) CleanupUsedRefreshToken(exp time.Time) error {
	return db.connection.Where("refresh_token_expiry < ?", exp).Delete(&models.UsedRefreshToken{}).Error
}

func (db *DB) CleanupGeneratedRefreshToken(exp time.Time) error {
	return db.connection.Where("refresh_token_expiry < ?", exp).Delete(&models.GeneratedRefreshToken{}).Error
}

func (db *DB) CleanupInvalidatedRefreshToken(exp time.Time) error {
	return db.connection.Where("refresh_token_expiry < ?", exp).Delete(&models.InvalidatedRefreshToken{}).Error
}
