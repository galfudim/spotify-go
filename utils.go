package main

import "encoding/json"

const (
	ClientId     = "SPOTIFY_CLIENT_ID"
	ClientSecret = "SPOTIFY_CLIENT_SECRET"
)

const (
	AuthorizeUserEndpoint         = "https://accounts.spotify.com/authorize"
	GenerateTokenEndpoint         = "https://accounts.spotify.com/api/token"
	GetCurrentUserProfileEndpoint = "https://api.spotify.com/v1/me"
)

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil
	}
	return result
}
