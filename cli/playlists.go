package cli

import (
	"context"
	"fmt"
	"github.com/chaosinthecrd/spotify-scraper/internal"
)

func GetPlaylists() error {
   ctx := context.Background()

   client, err := internal.GetClient()
   if err != nil {
      return err
   }

   playlists, err := client.CurrentUsersPlaylists(ctx)
   if err != nil {
      return err
   }

   for i := range playlists.Playlists {
      fmt.Printf("%s:\n", playlists.Playlists[i].Name)

      playlist, err := client.GetPlaylistTracks(ctx, playlists.Playlists[i].ID)
      if err != nil {
         return err
      }

      tracks := playlist.Tracks

      for i := range tracks {
       fmt.Printf("      %s (%s) \n", tracks[i].Track.Name, tracks[i].Track.PreviewURL)
      }
   }
   return nil
}
