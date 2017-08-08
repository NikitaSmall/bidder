package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	prepareDotEnv()
	configureDatabase()
	migrateDatabase()
}

func prepareDotEnv() {
	// make sure we load dev config only in default (debug) mode
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()

		if err != nil {
			log.Println("Error loading .env file")
		}
	}
}
