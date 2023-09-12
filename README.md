#### *Note: This implementation is not totally secure. I am sure that using refresh tokens as they are below is an insecure practice, and I am aware of other insecurities in the code. However, this is for personal use and I am happy with the compromises that I have made to make this work. Please feel free to use the code at your leisure and if you want to make it more secure, contribute back!*
#### *Another Note: This code is hodge podge of calling environment variables and command-line flags. Sorry about that, I'll neaten things up at some point (maybe)*

# spotctl
This project was created because I wanted to create a way of displaying information on [my personal website](https://chaosinthe.dev) about the music I listen to using the [Spotify Web API](https://developer.spotify.com/documentation/web-api).

## CLI
There is a simple CLI that you can use for the initial authentication flow and for making follow up requests to the Spotify Web API. I plan on improving this over time.

### Authentication
Of course you need to authenticate to the Spotify API, so the [cli](./cmd/cli) has an `auth` command that lets you get an access token after going through an authentication in the web ui. the command then stashes the refresh token (since it's the only thing we really need to repeat
 requests) on the disk to be used by other commands or the server later on.

To use the command you must first [create an "App"](https://developer.spotify.com/documentation/web-api/concepts/apps) in the [Spotify Developer Dashboard](https://developer.spotify.com/). From this you will get back a Client ID and Client Secret. These
 can then be used as follows:

```bash
$ SPOTIFY_ID=<CLIENT_ID> SPOTIFY_SECRET=<CLIENT_SECRET> go run ./cmd/main.go auth
```

### Other commands
So far, the only thing you can do with the CLI is get the current track with the `status` command. I intend on improving this over time.

## Server
There is also an [API Server](./cmd/server) implementation that I am using for my personal website as a secure way of serving [my personal website](https://chaosinthe.dev).

The server can be started by executing `go run ./server/main.go --spotify-client-id=<CLIENT_SECRET> --spotify-client-secret=<CLIENT_SECRET> --refresh-token=<REFRESH_TOKEN>`.

**NOTE: The server works by using the refresh token in `--refresh-token` at the initial startup only. From then on, it stores the refresh token in the `REFRESH_TOKEN` environment variable (yes I know it isn't the most secure way of doing things), and
 so if the server restarts, the refresh token will need to be replaced. This was my workaround as I didn't want to have to configure things so it could access some secret store and overwrite the secret each time. It is what it is.
