package internal

import (
	"context"
	"os"
        "fmt"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const RedirectURI = "http://localhost:8080/callback"

var (
	Auth  = spotifyauth.New(spotifyauth.WithRedirectURL(RedirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopePlaylistReadPrivate))
)

func GetClient() (*spotify.Client, error) {
   ctx := context.TODO()
  
   fmt.Println("Reading Token")
   t := os.Getenv("REFRESH_TOKEN")

   token := new(oauth2.Token)
   token.Expiry = time.Now().Add(time.Second * -5)

   fmt.Println(string(t))
   token.RefreshToken = string(t)

   // use the token to get an authenticated client
   client := spotify.New(Auth.Client(ctx, token))

   // Need to set the new refresh token for the next request
   newToken, err := client.Token()
   if err != nil {
      return nil, err
   }

   fmt.Println("Writing Token")
   os.Setenv("REFRESH_TOKEN", newToken.RefreshToken)

   return client, nil
}

// open opens the specified URL in the default browser of the user.
func Open(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}
