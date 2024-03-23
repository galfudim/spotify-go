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

var bearerToken string

func handleRoot(w http.ResponseWriter, r *http.Request) {
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	URL := "https://accounts.spotify.com/authorize"
	authURL := fmt.Sprintf("%s?client_id=%s&grant_type=%s&response_type=%s&redirect_uri=%s&scope=%s&state=%s", URL, clientId, GrantType, ResponseType, RedirectURI, Scope, State)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code, state := r.URL.Query().Get("code"), r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Missing auth code", http.StatusBadRequest)
		log.Fatal("Missing authentication code")
		return
	}

	log.Printf("Code: %s\nState: %s", code, state)

	clientId, clientSecret := os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET")
	authReqBody := url.Values{}
	authReqBody.Set("code", code)
	authReqBody.Set("redirect_uri", RedirectURI)
	authReqBody.Set("grant_type", GrantType)
	authRedirectReq, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(authReqBody.Encode()))
	if err != nil {
		log.Fatal("Failed to create authentication request", err)
	}

	secrets := clientId + ":" + clientSecret
	encodedSecrets := base64.StdEncoding.EncodeToString([]byte(secrets))
	authorizationHeader := fmt.Sprintf("Basic %s", encodedSecrets)
	authRedirectReq.Header.Set("Content-Type", ContentType)
	authRedirectReq.Header.Set("Authorization", authorizationHeader)

	client := http.Client{}
	resp, err := client.Do(authRedirectReq)
	if err != nil {
		log.Fatal("Failed to make authentication request", err)
	}

	responseBody, _ := io.ReadAll(resp.Body)
	log.Println(string(responseBody))
	authToken := jsonToMap(string(responseBody))
	bearerToken = "Bearer " + fmt.Sprintf("%v", authToken["access_token"])
	log.Printf("Bearer token:\n%s", bearerToken)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", bearerToken)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to make \"me\" request")
	}
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/callback", handleRedirect)
	http.HandleFunc("/me", handleUser)

	log.Println("Server listening on localhost:8000")
	serverError := http.ListenAndServe(":8000", nil)
	log.Fatalf("Server killed, error: %s", serverError)
}
