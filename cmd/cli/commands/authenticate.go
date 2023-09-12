package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chaosinthecrd/spotify-scraper/internal"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

var (
	ch    = make(chan *spotify.Client)
	state = internal.Generate(30)
	// These should be randomly generated for each request
	//  More information on generating these can be found here,
	// https://www.oauth.com/playground/authorization-code-with-pkce.html
	codeVerifier, err = cv.CreateCodeVerifier()
)

func Authenticate() error {
	// first start an HTTP server
	http.HandleFunc("/callback", CompleteAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := internal.Auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeVerifier.CodeChallengeS256()),
	)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	internal.Open(url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	token, err := client.Token()
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("refresh.token", []byte(token.RefreshToken), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := internal.Auth.Token(r.Context(), state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier.Value))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(internal.Auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
