package models

import (
	"database/sql"
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
	return result.Scan(
		&channel.ID,
		&channel.OwnerID,
		&channel.ParentID,
		&channel.Name,
		&channel.Type,
		&channel.CreatedAt,
	)
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

func (routerDec *RouterDec) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	if origin != "" {
		rw.Header().Set(
			"Access-Control-Allow-Origin", origin,
		)
		rw.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET, PUT, DELETE, PATCH",
		)
		rw.Header().Add(
			"Access-Control-Allow-Headers",
			"Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With",
		)
	}

	if req.Method == "OPTIONS" {
		return
	}

	routerDec.Router.ServeHTTP(rw, req)
}
