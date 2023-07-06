package endpoints

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kutoru/chanl-backend/glb"
	"github.com/kutoru/chanl-backend/models"
)

type ChannelResponse struct {
	ParentChannel  *models.JoinedChannel   `json:"parentChannel"`
	CurrentChannel *models.JoinedChannel   `json:"currentChannel"`
	ChildChannels  []*models.JoinedChannel `json:"childChannels"`
}

func getChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("getChannel called")

	// Getting and checking channel id

	vars := mux.Vars(r)
	channelIdString, ok := vars["CHANNEL_ID"]
	if !ok {
		http.Error(w, "Could not get the channel id", http.StatusBadRequest)
		return
	}

	channelId, err := strconv.Atoi(channelIdString)
	if err != nil {
		http.Error(w, "Could not convert channel id to an int", http.StatusBadRequest)
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

	// Getting channel info

	var currentChannel *models.JoinedChannel = &models.JoinedChannel{}
	currentChannel.Channel = &models.Channel{}

	result, err := glb.DB.Query(`
		SELECT * FROM channels WHERE id = ?;
	`, channelId)
	if err != nil || !result.Next() {
		http.Error(w, "Could not get db result for channel", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = currentChannel.Channel.ScanFromResult(result)
	result.Close()
	if err != nil {
		http.Error(w, "Could not scan channel from result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// If user id is 0, the user is not logged in. In that case just returning the channel info

	if userId == 0 {
		response := ChannelResponse{
			ParentChannel:  nil,
			CurrentChannel: currentChannel,
			ChildChannels:  []*models.JoinedChannel{},
		}

		sendJSONResponse(w, response)
		return
	}

	// Otherwise getting parentChannel, currentChannel info and childChannels

	var parentChannel *models.JoinedChannel = nil

	// Have to do an extra check for parentChannel because global channel doesn't have a parent
	if currentChannel.Channel.ID != 1 {
		parentChannel = &models.JoinedChannel{}
		parentChannel.Channel = &models.Channel{}

		result, err = glb.DB.Query(`
			SELECT * FROM joined_channels WHERE user_id = ? AND channel_id = ?;
		`, userId, currentChannel.Channel.ParentID)
		if err != nil || !result.Next() {
			http.Error(w, "Could not get db result for parentChannel", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		err = parentChannel.ScanFromResult(result)
		result.Close()
		if err != nil {
			http.Error(w, "Could not scan parentChannel from result", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		result, err = glb.DB.Query(`
			SELECT * FROM channels WHERE id = ?;
		`, parentChannel.ChannelID)
		if err != nil || !result.Next() {
			http.Error(w, "Could not get db result for parentChannel.Channel", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		err = parentChannel.Channel.ScanFromResult(result)
		result.Close()
		if err != nil {
			http.Error(w, "Could not scan parentChannel.Channel from result", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	// Getting currentChannel info

	result, err = glb.DB.Query(`
		SELECT * FROM joined_channels WHERE user_id = ? AND channel_id = ?;
	`, userId, channelId)
	if err != nil || !result.Next() {
		http.Error(w, "Could not get db result for joinedChannel", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = currentChannel.ScanFromResult(result)
	result.Close()
	if err != nil {
		http.Error(w, "Could not scan joinedChannel from result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Getting channel's children info

	var childChannels []*models.JoinedChannel = make([]*models.JoinedChannel, 0)

	result, err = glb.DB.Query(`
		SELECT joined_channels.* FROM joined_channels
		INNER JOIN channels ON
		(joined_channels.user_id = ?) AND
		(joined_channels.channel_id = channels.id) AND
		(channels.parent_id = ?);
	`, userId, channelId)
	if err != nil {
		http.Error(w, "Could not get db result for childChannels", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for result.Next() {
		var joinedChannel models.JoinedChannel
		joinedChannel.Channel = &models.Channel{}

		err := joinedChannel.ScanFromResult(result)
		if err != nil {
			http.Error(w, "Could not parse child joinedChannel from result", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		channelResult, err := glb.DB.Query(`
			SELECT * FROM channels WHERE id = ?;
		`, joinedChannel.ChannelID)
		if err != nil || !channelResult.Next() {
			http.Error(w, "Could not get child joinedChannel.Channel result from db. The channel id exists in joined_channels but doesn't exist in channels?", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		err = joinedChannel.Channel.ScanFromResult(channelResult)
		channelResult.Close()
		if err != nil {
			http.Error(w, "Could not parse child joinedChannel.Channel from result", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		childChannels = append(childChannels, &joinedChannel)
	}

	result.Close()

	// And finally constructing and sending the response

	response := ChannelResponse{
		ParentChannel:  parentChannel,
		CurrentChannel: currentChannel,
		ChildChannels:  childChannels,
	}

	sendJSONResponse(w, response)
}
