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

var (
	spotifyClientId     string
	spotifyClientSecret string
	bearerToken         string
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// TODO change state to random 16-char string to prevent CSRF
	state := "29384dz8ag823fhh"
	authURL := fmt.Sprintf("%s?client_id=%s&grant_type=%s&response_type=%s&redirect_uri=%s&scope=%s&state=%s",
		AuthorizeUserEndpoint,
		spotifyClientId,
		"authorization_code",
		"code",
		"http://localhost:8000/callback",
		"user-read-private+user-read-email+user-library-read+playlist-read-private",
		state)
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
	authRedirectReq, err := http.NewRequest(http.MethodPost, GenerateTokenEndpoint, strings.NewReader(authReqBody.Encode()))
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
		log.Fatal("Failed to make authentication request", err)
	}

	responseBody, _ := io.ReadAll(resp.Body)
	authToken := jsonToMap(responseBody)
	bearerToken = "Bearer " + fmt.Sprintf("%v", authToken["access_token"])
	//log.Printf("Bearer token:\n%s", bearerToken)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, GetCurrentUserProfileEndpoint, nil)
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
	var likedSongs []Track
	endpoint := "https://api.spotify.com/v1/me/tracks?limit=50&offset=0"

	for {
		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req.Header.Set("Authorization", bearerToken)
		client := http.Client{}
		resp, _ := client.Do(req)

		b, _ := io.ReadAll(resp.Body)
		trackResponse := TrackResponse{}
		err := json.Unmarshal(b, &trackResponse)
		if err != nil {
			log.Fatal("Method: handleLikedSongs, error: Unable to read json")
		}

		log.Printf("%d items in this iteration", len(trackResponse.Items))
		for _, item := range trackResponse.Items {
			track := item.Track
			fmt.Fprintf(w, "song: %s\n", track.Name)
			likedSongs = append(likedSongs, track)
		}

		if trackResponse.Next != "" {
			endpoint = trackResponse.Next
			log.Printf("endpoint: %s", endpoint)
		} else {
			log.Printf("%d total liked song count\n", len(likedSongs))
			return
		}
	}
}

// handlePlaylists first gets all playlist IDs
func handlePlaylists(w http.ResponseWriter, _ *http.Request) {
	var playlistRefs []any
	endpoint := "https://api.spotify.com/v1/me/playlists?limit=50&offset=0"

	for {
		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req.Header.Set("Authorization", bearerToken)
		client := http.Client{}
		resp, _ := client.Do(req)

		b, _ := io.ReadAll(resp.Body)
		playlistResponse := PlaylistResponse{}
		err := json.Unmarshal(b, &playlistResponse)
		if err != nil {
			log.Fatal("Method: handlePlaylists, error: Unable to read json")
		}

		log.Printf("%d items in this iteration", len(playlistResponse.Items))
		for _, playlist := range playlistResponse.Items {
			fmt.Fprintf(w, "playlist: %s, song count: %d\n", playlist.Name, playlist.Tracks.Total)
			playlistRefs = append(playlistRefs, playlist.Tracks.HREF)
		}

		if playlistResponse.Next != "" {
			endpoint = playlistResponse.Next
			log.Printf("endpoint: %s", endpoint)
		} else {
			log.Printf("%d total playlist count\n", len(playlistRefs))
			return
		}
	}
}

func main() {
	spotifyClientId = os.Getenv(ClientId)
	spotifyClientSecret = os.Getenv(ClientSecret)

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleRedirect)
	http.HandleFunc("/", handleUser)
	http.HandleFunc("/liked-songs", handleLikedSongs)
	http.HandleFunc("/playlists", handlePlaylists)

	log.Println("Server listening on localhost:8000")
	serverError := http.ListenAndServe("localhost:8000", nil)
	log.Fatalf("Server killed: %s", serverError)
}
