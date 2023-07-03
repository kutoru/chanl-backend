package main

import (
	"log"
	"main/dbmanager"
	"main/endpoints"
	"main/glb"
	"main/models"
	"net/http"
)

func main() {
	glb.LoadEnv()
	dbmanager.ConnectToDB()
	dbmanager.ResetDB()
	dbmanager.TestDB()
	defer dbmanager.DisconnectFromDB()

	r := endpoints.GetRouter()
	http.Handle("/", &models.RouterDec{Router: r})
	log.Println("API server is listening on http://localhost:4000")
	log.Panicln(http.ListenAndServe(":4000", nil))
}
