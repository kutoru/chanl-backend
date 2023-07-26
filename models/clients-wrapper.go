package models

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kutoru/chanl-backend/tokens"
)

// In theory this wait system helps avoiding any async problems
var busy bool = false

// These two functions should always be used together
func waitForBusy() {
	for busy {
		time.Sleep(time.Millisecond)
	}
	busy = true
}

func freeBusy() {
	busy = false
}

// Individual client
type ClientWrapper struct {
	Cookie *http.Cookie
	Conn   *websocket.Conn
}

// All active clients
// map[channelId]map[userId]*ClientWrapper
type ActiveClientsWrapper struct {
	clientsWrapper map[int]map[int]*ClientWrapper
}

func (acw *ActiveClientsWrapper) Initialize() {
	acw.clientsWrapper = make(map[int]map[int]*ClientWrapper)
}

func (acw *ActiveClientsWrapper) AddChannel(channelId int, wait bool) {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

	_, exists := acw.clientsWrapper[channelId]
	if !exists {
		acw.clientsWrapper[channelId] = make(map[int]*ClientWrapper)
	}
}

func (acw *ActiveClientsWrapper) RemoveChannel(channelId int, wait bool) {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

	channel, exists := acw.clientsWrapper[channelId]
	if exists {
		for userId, userInfo := range channel {
			if userInfo.Conn != nil {
				userInfo.Conn.Close()
			}
			delete(acw.clientsWrapper[channelId], userId)
		}
		delete(acw.clientsWrapper, channelId)
	}
}

// If userId is 0, changes it to the lowest negative number that is not present in the map.
// Returns the user id that it added
func (acw *ActiveClientsWrapper) checkUserID(channelId int, userId int, wait bool) int {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

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
			acw.RemoveUser(channelId, userId, false)
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

func (acw *ActiveClientsWrapper) AddNewUser(channelId int, userId int, wait bool) (*http.Cookie, error) {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

	acw.AddChannel(channelId, false)
	userId = acw.checkUserID(channelId, userId, false)
	token, err := tokens.CreateWebsocketToken(channelId, userId)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "_wst_",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   5,
		SameSite: http.SameSiteNoneMode,
	}

	acw.RemoveUser(channelId, userId, false) // Calling this just in case there are previous user sessions

	acw.clientsWrapper[channelId][userId] = &ClientWrapper{
		Cookie: cookie,
		Conn:   nil,
	}

	return cookie, nil
}

// Assign a connection object to a user
func (acw *ActiveClientsWrapper) AssignConnection(channelId int, userId int, conn *websocket.Conn, wait bool) error {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

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

func (acw *ActiveClientsWrapper) RemoveCookie(channelId int, userId int, wait bool) {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

	channel, exists := acw.clientsWrapper[channelId]
	if exists {
		userInfo, exists := channel[userId]
		if exists {
			userInfo.Cookie = nil
		}
	}
}

func (acw *ActiveClientsWrapper) RemoveUser(channelId int, userId int, wait bool) {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

	channel, channelExists := acw.clientsWrapper[channelId]
	if channelExists {
		userInfo, userExists := channel[userId]
		if userExists && userInfo.Conn != nil { // If the connection is nil, that most likely means that the user is about to connect
			userInfo.Conn.Close()
			delete(acw.clientsWrapper[channelId], userId)
		}
	}
}

func (acw *ActiveClientsWrapper) SendMessagesToChannel(channelId int, message []Message, wait bool) error {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

	channel, exists := acw.clientsWrapper[channelId]
	if !exists {
		return fmt.Errorf("channel id (%v) does not exits in the map", channelId)
	}

	var savedError error = nil

	for userId, userInfo := range channel {
		if userInfo.Conn != nil {
			err := userInfo.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Could not send message to: %v. %v\n", userId, err)
				savedError = err
			}
		}
	}

	return savedError
}

func (acw *ActiveClientsWrapper) Print(wait bool) {
	if wait {
		waitForBusy()
		defer freeBusy()
	}

	log.Println("Current active clients:")
	for channelId, channel := range acw.clientsWrapper {
		log.Printf("  Channel %v:\n", channelId)
		for userId, user := range channel {
			log.Printf("    User %v: %v\n", userId, *user)
		}
	}
}
