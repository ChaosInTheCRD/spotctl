package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/chaosinthecrd/spotify-scraper/cli"
)

const (
	defaultRefreshTokenPath = "/tmp/spotify/refresh.token"
)

var (
	refreshTokenPath = flag.String("refresh-token-path", defaultRefreshTokenPath, "the path to the file containing the reresh token")
	clientID         = flag.String("spotify-client-id", "", "the spotify client ID")
	clientSecret     = flag.String("spotify-client-secret", "", "the spotify client secret")
)

func main() {
	http.HandleFunc("/status", statusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	flag.Parse()
	if *clientID == "" || *clientSecret == "" {
		log.Printf("client ID and client secret not set correctly. Exiting...")
		return
	}
	response, err := cli.GetCurrentTrack(*clientID, *clientSecret)
	if err != nil {
		log.Printf("Error getting spotify track: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Serialize the response data to JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
