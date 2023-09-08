package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"

	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const RedirectURI = "http://localhost:8080/callback"

var Auth = spotifyauth.New(spotifyauth.WithRedirectURL(RedirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopePlaylistReadPrivate))

func GetClient(clientID, clientSecret, refreshToken string) (*spotify.Client, error) {
	auth := spotifyauth.New(spotifyauth.WithRedirectURL(RedirectURI), spotifyauth.WithClientID(clientID), spotifyauth.WithClientSecret(clientSecret), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopePlaylistReadPrivate))
	ctx := context.TODO()

	tByte, err := os.ReadFile(refreshToken)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("cannot read refresh token file"))
	}

	t := string(tByte)

	token := new(oauth2.Token)
	token.Expiry = time.Now().Add(time.Second * -5)

	token.RefreshToken = string(t)

	fmt.Println(token.RefreshToken)

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(ctx, token))

	// Need to set the new refresh token for the next request
	newToken, err := client.Token()
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(refreshToken, []byte(newToken.RefreshToken), 0644)
	if err != nil {
		return nil, err
	}
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
