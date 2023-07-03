package dbmanager

import (
	"database/sql"
	"fmt"
	"log"
	"main/glb"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDB() {
	dbInfo := fmt.Sprintf(
		"root:%s@tcp(localhost:3306)/%s?multiStatements=true&loc=Europe%%2FLondon",
		os.Getenv("DB_PASS"), os.Getenv("DB_NAME"),
	)

	var err error
	glb.DB, err = sql.Open("mysql", dbInfo)
	if err != nil {
		log.Fatalf("Could not open the database: %v", err)
	}

	// go can start before the database sometimes, this avoids any issues related to that
	for glb.DB.Ping() != nil {
		log.Println("Attempting connection to db")
		time.Sleep(3 * time.Second)
	}

	log.Println("Connected to db")
}

func DisconnectFromDB() {
	glb.DB.Close()
}

func ResetDB() {
	script, err := os.ReadFile("./create_db.sql")
	if err != nil {
		log.Fatalf("Could not open the db script: %v", err)
	}

	_, err = glb.DB.Exec(string(script))
	if err != nil {
		log.Fatalf("Could not execute the db script: %v", err)
	}

	log.Println("Initialized the DB")
}

func TestDB() {
	script, err := os.ReadFile("./test_db.sql")
	if err != nil {
		log.Fatalf("Could not open the db script: %v", err)
	}

	_, err = glb.DB.Exec(string(script))
	if err != nil {
		log.Fatalf("Could not execute the db script: %v", err)
	}

	log.Println("Executed the test script")
}
