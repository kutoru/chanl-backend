package endpoints

import (
	"log"
	"net/http"
)

func connectToChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("connectToChannel called")
}
