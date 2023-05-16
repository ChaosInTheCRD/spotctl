package cli

import (
	"context"
	"fmt"

	"github.com/chaosinthecrd/spotify-scraper/internal"
)

type Track struct {
	Name       string `json:"name"`
	Artist     string `json:"artist"`
	PreviewURL string `json:"previewURL"`
}

func GetCurrentTrack(clientID, clientSecret string) (Track, error) {
	ctx := context.Background()

	client, err := internal.GetClient(clientID, clientSecret)
	if err != nil {
		return Track{}, err
	}

	currentlyPlaying, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return Track{}, err
	}

	if !currentlyPlaying.Playing {
		fmt.Println("null")
		return Track{}, err
	}

	return Track{Name: currentlyPlaying.Item.Name, Artist: currentlyPlaying.Item.Artists[0].Name, PreviewURL: currentlyPlaying.Item.PreviewURL}, nil
}
