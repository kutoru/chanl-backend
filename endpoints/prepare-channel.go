package endpoints

import (
	"log"
	"net/http"
	"strconv"
)

func prepareWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("prepareWebsocketConnection called")

	// Getting user id

	userIdString := r.Header.Get("User-ID")
	userId := 0
	var err error
	if userIdString != "" {
		userId, err = strconv.Atoi(userIdString)
		if err != nil {
			http.Error(w, "Could not convert user id to an int", http.StatusBadRequest)
			return
		}
	}

	// Getting channel id

	channelId, err := getMuxVar(r, "CHANNEL_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Creating and sending the cookie

	cookie, err := activeClients.AddNewUser(channelId, userId, true)
	if err != nil {
		http.Error(w, "Could not create required cookie", http.StatusBadRequest)
		log.Println(err)
		return
	}

	log.Printf("Created cookie: %v\n", cookie)

	http.SetCookie(w, cookie)
	sendJSONResponse(w, true)
}
