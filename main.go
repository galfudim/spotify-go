package main

import (
	"encoding/base64"
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
		"user-read-private+user-read-email+user-library-read",
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
	authToken := jsonToMap(string(responseBody))
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

func handleUserTracks(w http.ResponseWriter, r *http.Request) {
	// arr to collect values
	arr := []interface{}{}
	endpoint := "https://api.spotify.com/v1/me/tracks?limit=50&offset=0"

	for {
		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req.Header.Set("Authorization", bearerToken)
		client := http.Client{}
		resp, _ := client.Do(req)

		b, _ := io.ReadAll(resp.Body)
		tracks := jsonToMap(string(b))
		for key, value := range tracks {
			// TODO ad items to array
			if key == "items" {
				values := value.([]interface{}) // arr of tracks with added date and track data
				for _, track := range values {
					trackData := track.(map[string]interface{})["track"]
					song := trackData.(map[string]interface{})
					fmt.Fprintf(w, "song: %s\n", song["name"])
					arr = append(arr, song)
				}
			} else if key == "next" {
				if value != nil {
					endpoint = fmt.Sprintf("%s", value.(string))
					log.Printf("endpoint: %s", endpoint)
				} else {
					log.Printf("%d total items collected\n", len(arr))
					return
				}
			}
		}
	}
}

func main() {
	spotifyClientId = os.Getenv(ClientId)
	spotifyClientSecret = os.Getenv(ClientSecret)

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleRedirect)
	http.HandleFunc("/", handleUser)
	http.HandleFunc("/tracks", handleUserTracks)

	log.Println("Server listening on localhost:8000")
	serverError := http.ListenAndServe("localhost:8000", nil)
	log.Fatalf("Server killed: %s", serverError)
}
