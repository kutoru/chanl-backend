package endpoints

import (
	"log"
	"net/http"

	"github.com/kutoru/chanl-backend/glb"
	"github.com/kutoru/chanl-backend/models"
	"github.com/kutoru/chanl-backend/tokens"
)

// Auth with already existing _ltat_ token
func authWithToken(w http.ResponseWriter, r *http.Request) {
	log.Println("authWithToken called")

	// Checking if the request's token is valid

	userId, token, err := tokens.CheckAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Getting the user's info

	result, err := glb.DB.Query(`
		SELECT * FROM users WHERE id = ?;
	`, userId)
	if err != nil {
		http.Error(w, "Could not fetch existing user from the DB", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer result.Close()

	if !result.Next() {
		http.Error(w, "Could not find the user in the DB", http.StatusBadRequest)
		return
	}

	var user models.User
	err = user.ScanFromResult(result)
	user.Password = ""
	if err != nil {
		http.Error(w, "Could not scan the user from the DB result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Updating the token and sending the response

	cookie, err := tokens.UpdateAuthToken(user.ID, token)
	if err != nil {
		http.Error(w, "Could not update the auth token", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	http.SetCookie(w, cookie)
	sendJSONResponse(w, user)
}
