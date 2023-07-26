package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kutoru/chanl-backend/glb"
	"github.com/kutoru/chanl-backend/models"
	"github.com/kutoru/chanl-backend/tokens"
	"golang.org/x/crypto/bcrypt"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("createUser called")

	// Getting the user info from request and checking if the info format is valid
	// For now there are only length checks, but they could be more complicated if necessary

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		log.Printf("Could not decode the user struct: %v\n", err)
		return
	}

	if len(user.Name) < 4 {
		http.Error(w, "The username is too short", http.StatusBadRequest)
		return
	}

	if len(user.Name) > 20 {
		http.Error(w, "The username is too long", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 4 {
		http.Error(w, "The password is too short", http.StatusBadRequest)
		return
	}

	if len(user.Password) > 20 {
		http.Error(w, "The password is too long", http.StatusBadRequest)
		return
	}

	// Checking if there is already a user in the db with the same username

	result, err := glb.DB.Query(`
		SELECT id FROM users WHERE name = ?;
	`, user.Name)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Printf("Could not fetch existing users from the DB: %v\n", err)
		return
	}
	defer result.Close()

	exists := result.Next()
	if exists {
		http.Error(w, "The name was already taken", http.StatusBadRequest)
		return
	}

	// Generating password hash

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Printf("Could not generate hash: %v", err)
		return
	}

	// Inserting the user into the DB

	_, err = glb.DB.Exec(`
		INSERT INTO users (name, password, created_at)
		VALUES (?, ?, now());
	`, user.Name, hash)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Printf("Could not insert the new user into the DB: %v\n", err)
		return
	}

	// Getting the newly inserted user from the DB

	result, err = glb.DB.Query(`
		SELECT * FROM users WHERE name = ?;
	`, user.Name)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}
	defer result.Close()

	if !result.Next() {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, nil)
		return
	}

	err = user.ScanFromResult(result)
	user.Password = ""
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	// Inserting the user's private and personal channels

	_, err = glb.DB.Exec(`
		INSERT INTO channels (owner_id, parent_id, name, type, created_at)
		VALUES (?, 1, ?, 'pr', now());
	`, user.ID, fmt.Sprintf("%v's private channel", user.Name))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	_, err = glb.DB.Exec(`
		INSERT INTO channels (owner_id, parent_id, name, type, created_at)
		VALUES (?, 1, ?, 'pe', now());
	`, user.ID, fmt.Sprintf("%v's personal channel", user.Name))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	// Getting the newly inserted channel IDs

	result, err = glb.DB.Query(`
		SELECT * FROM channels WHERE owner_id = ?;
	`, user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}
	defer result.Close()

	if !result.Next() {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, nil)
		return
	}

	var privateChannel models.Channel
	err = privateChannel.ScanFromResult(result)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	if !result.Next() {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, nil)
		return
	}

	var personalChannel models.Channel
	err = personalChannel.ScanFromResult(result)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	// Inserting the initial joined channels (global, private and personal) for the user

	_, err = glb.DB.Exec(`
		INSERT INTO joined_channels VALUES (?, 1, 0, now());
	`, user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	_, err = glb.DB.Exec(`
		INSERT INTO joined_channels VALUES (?, ?, 1, now());
	`, user.ID, privateChannel.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	_, err = glb.DB.Exec(`
		INSERT INTO joined_channels VALUES (?, ?, 1, now());
	`, user.ID, personalChannel.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	// Creating the auth token and sending the response

	cookie, err := tokens.CreateAuthToken(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		deleteUserInfo(user.Name, err)
		return
	}

	http.SetCookie(w, cookie)
	sendJSONResponse(w, user)
}

// Attempts deletion of user info that gets inserted during the createUser function from the DB
func deleteUserInfo(username string, err error) {
	log.Printf("Deleting user info (%v). Error: %v\n", username, err)

	// First getting the user's id. In theory, if this fails, it would mean that nothing got inserted into the DB

	result, err := glb.DB.Query(`
		SELECT * FROM users WHERE name = ?;
	`, username)
	if err != nil {
		log.Println("Could not get the user from the DB")
		log.Println(err)
		return
	}
	defer result.Close()

	if !result.Next() {
		log.Println("Could not get the user from the DB")
		return
	}

	var user models.User
	err = user.ScanFromResult(result)
	user.Password = ""
	if err != nil {
		log.Println("Could not scan the user from the result")
		log.Println(err)
		return
	}

	// Deleting the user's joined channel info

	_, err = glb.DB.Exec(`
		DELETE FROM joined_channels WHERE user_id = ?;
	`, user.ID)
	if err != nil {
		log.Println("Could not delete the user's joined channels info from the DB")
		log.Println(err)
	} else {
		log.Println("Deleted the user's joined channels info")
	}

	// Deleting the user's private and personal channels

	_, err = glb.DB.Exec(`
		DELETE FROM channels WHERE owner_id = ? AND type = 'pr';
	`, user.ID, user.ID)
	if err != nil {
		log.Println("Could not delete the user's private channel")
		log.Println(err)
	} else {
		log.Println("Deleted the user's private channel")
	}

	_, err = glb.DB.Exec(`
		DELETE FROM channels WHERE owner_id = ? AND type = 'pe';
	`, user.ID, user.ID)
	if err != nil {
		log.Println("Could not delete the user's personal channel")
		log.Println(err)
	} else {
		log.Println("Deleted the user's personal channel")
	}

	// Deleting the user

	_, err = glb.DB.Exec(`
		DELETE FROM users WHERE id = ?;
	`, user.ID)
	if err != nil {
		log.Println("Could not delete the user")
		log.Println(err)
	} else {
		log.Println("Deleted the user")
	}
}
