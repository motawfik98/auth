package test

import (
	"backend-auth/pkg/models"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var totalCreatedUsers int64 = 0

func getRandomEmailHandle() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := rand.New(rand.NewSource(time.Now().UnixNano() + totalCreatedUsers))
	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

func createUser(overrides map[string]interface{}, saveToDB bool, db *gorm.DB) *models.User {
	user := &models.User{
		Email:    fmt.Sprintf("%s@a.com", getRandomEmailHandle()),
		Password: "Aa!12345",
		FullName: "Mohamed Tawfik",
	}
	if val, found := overrides["email"]; found {
		user.Email = val.(string)
	}
	if val, found := overrides["password"]; found {
		user.Password = val.(string)
	}
	if val, found := overrides["full_name"]; found {
		user.FullName = val.(string)
	}
	if saveToDB {
		db.Create(user)
	}
	totalCreatedUsers++
	return user
}
