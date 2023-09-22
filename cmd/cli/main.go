package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	cmd "github.com/chaosinthecrd/spotctl/cmd/cli/commands"
	urcli "github.com/urfave/cli/v2"
)

func main() {
	app := &urcli.App{
		Name:     "spotctl",
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
						err := cmd.GetPlaylists()
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
				rt, err := os.ReadFile("./refresh.token")
				if err != nil {
					log.Fatalf(errors.Join(err, fmt.Errorf("could not read refresh token file")).Error())
					return err
				}
				track, nrt, err := cmd.GetCurrentTrack(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"), string(rt))
				if err != nil {
					log.Fatalf(errors.Join(err, fmt.Errorf("failed to get current track")).Error())
					return err
				}

				err = os.WriteFile("./refresh.token", []byte(nrt), 0644)
				if err != nil {
					log.Fatalf(errors.Join(err, fmt.Errorf("failed to get current track")).Error())
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
				err := cmd.Authenticate()
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}
