package commands

import (
	"context"
	"fmt"

	"github.com/chaosinthecrd/spotctl/internal"
)

type Track struct {
	Name       string `json:"name"`
	Artist     string `json:"artist"`
	PreviewURL string `json:"previewURL"`
}

func GetCurrentTrack(clientID, clientSecret, refreshToken string) (Track, string, error) {
	ctx := context.Background()

	client, rt, err := internal.GetClient(clientID, clientSecret, refreshToken)
	if err != nil {
		return Track{}, "", err
	} else if client == nil {
		fmt.Println("Got a 503... skipping this for now")
		return Track{}, rt, nil
	}

	currentlyPlaying, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return Track{}, rt, err
	}

	if !currentlyPlaying.Playing || currentlyPlaying.Item == nil {
		fmt.Println("null")
		return Track{}, rt, nil
	}

	if len(currentlyPlaying.Item.Artists) > 0 {
		return Track{Name: currentlyPlaying.Item.Name, Artist: currentlyPlaying.Item.Artists[0].Name, PreviewURL: currentlyPlaying.Item.PreviewURL}, rt, nil
	} else {
		return Track{Name: currentlyPlaying.Item.Name, Artist: "", PreviewURL: currentlyPlaying.Item.PreviewURL}, rt, nil
	}
}
