package endpoints

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kutoru/chanl-backend/glb"
	"github.com/kutoru/chanl-backend/models"
	"github.com/kutoru/chanl-backend/tokens"
)

var activeClients models.ActiveClientsWrapper

func InitializeActiveClients() {
	activeClients = models.ActiveClientsWrapper{}
	activeClients.Initialize()
}

func webSocketCheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("origin")
	for _, allowOrigin := range glb.AllowedOrigins {
		if origin == allowOrigin {
			return true
		}
	}
	return false
}

// Might need to change the buffer sizes later
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     webSocketCheckOrigin,
}

func connectToChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("connectToChannel called")

	// Getting user id and channel id from cookies

	cookie, err := r.Cookie("_wst_")
	if err != nil {
		http.Error(w, "Could not get the required cookie", http.StatusBadRequest)
		log.Println(err)
		return
	}

	log.Printf("Got cookie: %v\n", cookie)

	expectedChannelId, userId, err := tokens.ParseWebsocketToken(cookie.Value)
	if err != nil {
		http.Error(w, "Could not parse the cookie", http.StatusBadRequest)
		log.Println(err)
		return
	}

	log.Printf("Got cookie values: %v, %v\n", expectedChannelId, userId)

	// Getting requested channel id and comparing it with cookie's channel id

	channelId, err := getMuxVar(r, "CHANNEL_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if channelId != expectedChannelId {
		http.Error(w, "Expected channel id != requested channel id", http.StatusBadRequest)
		return
	}

	activeClients.RemoveCookie(channelId, userId, true)

	// Checking if the channel requires auth

	result, err := glb.DB.Query(`
		SELECT * FROM channels WHERE id = ?;
	`, channelId)
	if err != nil || !result.Next() {
		http.Error(w, "Could not get channel info from db", http.StatusBadRequest)
		log.Println(err)
		return
	}

	var channel models.Channel
	err = channel.ScanFromResult(result)
	result.Close()
	if err != nil {
		http.Error(w, "Could not parse channel info from db result", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if userId <= 0 &&
		(channel.Type != models.Global && channel.Type != models.Server) {
		http.Error(w, "The channel requires auth, but the user id is 0", http.StatusBadRequest)
		log.Println(userId, channelId, channel)
		return
	}

	// If everything is fine, upgrading to a websocket connection and saving the client

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed")
		log.Println(err)
		return
	}

	err = activeClients.AssignConnection(channelId, userId, conn, true)
	if err != nil {
		log.Println("Could not assign the connection")
		log.Println(err)
		return
	}

	// RemoveUser also closes the connection
	defer activeClients.RemoveUser(channelId, userId, true)

	log.Println("Connection upgraded")
	activeClients.Print(true)

	// Sending the existing channel messages to the client

	messages := []models.Message{}

	result, err = glb.DB.Query(`
		SELECT messages.*, users.name FROM messages
		INNER JOIN users ON
		(messages.channel_id = ?) AND
		(messages.user_id = users.id)
		ORDER BY messages.id DESC;
	`, channelId)
	if err != nil {
		log.Println("Could not get messages result")
		log.Println(err)
		return
	}

	for result.Next() {
		var message models.Message
		err := message.ScanFromResult(result, true)
		if err != nil {
			log.Println("Could not scan a message from db result")
			log.Println(err)
			return
		}

		messages = append(messages, message)
	}

	result.Close()

	err = conn.WriteJSON(messages)
	if err != nil {
		log.Println("Could not send the messages response")
		log.Println(err)
		return
	}

	for {
		// The for loop pauses on the ReadJSON line and waits for a new message from the client
		// Breaking out of this loop basically means closing the connection

		err = listenForMessage(userId, channelId, conn)
		if err != nil {
			log.Println(err)
			break
		}
	}

	log.Printf("connectToChannel finished\n\n")
}

func listenForMessage(userId int, channelId int, conn *websocket.Conn) error {

	// Getting a message

	var message models.Message
	err := conn.ReadJSON(&message)
	if err != nil {
		return fmt.Errorf("could not read the received message: %v", err)
	}

	if userId <= 0 ||
		userId != message.UserID || channelId != message.ChannelID ||
		len(message.Text) < 1 || len(message.Text) > 1024 {
		return fmt.Errorf("got invalid message. %v, %v, %v", userId, channelId, message)
	}

	// Inserting it into the db

	_, err = glb.DB.Exec(`
		INSERT INTO messages (user_id, channel_id, text, sent_at)
		VALUES (?, ?, ?, now());
	`, message.UserID, message.ChannelID, message.Text)
	if err != nil {
		return fmt.Errorf("could not insert the message into the db: %v", err)
	}

	log.Printf("Message processed: %v\n", message)

	// Getting additional info for the message

	result, err := glb.DB.Query(`
		SELECT messages.id, messages.sent_at, users.name FROM messages
		INNER JOIN users ON (messages.channel_id = ?) AND (messages.user_id = users.id)
		WHERE users.id = ? ORDER BY messages.id DESC;
	`, channelId, userId)
	if err != nil || !result.Next() {
		return fmt.Errorf("could not get the additional message info from the db: %v", err)
	}

	err = result.Scan(&message.ID, &message.SentAt, &message.UserName)
	result.Close()
	if err != nil {
		return fmt.Errorf("could not parse additional message info from db result: %v", err)
	}

	log.Printf("Additional message info processed: %v\n", message)

	// Sending the message to all clients in this channel
	// Making an array of messages because that is what the frontend accepts

	messageArray := []models.Message{message}
	err = activeClients.SendMessagesToChannel(channelId, messageArray, true)
	if err != nil {
		return fmt.Errorf("could not send message to all activeClients: %v", err)
	}

	log.Print("connectToChannel for loop finished without errors\n\n")
	return nil
}
