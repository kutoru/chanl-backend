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

	r.HandleFunc("/api/global", getGlobalChannel).Methods("GET")
	// r.HandleFunc("/api/connect-to-channel/{CHANNEL_ID}", connectToChannel)
	r.HandleFunc("/api/private/{USER_ID}", getPrivateChannel).Methods("GET")

	// taken from other project, delete later
	// r.HandleFunc("/api/posts", getPosts).Methods("GET")
	// r.HandleFunc("/api/posts", checkTokenWrapper(addPost)).Methods("POST")
	// r.HandleFunc("/api/posts/{POST_ID}", checkTokenWrapper(deletePost)).Methods("DELETE")
	// r.HandleFunc("/api/posts/{POST_ID}/comments", checkTokenWrapper(addComment)).Methods("POST")
	// r.HandleFunc("/api/auth/login", getTokenFromPassword).Methods("POST")
	// r.HandleFunc("/api/auth/create-user", addUser).Methods("POST")
	// r.HandleFunc("/api/auth/token", checkTokenWrapper(getTokenFromToken)).Methods("GET")
	// r.HandleFunc("/api/users/{USERNAME}", checkTokenWrapper(getUser)).Methods("GET")

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
