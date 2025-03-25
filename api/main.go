package main

import (
	"log"
	"net/http"
	"os"
	// "tickets/models"
)

type application struct{}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	port := getEnv("PORT", "5555")

	app := &application{}

	log.Printf("Server starting on port %s", port)
	err := http.ListenAndServe(":"+port, app.routes())
	log.Fatal(err.Error())
	os.Exit(1)

}
