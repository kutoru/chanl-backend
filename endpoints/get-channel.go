package endpoints

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kutoru/chanl-backend/glb"
	"github.com/kutoru/chanl-backend/models"
)

func getChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("getChannel called")

	// Getting and checking channel id

	channelId, err := getMuxVar(r, "CHANNEL_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Checking user id. TODO: do some kind of authorization check instead

	userIdString := r.Header.Get("User-ID")
	userId := 0
	if userIdString != "" {
		userId, err = strconv.Atoi(userIdString)
		if err != nil {
			http.Error(w, "Could not convert user id to an int", http.StatusBadRequest)
			return
		}
	}

	// Getting current channel

	currentChannel, err := getCurrentChannel(userId, channelId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting parent channel

	parentChannel, err := getParentChannel(userId, currentChannel.Channel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting child channels

	childChannels, err := getChildChannels(userId, currentChannel.Channel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// And finally constructing and sending the response

	response := struct {
		ParentChannel  *models.JoinedChannel   `json:"parentChannel"`
		CurrentChannel *models.JoinedChannel   `json:"currentChannel"`
		ChildChannels  []*models.JoinedChannel `json:"childChannels"`
	}{
		ParentChannel:  parentChannel,
		CurrentChannel: currentChannel,
		ChildChannels:  childChannels,
	}

	sendJSONResponse(w, response)
}

func getCurrentChannel(userId int, channelId int) (*models.JoinedChannel, error) {
	result, err := glb.DB.Query(`
		SELECT * FROM channels WHERE id = ?;
	`, channelId)
	if err != nil || !result.Next() {
		return nil, fmt.Errorf("could not get db result for currentChannel. %v", err)
	}

	var currentChannel *models.JoinedChannel = &models.JoinedChannel{}
	currentChannel.Channel = &models.Channel{}

	err = currentChannel.Channel.ScanFromResult(result)
	result.Close()
	if err != nil {
		return nil, fmt.Errorf("could not scan channel from result. %v", err)
	}

	// If user id is 0, the user is not logged in. In that case just returning the channel info
	if userId <= 0 {
		return currentChannel, nil
	}

	result, err = glb.DB.Query(`
		SELECT * FROM joined_channels WHERE user_id = ? AND channel_id = ?;
	`, userId, channelId)
	if err != nil || !result.Next() {
		return nil, fmt.Errorf("could not get db result for joinedChannel. %v", err)
	}

	err = currentChannel.ScanFromResult(result)
	result.Close()
	if err != nil {
		return nil, fmt.Errorf("could not scan joinedChannel from result. %v", err)
	}

	return currentChannel, nil
}

func getParentChannel(userId int, channel *models.Channel) (*models.JoinedChannel, error) {
	if userId <= 0 {
		return nil, nil
	}

	var parentChannelType models.ChannelType

	switch channel.Type {
	case models.Global: // global doesn't have a parent
		return nil, nil
	case models.Private:
		parentChannelType = models.Global
	case models.Server:
		parentChannelType = models.Private
	case models.Room:
		parentChannelType = models.Server
	case models.Personal:
		parentChannelType = models.Global
	case models.Friend:
		parentChannelType = models.Personal
	default:
		return nil, fmt.Errorf("unknown channel type (%v)", channel.Type)
	}

	parentChannel := &models.JoinedChannel{}
	parentChannel.Channel = &models.Channel{}

	var result *sql.Rows = &sql.Rows{}
	var err error = nil

	if parentChannelType == models.Global {
		result, err = glb.DB.Query(`
			SELECT * FROM joined_channels WHERE user_id = ? AND channel_id = 1;
		`, userId)
	} else if parentChannelType == models.Private || parentChannelType == models.Personal {
		result, err = glb.DB.Query(`
			SELECT joined_channels.* FROM joined_channels
			INNER JOIN channels ON
			(channels.owner_id = ?) AND
			(channels.type = ?) AND
			(joined_channels.channel_id = channels.id);
		`, userId, parentChannelType)
	} else {
		result, err = glb.DB.Query(`
			SELECT * FROM joined_channels WHERE user_id = ? AND channel_id = ?;
		`, userId, channel.ParentID)
	}

	if err != nil || !result.Next() {
		return nil, fmt.Errorf("could not get db result for parentChannel. %v", err)
	}

	err = parentChannel.ScanFromResult(result)
	result.Close()
	if err != nil {
		return nil, fmt.Errorf("could not scan parentChannel from result. %v", err)
	}

	result, err = glb.DB.Query(`
		SELECT * FROM channels WHERE id = ?;
	`, parentChannel.ChannelID)
	if err != nil || !result.Next() {
		return nil, fmt.Errorf("could not get db result for parentChannel.Channel. %v", err)
	}

	err = parentChannel.Channel.ScanFromResult(result)
	result.Close()
	if err != nil {
		return nil, fmt.Errorf("could not scan parentChannel.Channel from result. %v", err)
	}

	return parentChannel, nil
}

func getChildChannels(userId int, channel *models.Channel) ([]*models.JoinedChannel, error) {
	var childChannels []*models.JoinedChannel = make([]*models.JoinedChannel, 0)
	if userId <= 0 {
		return childChannels, nil
	}

	var childChannelType models.ChannelType

	switch channel.Type {
	case models.Global: // Global has two types of children, it has its own special case later
		childChannelType = ""
	case models.Private:
		childChannelType = models.Server
	case models.Server:
		childChannelType = models.Room
	case models.Room:
		return childChannels, nil
	case models.Personal:
		childChannelType = models.Friend
	case models.Friend:
		return childChannels, nil
	default:
		return nil, fmt.Errorf("unknown channel type (%v)", channel.Type)
	}

	var result *sql.Rows = &sql.Rows{}
	var err error = nil
	if channel.Type == models.Global {
		result, err = glb.DB.Query(`
			SELECT joined_channels.* FROM joined_channels
			INNER JOIN channels ON
			(joined_channels.user_id = ?) AND
			(joined_channels.channel_id = channels.id) AND
			(channels.parent_id = 1) AND
			(channels.owner_id = ?);
		`, userId, userId)
	} else {
		result, err = glb.DB.Query(`
			SELECT joined_channels.* FROM joined_channels
			INNER JOIN channels ON
			(joined_channels.user_id = ?) AND
			(joined_channels.channel_id = channels.id) AND
			(channels.type = ?);
		`, userId, childChannelType)
	}

	if err != nil {
		return nil, fmt.Errorf("could not get db result for childChannels. %v", err)
	}

	for result.Next() {
		var joinedChannel models.JoinedChannel
		joinedChannel.Channel = &models.Channel{}

		err := joinedChannel.ScanFromResult(result)
		if err != nil {
			return nil, fmt.Errorf("could not parse child joinedChannel from result. %v", err)
		}

		channelResult, err := glb.DB.Query(`
			SELECT * FROM channels WHERE id = ?;
		`, joinedChannel.ChannelID)
		if err != nil || !channelResult.Next() {
			return nil, fmt.Errorf("could not get child joinedChannel.Channel result from db. The channel id exists in joined_channels but doesn't exist in channels?. %v", err)
		}

		err = joinedChannel.Channel.ScanFromResult(channelResult)
		channelResult.Close()
		if err != nil {
			return nil, fmt.Errorf("could not parse child joinedChannel.Channel from result. %v", err)
		}

		childChannels = append(childChannels, &joinedChannel)
	}

	result.Close()

	return childChannels, nil
}
