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

	ip, err := glb.GetIP()
	if err != nil {
		log.Println(err)
	}
	port := ":4000"

	endpoints.InitializeActiveClients()
	r := endpoints.GetRouter()
	http.Handle("/", &models.RouterDec{Router: r})
	log.Printf("API server is listening on http://%v%v\n", ip, port)
	log.Panicln(http.ListenAndServe(port, nil))
}
