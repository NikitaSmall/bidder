package util

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	prepareDotEnv()
	configureDatabase()
}

func prepareDotEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Print("Error loading .env file")
	}
}
