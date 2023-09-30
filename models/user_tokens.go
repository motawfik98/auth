package models

import "gorm.io/gorm"

type UserTokens struct {
	gorm.Model
	UserID       uint   `gorm:"index:idx_user_device"`
	DeviceID     string `gorm:"index:idx_user_device"`
	AccessToken  string
	RefreshToken string
}
