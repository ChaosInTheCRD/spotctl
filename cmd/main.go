package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/chaosinthecrd/spotify-scraper/cli"
	urcli "github.com/urfave/cli/v2"
)

func main() {
	app := &urcli.App{
		Name:     "spotify-scraper",
		Usage:    "Scraping the spotify API for important information",
		Commands: getCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getCommands() []*urcli.Command {
	return []*urcli.Command{
		{
			Name:  "get",
			Usage: "get",
			Subcommands: []*urcli.Command{
				{
					Name: "playlists",
					Action: func(cCtx *urcli.Context) error {
						err := cli.GetPlaylists()
						if err != nil {
							return err
						}
						return nil
					},
				},
			},
		},
		{
			Name:  "status",
			Usage: "status",
			Action: func(cCtx *urcli.Context) error {
				track, err := cli.GetCurrentTrack(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
				if err != nil {
					return err
				}
				trackJSON, _ := json.Marshal(track)
				fmt.Println(string(trackJSON))
				return nil
			},
		},
		{
			Name:  "auth",
			Usage: "auth",
			Action: func(cCtx *urcli.Context) error {
				err := cli.Authenticate()
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}
