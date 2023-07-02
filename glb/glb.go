package glb

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
}
