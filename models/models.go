package models

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ChannelType string

const (
	Global   ChannelType = "gl"
	Private  ChannelType = "pr"
	Server   ChannelType = "se"
	Room     ChannelType = "ro"
	Personal ChannelType = "pe"
	Friend   ChannelType = "fr"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}

type Channel struct {
	ID        int         `json:"id"`
	OwnerID   int         `json:"ownerId"`
	ParentID  int         `json:"parentId"`
	Name      string      `json:"name"`
	Type      ChannelType `json:"type"`
	CreatedAt string      `json:"createdAt"`
}

type JoinedChannel struct {
	UserID    int      `json:"userId"`
	ChannelID int      `json:"channelId"`
	CanWrite  bool     `json:"canWrite"`
	JoinedAt  string   `json:"joinedAt"`
	Channel   *Channel `json:"channel"`
}

type Message struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	ChannelID int    `json:"channelId"`
	Text      string `json:"text"`
	SentAt    string `json:"sentAt"`
}

type RouterDec struct {
	Router *mux.Router
}

func (channel *Channel) ScanFromResult(result *sql.Rows) error {
	var parentId *int = nil

	err := result.Scan(
		&channel.ID,
		&channel.OwnerID,
		&parentId,
		&channel.Name,
		&channel.Type,
		&channel.CreatedAt,
	)
	if err != nil {
		return err
	}

	if channel.ID != 1 && parentId == nil {
		return fmt.Errorf("channel id is not 1 but the parent id is null")
	} else if channel.ID != 1 {
		channel.ParentID = *parentId
	}

	return nil
}

func (user *User) ScanFromResult(result *sql.Rows) error {
	return result.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
	)
}

func (joinedChannel *JoinedChannel) ScanFromResult(result *sql.Rows) error {
	return result.Scan(
		&joinedChannel.UserID,
		&joinedChannel.ChannelID,
		&joinedChannel.CanWrite,
		&joinedChannel.JoinedAt,
	)
}

func (message *Message) ScanFromResult(result *sql.Rows) error {
	return result.Scan(
		&message.ID,
		&message.UserID,
		&message.ChannelID,
		&message.Text,
		&message.SentAt,
	)
}

// Frontend origins
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
		w.Header().Set(
			"Access-Control-Allow-Origin", origin,
		)
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET",
		)
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
