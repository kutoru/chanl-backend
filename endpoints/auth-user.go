package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kutoru/chanl-backend/glb"
	"github.com/kutoru/chanl-backend/models"
	"github.com/kutoru/chanl-backend/tokens"
	"golang.org/x/crypto/bcrypt"
)

// Auth with username and password
func authUser(w http.ResponseWriter, r *http.Request) {
	log.Println("authUser called")

	// Getting user from the request

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		log.Printf("Could not decode the user struct: %v\n", err)
		return
	}

	// Getting the user's info from the DB

	result, err := glb.DB.Query(`
		SELECT * FROM users WHERE name = ?;
	`, user.Name)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Printf("Could not fetch existing user from the DB: %v\n", err)
		return
	}
	defer result.Close()

	if !result.Next() {
		http.Error(w, "This user does not exist", http.StatusBadRequest)
		return
	}

	var existingUser models.User
	err = existingUser.ScanFromResult(result)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Printf("Could not scan the user from the DB result: %v\n", err)
		return
	}

	// Checking if the password is valid

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	existingUser.Password = ""
	if err != nil {
		http.Error(w, "The password is wrong", http.StatusBadRequest)
		return
	}

	// Creating the auth token and sending the response

	cookie, err := tokens.CreateAuthToken(existingUser.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Printf("Could not create auth token: %v\n", err)
		return
	}

	http.SetCookie(w, cookie)
	sendJSONResponse(w, existingUser)
}
