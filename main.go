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

	. "spotify-go/constants/http"
	. "spotify-go/constants/spotify"
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	URL := "https://accounts.spotify.com/authorize"
	state := "29384dz8ag823fhh" // TODO change state to random 16-char string to prevent CSRF
	authURL := fmt.Sprintf("%s?client_id=%s&grant_type=%s&response_type=%s&redirect_uri=%s&scope=%s&state=%s",
		URL,
		spotifyClientId,
		ReqQueryParamAuthorizationCode,
		ReqQueryParamCode,
		RedirectURI,
		fmt.Sprintf("%s+%s", ScopeReadPrivate, ScopeReadEmail),
		state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get(ReqQueryParamCode)
	state := r.URL.Query().Get(ReqQueryParamState)

	if code == "" {
		http.Error(w, "Missing auth code", http.StatusBadRequest)
		log.Fatal("Missing authentication code")
		return
	}

	log.Printf("Code: %s\nState: %s", code, state)

	authReqBody := url.Values{}
	authReqBody.Set(ReqBodyParamCode, code)
	authReqBody.Set(ReqBodyParamRedirectURI, RedirectURI)
	authReqBody.Set(ReqBodyParamGrantType, ReqQueryParamAuthorizationCode)
	authRedirectReq, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(authReqBody.Encode()))
	if err != nil {
		log.Fatal("Failed to create authentication request", err)
	}

	secrets := fmt.Sprintf("%s:%s", spotifyClientId, spotifyClientSecret)
	encodedSecrets := base64.StdEncoding.EncodeToString([]byte(secrets))
	authorizationHeader := fmt.Sprintf("Basic %s", encodedSecrets)
	authRedirectReq.Header.Set(HeaderContentType, ContentTypeApplicationFormUrlEncoded)
	authRedirectReq.Header.Set(HeaderAuthorization, authorizationHeader)

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
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	req.Header.Set(HeaderAuthorization, bearerToken)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to make \"me\" request")
	}
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func main() {
	spotifyClientId = os.Getenv(ClientId)
	spotifyClientSecret = os.Getenv(ClientSecret)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/callback", handleRedirect)
	http.HandleFunc("/me", handleUser)

	log.Printf("Server listening on %s\n", Host)
	serverError := http.ListenAndServe(Host, nil)
	log.Fatalf("Server killed, error: %s", serverError)
}
