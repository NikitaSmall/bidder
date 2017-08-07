package util

import (
	"database/sql"
	"log"
	"os"

	// we want to ensure we use postgresql database
	_ "github.com/lib/pq"
)

// DBConnect is a single connection for this application.
// It should have it's own connection pull.
var DBConnect *sql.DB

func configureDatabase() {
	var err error

	DBConnect, err = sql.Open("postgres", os.Getenv("POSTGRES"))

	if err != nil {
		log.Fatalf("Cannot create DB connection due to error: %s", err)
	}

	err = DBConnect.Ping()
	if err != nil {
		log.Fatalf("Cannot reach DB due to error: %s", err)
	}

	log.Println("Succesfully connected to the DB.")
}
