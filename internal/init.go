package internal

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const RedirectURI = "http://localhost:8080/callback"
const checkURL = "https://api.spotify.com/v1"

var Auth = spotifyauth.New(spotifyauth.WithRedirectURL(RedirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopePlaylistReadPrivate))

func GetClient(clientID, clientSecret, refreshToken string) (*spotify.Client, string, error) {

	// this is dirty, but I want to check if I am getting a 503... might work, might not.
	c := &http.Client{}

	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		fmt.Println("Got a 503 Service Unavailable error. Returning...")
		return nil, refreshToken, nil
	}
	auth := spotifyauth.New(spotifyauth.WithRedirectURL(RedirectURI), spotifyauth.WithClientID(clientID), spotifyauth.WithClientSecret(clientSecret), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopePlaylistReadPrivate))
	ctx := context.TODO()

	token := new(oauth2.Token)
	token.Expiry = time.Now().Add(time.Second * -5)

	token.RefreshToken = refreshToken

	fmt.Println(token.RefreshToken)

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(ctx, token))

	// Need to set the new refresh token for the next request
	newToken, err := client.Token()
	if err != nil {
		return nil, "", err
	}

	return client, newToken.RefreshToken, nil
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
