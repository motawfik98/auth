package models

import (
	"gorm.io/gorm"
	"time"
)

type UsedRefreshToken struct {
	gorm.Model
	UserID             uint
	RefreshToken       string
	RefreshTokenExpiry time.Time
}
