package dev

import (
	"fmt"
	"github.com/joho/godotenv"
)

var directory = "configs/dev"

func LoadGlobalEnvFile() {
	if err := godotenv.Load(fmt.Sprintf("%s/.env", directory)); err != nil {
		panic(err)
	}
}

func LoadWorkersEnvFile() {
	if err := godotenv.Load(fmt.Sprintf("%s/.worker.env", directory)); err != nil {
		panic(err)
	}
}
