package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/get-channel/{CHANNEL_ID}", getChannel).Methods("GET")
	r.HandleFunc("/api/prepare-channel/{CHANNEL_ID}", prepareWebsocketConnection).Methods("GET")
	r.HandleFunc("/api/connect-to-channel/{CHANNEL_ID}", connectToChannel).Methods("GET")

	// todos
	// r.HandleFunc("/api/create-channel", createChannel).Methods("POST")

	// r.HandleFunc("/api/create-message", createMessage).Methods("POST")
	// r.HandleFunc("/api/get-message", getMessage).Methods("GET")

	// r.HandleFunc("/api/create-user", createUser).Methods("POST")
	// r.HandleFunc("/api/get-user/{USER_ID}", getUser).Methods("GET")

	return r
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode a JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write the response body: %v\n", err)
	}
}

func getMuxVar(r *http.Request, varKey string) (int, error) {
	vars := mux.Vars(r)
	varString, ok := vars[varKey]
	if !ok {
		return 0, fmt.Errorf("could not get mux var (%v)", varKey)
	}

	varInt, err := strconv.Atoi(varString)
	if err != nil {
		return 0, fmt.Errorf("could not convert mux var (%v) to an int", varKey)
	}

	return varInt, nil
}
