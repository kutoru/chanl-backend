package endpoints

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "websocket-test.html")
}

var origins = []string{"http://localhost:4000", "http://192.168.1.12:4000"}

func webSocketCheckOrigin(r *http.Request) bool {
	origin := r.Header.Get("origin")
	for _, allowOrigin := range origins {
		if origin == allowOrigin {
			return true
		}
	}
	return false
}

func addClient(newClient *websocket.Conn) int {
	id := maxClientId
	maxClientId++
	clients = append(clients, ClientWrapper{
		ID:   id,
		Conn: newClient,
	})
	return id
}

func removeClient(clientId int) {
	for index, client := range clients {
		if client.ID == clientId {
			clen := len(clients)
			clients[index] = clients[clen-1]
			clients = clients[:clen-1]
			return
		}
	}
}

type ClientWrapper struct {
	ID   int
	Conn *websocket.Conn
}

type Message struct {
	ClientID int    `json:"clientId"`
	Text     string `json:"text"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     webSocketCheckOrigin,
}
var clients []ClientWrapper
var maxClientId int = 1
var messages []Message

func webSocketTest(w http.ResponseWriter, r *http.Request) {
	log.Println("webSocketThing called")

	// Upgrading the connection to a websocket connection

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %v\n", err)
		return
	}

	// Saving the client and sending the existing messages

	clientId := addClient(conn)
	defer conn.Close()
	defer removeClient(clientId)

	if len(messages) > 0 {
		err = conn.WriteJSON(messages)
		if err != nil {
			log.Printf("Write failed: %v", err)
			return
		}
	}

	log.Println("webSocketThing upgraded")

	// The for loop starts when there is a new message from the client
	for {
		// Getting, parsing and saving the client's message

		messageType, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read failed: %v\n", err)
			break
		} else if messageType != 1 {
			log.Printf("Read failed, got unknown message type: %v\n", messageType)
			break
		}

		message := strings.TrimSpace(string(messageBytes))
		if len(message) < 1 || len(message) > 1024 {
			continue
		}
		messages = append(messages, Message{ClientID: clientId, Text: message})

		log.Printf("Message read from %v: \"%v\"\n", clientId, message)

		// Sending the updated messages list to all clients

		log.Println("===== clients: =====")
		log.Println(clients)
		log.Println("===== messages: =====")
		log.Println(messages)

		for _, client := range clients {
			err = client.Conn.WriteJSON(messages)
			if err != nil {
				log.Printf("Write failed: %v", err)
				continue
			}
		}

		log.Print("Loop done\n\n")
	}

	log.Printf("webSocketThing done\n\n")
}
