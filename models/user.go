package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email           string       `gorm:"unique" json:"email" validate:"required,email"`
	EmailVerified   bool         `gorm:"default:0" json:"-"`
	Password        string       `json:"password" validate:"required,strongPassword,eqfield=ConfirmPassword"`
	ConfirmPassword string       `json:"confirm_password" gorm:"-" validate:"required"`
	FullName        string       `json:"full_name"`
	PhoneNumber     string       `json:"phone_number"`
	UserTokens      []UserTokens `json:"-"`
}
