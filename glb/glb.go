package glb

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
)

// Frontend origins
var AllowedOrigins = []string{"http://localhost:5000", "http://192.168.1.12:5000"}
var DB *sql.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
}
