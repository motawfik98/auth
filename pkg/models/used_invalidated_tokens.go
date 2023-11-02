package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type UsedRefreshToken struct {
	gorm.Model
	UserID             uint
	RefreshToken       string
	RefreshTokenExpiry time.Time `gorm:"index"`
}

type GeneratedRefreshToken struct {
	gorm.Model
	UserID               uint
	RefreshToken         string
	RefreshTokenExpiry   time.Time     `gorm:"index"`
	ParentRefreshTokenID sql.NullInt64 // id for the parent refresh token
}

type InvalidatedRefreshToken struct {
	gorm.Model
	UserID             uint
	RefreshToken       string
	RefreshTokenExpiry time.Time `gorm:"index"`
}
