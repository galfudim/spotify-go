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

	. "spotify-go/common/http"
	. "spotify-go/common/spotify"
)

const (
	Host        = "localhost:8000"
	RedirectURI = "http://" + Host + "/callback"
)

var (
	spotifyClientId     string
	spotifyClientSecret string
	bearerToken         string
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	state := "29384dz8ag823fhh" // TODO change state to random 16-char string to prevent CSRF
	authURL := fmt.Sprintf("%s?client_id=%s&grant_type=%s&response_type=%s&redirect_uri=%s&scope=%s&state=%s",
		AuthorizeUserEndpoint,
		spotifyClientId,
		AuthorizationCodeReqQueryParam,
		CodeReqQueryParam,
		RedirectURI,
		getScopeString(),
		state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get(CodeReqQueryParam)
	state := r.URL.Query().Get(StateReqQueryParam)

	if code == "" {
		http.Error(w, "Missing auth code", http.StatusBadRequest)
		log.Fatal("Missing authentication code")
		return
	}

	log.Printf("Code: %s\nState: %s", code, state)

	authReqBody := url.Values{}
	authReqBody.Set(CodeReqBodyParam, code)
	authReqBody.Set(RedirectURIReqBodyParam, RedirectURI)
	authReqBody.Set(GrantTypeReqBodyParam, AuthorizationCodeReqQueryParam)
	authRedirectReq, err := http.NewRequest(http.MethodPost, GenerateTokenEndpoint, strings.NewReader(authReqBody.Encode()))
	if err != nil {
		log.Fatal("Failed to create authentication request", err)
	}

	secrets := fmt.Sprintf("%s:%s", spotifyClientId, spotifyClientSecret)
	encodedSecrets := base64.StdEncoding.EncodeToString([]byte(secrets))
	authorizationHeader := fmt.Sprintf("Basic %s", encodedSecrets)
	authRedirectReq.Header.Set(ContentTypeHeader, ApplicationFormUrlEncodedContentType)
	authRedirectReq.Header.Set(AuthorizationHeader, authorizationHeader)

	client := http.Client{}
	resp, err := client.Do(authRedirectReq)
	if err != nil {
		log.Fatal("Failed to make authentication request", err)
	}

	responseBody, _ := io.ReadAll(resp.Body)
	authToken := jsonToMap(string(responseBody))
	bearerToken = "Bearer " + fmt.Sprintf("%v", authToken["access_token"])
	log.Printf("Bearer token:\n%s", bearerToken)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, GetCurrentUserProfileEndpoint, nil)
	req.Header.Set(AuthorizationHeader, bearerToken)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to hit \"me\" endpoint")
	}
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func handleUserTracks(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, GetCurrentUserSavedTracksEndpoint, nil)
	req.Header.Set(AuthorizationHeader, bearerToken)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to hit \"tracks\" endpoint")
	}
	b, _ := io.ReadAll(resp.Body)
	tracks := jsonToMap(string(b))
	fmt.Fprintf(w, "tracks is %s\n", tracks)
}

func main() {
	spotifyClientId = os.Getenv(ClientId)
	spotifyClientSecret = os.Getenv(ClientSecret)

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleRedirect)
	http.HandleFunc("/", handleUser)
	http.HandleFunc("/tracks", handleUserTracks)

	log.Printf("Server listening on: %s", Host)
	serverError := http.ListenAndServe(Host, nil)
	log.Fatalf("Server killed: %s", serverError)
}
