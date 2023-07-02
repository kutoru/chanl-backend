package endpoints

import (
	"log"
	"main/glb"
	"main/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getPrivateChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("getPrivateChannel called")

	// Getting user id

	vars := mux.Vars(r)
	userIdString, ok := vars["USER_ID"]
	if !ok {
		http.Error(w, "Could not get the user id", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "Could not convert user id to an int", http.StatusBadRequest)
		return
	}

	// Getting current channel info

	result, err := glb.DB.Query(`
		SELECT * FROM channels WHERE owner_id = ? AND type = 'pr';
	`, userId)
	if err != nil || !result.Next() {
		http.Error(w, "Could not get db result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var currentChannel models.JoinedChannel
	currentChannel.Channel = &models.Channel{}
	err = currentChannel.Channel.ScanFromResult(result)
	if err != nil {
		http.Error(w, "Could not scan channel from result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	result, err = glb.DB.Query(`
		SELECT * FROM joined_channels WHERE channel_id = ?;
	`, currentChannel.Channel.ID)
	if err != nil || !result.Next() {
		http.Error(w, "Could not get db result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = currentChannel.ScanFromResult(result)
	if err != nil {
		http.Error(w, "Could not scan channel from result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Getting channel's children info

	result, err = glb.DB.Query(`
		SELECT * FROM joined_channels WHERE user_id = ?;
	`, userId)
	if err != nil {
		http.Error(w, "Could not get db result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var childChannels []*models.JoinedChannel = make([]*models.JoinedChannel, 0)

	for result.Next() {
		var joinedChannel models.JoinedChannel
		err := joinedChannel.ScanFromResult(result)
		if err != nil {
			log.Println(err)
			continue
		}

		channelResult, err := glb.DB.Query(`
			SELECT * FROM channels WHERE id = ? AND type = 'se';
		`, joinedChannel.ChannelID)
		if err != nil {
			log.Println(err)
			continue
		} else if !channelResult.Next() {
			continue
		}

		joinedChannel.Channel = &models.Channel{}
		err = joinedChannel.Channel.ScanFromResult(channelResult)
		if err != nil {
			log.Println(err)
			continue
		}

		childChannels = append(childChannels, &joinedChannel)
	}

	// And finally constructing and sending the response

	response := struct {
		CurrentChannel *models.JoinedChannel   `json:"currentChannel"`
		ChildChannels  []*models.JoinedChannel `json:"childChannels"`
	}{
		CurrentChannel: &currentChannel,
		ChildChannels:  childChannels,
	}

	sendJSONResponse(w, response)
}
