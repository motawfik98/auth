package models

import "gorm.io/gorm"

type UserTokens struct {
	gorm.Model
	UserID       uint   `gorm:"uniqueIndex:idx_user_device"`
	DeviceID     string `gorm:"uniqueIndex:idx_user_device;size:40"`
	AccessToken  string
	RefreshToken string
}
