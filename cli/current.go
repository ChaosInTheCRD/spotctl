package cli

import (
   "fmt"
      "context"
      "github.com/chaosinthecrd/spotify-scraper/internal"
)

func GetCurrentTrack() error {
   ctx := context.Background()

   client, err := internal.GetClient()
   if err != nil {
      return err
   }

   currentlyPlaying, err := client.PlayerCurrentlyPlaying(ctx)
   if err != nil {
      return err
   }

   if !currentlyPlaying.Playing {
      fmt.Println("null")
      return nil
   }

   fmt.Println("Track: ", currentlyPlaying.Item.Name)
   fmt.Println("Artist: ", currentlyPlaying.Item.Artists[0].Name)
   fmt.Println("Preview: ", currentlyPlaying.Item.PreviewURL)
   return nil
}
