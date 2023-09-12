package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cli "github.com/chaosinthecrd/spotctl/cmd/cli/commands"
)

var (
	refreshToken = flag.String("refresh-token", "", "the refresh token to use at boot")
	clientID     = flag.String("spotify-client-id", "", "the spotify client ID")
	clientSecret = flag.String("spotify-client-secret", "", "the spotify client secret")
	currentSong  = cli.Track{}
)

func main() {
	flag.Parse()
	log.Printf("refresh token path set to %s", *refreshToken)
	log.Printf("client ID set to %s", *clientID)
	log.Printf("client Secret set to %s", *clientSecret)

	refreshTrack()

	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				refreshTrack()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	http.HandleFunc("/status", statusHandler)
	log.Println(http.ListenAndServe(":8080", nil))
	close(quit)
	os.Exit(1)
}

var CurrentSong string

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func refreshTrack() {
	var err error
	if *clientID == "" || *clientSecret == "" {
		log.Printf("client ID and client secret not set correctly. Exiting...")
		return
	}

	rt := os.Getenv("REFRESH_TOKEN")
	if rt == "" && *refreshToken != "" {
		fmt.Println("using initial refresh token")
		rt = *refreshToken
	} else {
		log.Printf("refreshToken not set correctly. Exiting...")
	}

	currentSong, err = cli.GetCurrentTrack(*clientID, *clientSecret, rt)
	if err != nil {
		log.Printf("Error getting spotify track: %s", err.Error())
		return
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Serialize the response data to JSON
	err := json.NewEncoder(w).Encode(currentSong)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
