package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string `gorm:"unique" json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
}
