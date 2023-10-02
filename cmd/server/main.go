package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cli "github.com/chaosinthecrd/spotctl/cmd/cli/commands"
	internal "github.com/chaosinthecrd/spotctl/internal"
)

const (
	projectID = "743756369742"
	secretID  = "spotify"
)

var (
	clientID     = flag.String("spotify-client-id", "", "the spotify client ID")
	clientSecret = flag.String("spotify-client-secret", "", "the spotify client secret")
	currentSong  = cli.Track{}
)

func main() {
	flag.Parse()
	log.Printf("client ID set to %s", *clientID)
	log.Printf("client Secret set to %s", *clientSecret)

	refreshTrack()

	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("time is", time.Now(), "refreshing track")
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

	rt, err := internal.GetLatestSecret(projectID, secretID)
	if err != nil {
		err := errors.Join(err, fmt.Errorf("Failed to get refresh token from secrets manager"))
		log.Fatalf(err.Error())
		return
	}

	currentSong, rt, err = cli.GetCurrentTrack(*clientID, *clientSecret, rt)
	if err != nil {
		log.Printf("Error getting spotify track: %s", err.Error())
		return
	}

	if rt != "" {
		err = internal.UpdateSecret(projectID, secretID, rt)
		if err != nil {
			err := errors.Join(err, fmt.Errorf("failed to update secret"))
			log.Fatalf(err.Error())
			return
		}
	}

	log.Println("Found current song ", currentSong.Name)
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
