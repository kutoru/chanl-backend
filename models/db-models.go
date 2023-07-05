package models

import (
	"database/sql"
	"fmt"
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
	ID        int     `json:"id"`
	UserID    int     `json:"userId"`
	ChannelID int     `json:"channelId"`
	Text      string  `json:"text"`
	SentAt    string  `json:"sentAt"`
	UserName  *string `json:"userName"`
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

func (message *Message) ScanFromResult(result *sql.Rows, scanUserName bool) error {
	if scanUserName {
		return result.Scan(
			&message.ID,
			&message.UserID,
			&message.ChannelID,
			&message.Text,
			&message.SentAt,
			&message.UserName,
		)
	} else {
		return result.Scan(
			&message.ID,
			&message.UserID,
			&message.ChannelID,
			&message.Text,
			&message.SentAt,
		)
	}
}
