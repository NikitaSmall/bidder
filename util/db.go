package util

import (
	"database/sql"
	"log"
	"os"

	// we want to ensure we use postgresql database
	_ "github.com/lib/pq"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"

	// we want to use file as a source for migration
	_ "github.com/mattes/migrate/source/file"
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

func migrateDatabase() {
	driver, err := postgres.WithInstance(DBConnect, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		log.Fatalf("Cannot migrate DB due to error: %s", err)
	}

	m.Steps(3)
}
