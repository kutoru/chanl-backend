package main

import (
	"log"
	"net/http"

	"github.com/kutoru/chanl-backend/dbmanager"
	"github.com/kutoru/chanl-backend/endpoints"
	"github.com/kutoru/chanl-backend/glb"
	"github.com/kutoru/chanl-backend/models"
)

func main() {
	glb.LoadEnv()
	dbmanager.ConnectToDB()
	dbmanager.ResetDB()
	dbmanager.TestDB()
	defer dbmanager.DisconnectFromDB()

	endpoints.InitializeActiveClients()
	r := endpoints.GetRouter()
	http.Handle("/", &models.RouterDec{Router: r})
	log.Println("API server is listening on http://192.168.1.12:4000")
	log.Panicln(http.ListenAndServe(":4000", nil))
}
