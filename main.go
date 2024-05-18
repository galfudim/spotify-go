package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var spotifyClientId, spotifyClientSecret, youtubeMusicClientId, youtubeMusicClientSecret, bearerToken string

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// TODO change state to random 16-char string to prevent CSRF
	URL := "https://accounts.spotify.com/authorize"
	state := "29384dz8ag823fhh"
	grantType := "authorization_code"
	responseType := "code"
	redirectURI := "http://localhost:8000/callback"
	scope := "user-read-private+user-read-email+user-library-read+playlist-read-private"
	authURL := fmt.Sprintf("%s?client_id=%s&grant_type=%s&response_type=%s&redirect_uri=%s&scope=%s&state=%s",
		URL, spotifyClientId, grantType, responseType, redirectURI, scope, state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Missing auth code", http.StatusBadRequest)
		log.Fatal("Missing authentication code")
		return
	}

	log.Printf("Code: %s\nState: %s", code, state)

	authReqBody := url.Values{}
	authReqBody.Set("code", code)
	authReqBody.Set("redirect_uri", "http://localhost:8000/callback")
	authReqBody.Set("grant_type", "authorization_code")
	authRedirectReq, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(authReqBody.Encode()))
	if err != nil {
		log.Fatal("Failed to create authentication request", err)
	}

	secrets := fmt.Sprintf("%s:%s", spotifyClientId, spotifyClientSecret)
	encodedSecrets := base64.StdEncoding.EncodeToString([]byte(secrets))
	authorizationHeader := fmt.Sprintf("Basic %s", encodedSecrets)
	authRedirectReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authRedirectReq.Header.Set("Authorization", authorizationHeader)

	client := http.Client{}
	resp, err := client.Do(authRedirectReq)
	if err != nil {
		log.Fatal("Failed to make authentication request: ", err)
	}

	responseBody, _ := io.ReadAll(resp.Body)
	authResponse := AuthResponse{}
	err = json.Unmarshal(responseBody, &authResponse)
	if err != nil {
		log.Fatal("Method: handleRedirect, error: Unable to read json")
	}
	bearerToken = "Bearer " + fmt.Sprintf("%v", authResponse.AccessToken)
}
func handleUser(_ http.ResponseWriter, _ *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", bearerToken)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to hit \"me\" endpoint")
	}
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func handleLikedSongs(w http.ResponseWriter, _ *http.Request) {
	var savedTracks []Track
	endpoint := "https://api.spotify.com/v1/me/tracks?limit=50&offset=0"

	for {
		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req.Header.Set("Authorization", bearerToken)
		client := http.Client{}
		resp, _ := client.Do(req)

		b, _ := io.ReadAll(resp.Body)
		trackResponse := SavedTracksResponse{}
		err := json.Unmarshal(b, &trackResponse)
		if err != nil {
			log.Fatal("Method: handleLikedSongs, error: Unable to read json")
		}

		log.Printf("%d items in this iteration", len(trackResponse.Items))
		for _, item := range trackResponse.Items {
			track := item.Track
			fmt.Fprintf(w, "song: %s\n", track.Name)
			savedTracks = append(savedTracks, track)
		}

		if trackResponse.Next != "" {
			endpoint = trackResponse.Next
			log.Printf("endpoint: %s", endpoint)
		} else {
			log.Printf("%d total liked song count\n", len(savedTracks))
			return
		}
	}
}

func handlePlaylists(w http.ResponseWriter, _ *http.Request) {
	var playlists []PlaylistDTO
	endpoint := "https://api.spotify.com/v1/me/playlists?limit=50&offset=0"

	for {
		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req.Header.Set("Authorization", bearerToken)
		client := http.Client{}
		resp, _ := client.Do(req)

		b, _ := io.ReadAll(resp.Body)
		playlistResponse := SimplifiedPlaylistsResponse{}
		err := json.Unmarshal(b, &playlistResponse)
		if err != nil {
			log.Fatal("Method: handlePlaylists, error: Unable to read json")
		}

		log.Printf("%d items in this iteration", len(playlistResponse.Items))
		for _, playlist := range playlistResponse.Items {
			fmt.Fprintf(w, "playlist: %s, song count: %d, url: %s\n", playlist.Name, playlist.Tracks.Total, playlist.Tracks.HREF)
			playlists = append(playlists, playlist)
		}

		if playlistResponse.Next != "" {
			endpoint = playlistResponse.Next
			log.Printf("endpoint: %s", endpoint)
		} else {
			log.Printf("%d total playlist count\n", len(playlists))
			getPlaylistTracks(w, playlists)
			return
		}
	}
}

func getPlaylistTracks(w http.ResponseWriter, playlists []PlaylistDTO) {
	playlistToTracks := make(map[string][]string)

	for _, playlist := range playlists {
		playlistUrl := fmt.Sprintf("%s?limit=50&offset=0", playlist.Tracks.HREF)
		log.Printf("playlist name: %s, url: %s", playlist.Name, playlistUrl)

		for {
			var tracks []string
			req, _ := http.NewRequest(http.MethodGet, playlistUrl, nil)
			req.Header.Set("Authorization", bearerToken)
			client := http.Client{}
			resp, _ := client.Do(req)

			b, _ := io.ReadAll(resp.Body)
			playlistTracksResponse := PlaylistTracksResponse{}
			err := json.Unmarshal(b, &playlistTracksResponse)
			if err != nil {
				log.Fatal("Method: getPlaylistTracks, error: Unable to read json")
			}

			log.Printf("%d items in this iteration", len(playlistTracksResponse.Items))

			for _, playlistTracks := range playlistTracksResponse.Items {
				track := playlistTracks.Track
				fmt.Fprintf(w, "song: %s\n", track.Name)
				tracks = append(tracks, track.Name)
			}

			if playlistTracksResponse.Next != "" {
				playlistUrl = playlistTracksResponse.Next
				log.Printf("url: %s", playlistUrl)
			} else {
				log.Printf("%d total liked song count\n", len(tracks))
				playlistToTracks[playlist.Name] = tracks
				break
			}
		}
	}
}

func main() {
	spotifyClientId = os.Getenv(SpotifyClientId)
	spotifyClientSecret = os.Getenv(SpotifyClientSecret)
	youtubeMusicClientId = os.Getenv(YoutubeMusicClientId)
	youtubeMusicClientSecret = os.Getenv(YoutubeMusicClientSecret)

	http.HandleFunc("/", handleUser)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleRedirect)
	http.HandleFunc("/liked-songs", handleLikedSongs)
	http.HandleFunc("/playlists", handlePlaylists)

	log.Println("Server listening on localhost:8000")
	serverError := http.ListenAndServe("localhost:8000", nil)
	log.Fatalf("Server killed: %s", serverError)
}
