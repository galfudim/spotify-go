package main

import "encoding/json"

type Scope string

const (
	UserReadPrivateScope Scope = "user-read-private"
	UserEmailReadScope   Scope = "user-read-email"

	UserLibraryReadScope Scope = "user-library-read"
)

const (
	ClientId     = "SPOTIFY_CLIENT_ID"
	ClientSecret = "SPOTIFY_CLIENT_SECRET"
)

const (
	SpotifyAccountEndpoint = "https://accounts.spotify.com"
	SpotifyApiEndpoint     = "https://api.spotify.com"
	v1ApiEndpoint          = SpotifyApiEndpoint + "/v1"
	AuthorizeUserEndpoint  = SpotifyAccountEndpoint + "/authorize"
	GenerateTokenEndpoint  = SpotifyAccountEndpoint + "/api/token"
)

// Resources
const (
	CurrentUserEndpoint = v1ApiEndpoint + "/me"
)

// v1 APIs
const (
	GetCurrentUserProfileEndpoint     = CurrentUserEndpoint
	GetCurrentUserSavedTracksEndpoint = CurrentUserEndpoint + "/tracks"
)

// HTTP request query params
const (
	CodeReqQueryParam              = "code"
	StateReqQueryParam             = "state"
	AuthorizationCodeReqQueryParam = "authorization_code"
)

// HTTP request body params
const (
	CodeReqBodyParam        = "code"
	RedirectURIReqBodyParam = "redirect_uri"
	GrantTypeReqBodyParam   = "grant_type"
)

// HTTP request headers - values
const ApplicationFormUrlEncodedContentType = "application/x-www-form-urlencoded"

// HTTP request headers
const (
	ContentTypeHeader   = "Content-Type"
	AuthorizationHeader = "Authorization"
)

func getAllScopes() []Scope {
	return []Scope{UserReadPrivateScope, UserEmailReadScope, UserLibraryReadScope}
}

func getScopeString() string {
	allScopes := getAllScopes()
	scopeString := ""
	for i := 0; i < len(allScopes); i++ {
		if i != len(allScopes)-1 {
			scopeString += string(allScopes[i]) + "+"
		} else {
			scopeString += string(allScopes[i])
		}
	}

	return scopeString
}

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil
	}
	return result
}
