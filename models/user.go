package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email           string `gorm:"unique" form:"email" validate:"required,email"`
	EmailVerified   bool   `gorm:"default:0"`
	Password        string `form:"password" validate:"required,strongPassword,eqfield=ConfirmPassword"`
	ConfirmPassword string `form:"confirmPassword" gorm:"-" validate:"required"`
	FullName        string `form:"fullName"`
	PhoneNumber     string `form:"phoneNumber"`
}
