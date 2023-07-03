package endpoints

import (
	"log"
	"main/glb"
	"main/models"
	"net/http"
	"strconv"
)

func getGlobalChannel(w http.ResponseWriter, r *http.Request) {
	log.Println("getGlobalChannel called")

	// Checking user id. TODO: do some kind of authorization check instead

	userIdString := r.Header.Get("User-ID")
	userId := 0
	var err error
	if len(userIdString) > 0 {
		userId, err = strconv.Atoi(userIdString)
		if err != nil {
			http.Error(w, "Could not convert user id to an int", http.StatusBadRequest)
			return
		}
	}

	// Getting global channel info

	result, err := glb.DB.Query(`
		SELECT * FROM channels WHERE id = 1;
	`)
	if err != nil || !result.Next() {
		http.Error(w, "Could not get db result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var currentChannel models.JoinedChannel
	currentChannel.Channel = &models.Channel{}
	err = currentChannel.Channel.ScanFromResult(result)
	result.Close()
	if err != nil {
		http.Error(w, "Could not scan channel from result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// If user id is 0, the user is not logged in. In that case just returning the channel info

	if userId == 0 {
		response := struct {
			CurrentChannel *models.JoinedChannel   `json:"currentChannel"`
			ChildChannels  []*models.JoinedChannel `json:"childChannels"`
		}{
			CurrentChannel: &currentChannel,
			ChildChannels:  []*models.JoinedChannel{},
		}

		sendJSONResponse(w, response)
		return
	}

	// Otherwise getting the currentChannel info and the childChannels

	result, err = glb.DB.Query(`
		SELECT * FROM joined_channels WHERE channel_id = 1 AND user_id = ?;
	`, userId)
	if err != nil || !result.Next() {
		http.Error(w, "Could not get db result", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = currentChannel.ScanFromResult(result)
	result.Close()
	if err != nil {
		http.Error(w, "Could not scan channel from result", http.StatusInternalServerError)
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
		(channels.parent_id = 1);
	`, userId)
	if err != nil {
		http.Error(w, "Could not get db result for childChannels", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for result.Next() {
		var joinedChannel models.JoinedChannel
		err := joinedChannel.ScanFromResult(result)
		if err != nil {
			http.Error(w, "Could not parse joinedChannel.Channel from result", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		channelResult, err := glb.DB.Query(`
			SELECT * FROM channels WHERE id = ?;
		`, joinedChannel.ChannelID)
		if err != nil || !channelResult.Next() {
			http.Error(w, "Could not get channel from db by joinedChannel.ChannelID", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		joinedChannel.Channel = &models.Channel{}
		err = joinedChannel.Channel.ScanFromResult(channelResult)
		channelResult.Close()
		if err != nil {
			http.Error(w, "Could not parse joinedChannel.Channel from result", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		childChannels = append(childChannels, &joinedChannel)
	}

	result.Close()

	// For each user there should be exactly 2 channels that have parentId 1. These channels are Private and Personal
	// Checking if the array is correct (technically the database tables and the queries should be safe enough, but just in case i also check here)

	if len(childChannels) != 2 ||
		!(childChannels[0].Channel.Type == models.Private || childChannels[1].Channel.Type == models.Private) ||
		!(childChannels[0].Channel.Type == models.Personal || childChannels[1].Channel.Type == models.Personal) {
		http.Error(w, "Invalid childChannels", http.StatusInternalServerError)
		log.Printf("len: %v", len(childChannels))
		return
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
