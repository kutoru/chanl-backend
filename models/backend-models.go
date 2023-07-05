package models

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ClientWrapper struct {
	Cookie *http.Cookie
	Conn   *websocket.Conn
}

// map[channelId]map[userId]*ClientWrapper
type ActiveClientsWrapper struct {
	clientsWrapper map[int]map[int]*ClientWrapper
}

type RouterDec struct {
	Router *mux.Router
}

func (acw *ActiveClientsWrapper) Initialize() {
	acw.clientsWrapper = make(map[int]map[int]*ClientWrapper)
}

func (acw *ActiveClientsWrapper) AddChannel(channelId int) {
	_, exists := acw.clientsWrapper[channelId]
	if !exists {
		acw.clientsWrapper[channelId] = make(map[int]*ClientWrapper)
	}
}

func (acw *ActiveClientsWrapper) RemoveChannel(channelId int) {
	channel, exists := acw.clientsWrapper[channelId]
	if exists {
		for userId, userInfo := range channel {
			userInfo.Conn.Close()
			delete(acw.clientsWrapper[channelId], userId)
		}
		delete(acw.clientsWrapper, channelId)
	}
}

// If userId is 0, changes it to the lowest negative number that is not present in the map.
// Returns the user id that it added
func (acw *ActiveClientsWrapper) checkUserID(channelId int, userId int) int {
	channel, exists := acw.clientsWrapper[channelId]
	if !exists {
		if userId > 0 {
			return userId
		} else {
			return -1
		}
	}

	if userId > 0 {
		_, exists := channel[userId]
		if exists {
			acw.RemoveUser(channelId, userId)
		}
		return userId
	}

	minId := -1
	for userId := range acw.clientsWrapper[channelId] {
		if userId < minId {
			minId = userId - 1
		}
	}

	return minId
}

var currentlyAdding bool = false

func (acw *ActiveClientsWrapper) AddNewUser(channelId int, userId int) *http.Cookie {
	for currentlyAdding {
		time.Sleep(10 * time.Millisecond)
	}

	currentlyAdding = true
	acw.AddChannel(channelId)
	userId = acw.checkUserID(channelId, userId)

	cookie := &http.Cookie{
		Name:     "_sc_",
		Value:    fmt.Sprintf("%d.%d", channelId, userId),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   5,
		SameSite: http.SameSiteNoneMode,
	}

	acw.clientsWrapper[channelId][userId] = &ClientWrapper{
		Cookie: cookie,
		Conn:   nil,
	}

	currentlyAdding = false
	return cookie
}

// Assign a connection object to a user
func (acw *ActiveClientsWrapper) AssignConnection(channelId int, userId int, conn *websocket.Conn) error {
	channel, exists := acw.clientsWrapper[channelId]
	if !exists {
		return fmt.Errorf("channel id (%v) does not exist in the map", channelId)
	}

	user, exists := channel[userId]
	if !exists {
		return fmt.Errorf("user id (%v) does not exist in the channel map (%v)", userId, channelId)
	} else {
		user.Conn = conn
		return nil
	}
}

func (acw *ActiveClientsWrapper) RemoveUser(channelId int, userId int) {
	channel, channelExists := acw.clientsWrapper[channelId]
	if channelExists {
		userInfo, userExists := channel[userId]
		if userExists {
			if userInfo.Conn != nil {
				userInfo.Conn.Close()
			}
			delete(acw.clientsWrapper[channelId], userId)
		}
	}
}

func (acw *ActiveClientsWrapper) SendMessagesToChannel(channelId int, message []Message) error {
	channel, exists := acw.clientsWrapper[channelId]
	if !exists {
		return fmt.Errorf("channel id (%v) does not exits in the map", channelId)
	}

	var savedError error = nil

	for userId, userInfo := range channel {
		if userInfo.Conn != nil {
			err := userInfo.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Could not send message to: %v. %v", userId, err)
				savedError = err
			}
		}
	}

	return savedError
}

func (acw *ActiveClientsWrapper) Print() {
	log.Println("Current active clients")
	for channelId, channel := range acw.clientsWrapper {
		log.Printf("  Channel %v:\n", channelId)
		for userId, user := range channel {
			log.Printf("    User %v: %v\n", userId, *user)
		}
	}
}

// Cannot import glb here, so unfortunately have to define the origins again
var allowedOrigins = []string{"http://localhost:5000", "http://192.168.1.12:5000"}

func checkOrigin(origin string) string {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == origin {
			return origin
		}
	}
	return ""
}

func (routerDec *RouterDec) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	origin = checkOrigin(origin)

	if origin != "" {
		// General backend access
		w.Header().Set(
			"Access-Control-Allow-Origin", origin,
		)
		// Ability to set cookies
		w.Header().Set(
			"Access-Control-Allow-Credentials", "true",
		)
		// Allowed methods
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET",
		)
		// Special allowed fields in headers
		w.Header().Add(
			"Access-Control-Allow-Headers",
			"User-ID",
			// "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, User-ID",
		)
	}

	if r.Method == "OPTIONS" {
		return
	}

	routerDec.Router.ServeHTTP(w, r)
}
