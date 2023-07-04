package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/websocket-test", webSocketTest).Methods("GET")

	r.HandleFunc("/api/get-channel/{CHANNEL_ID}", getChannel).Methods("GET")
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
