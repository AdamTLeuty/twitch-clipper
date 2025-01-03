package main

import (
	"fmt"
	"net/http"
)

func streamPickerHandler(w http.ResponseWriter, r *http.Request) {
	//Reads the url of the Twitch stream the user picked
	if r.Method == http.MethodPost {
		// Retrieve the form value
		userInput := r.FormValue("channel")
		// Respond with the user's input
		fmt.Fprintf(w, "You entered: %s", userInput)

		channelName, err := StreamlinkChannelName(userInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Channel Name: %s\n", channelName)
		}

		go RecordTwitchStream(userInput)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
