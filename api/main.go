package main

import (
	"log"
	"net/http"
	"os"
	"tickets/models"
)

type application struct {
	db *models.Database
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	httpPort := getEnv("PORT", "5555")

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "tickets")

	db, err := models.DB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatal("Could not establish connection to database")
		os.Exit(1)
	}

	app := &application{db: db}

	log.Printf("Server starting on port %s", httpPort)
	err = http.ListenAndServe(":"+httpPort, app.routes())
	log.Fatal(err.Error())
	os.Exit(1)

}
