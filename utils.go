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

func jsonToMap(jsonBytes []byte) map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil
	}
	return result
}

type TrackResponse struct {
	HREF     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []Item `json:"items"`
}

type Item struct {
	AddedAt string `json:"added_at"`
	Track   Track  `json:"track"`
}

type Track struct {
	Album            Album    `json:"album"`
	Artists          []Artist `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	DurationMS       int      `json:"duration_ms"`
	Explicit         bool     `json:"explicit"`
	ExternalIDs      struct {
		ISRC string `json:"isrc"`
		EAN  string `json:"ean"`
		UPC  string `json:"upc"`
	} `json:"external_ids"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	HREF        string   `json:"href"`
	ID          string   `json:"id"`
	IsPlayable  bool     `json:"is_playable"`
	LinkedFrom  struct{} `json:"linked_from"`
	Name        string   `json:"name"`
	Popularity  int      `json:"popularity"`
	PreviewURL  string   `json:"preview_url"`
	TrackNumber int      `json:"track_number"`
	Type        string   `json:"type"`
	URI         string   `json:"uri"`
	IsLocal     bool     `json:"is_local"`
}

type Album struct {
	AlbumType        string   `json:"album_type"`
	TotalTracks      int      `json:"total_tracks"`
	AvailableMarkets []string `json:"available_markets"`
	ExternalURLs     struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	HREF                 string  `json:"href"`
	ID                   string  `json:"id"`
	Images               []Image `json:"images"`
	Name                 string  `json:"name"`
	ReleaseDate          string  `json:"release_date"`
	ReleaseDatePrecision string  `json:"release_date_precision"`
	Restrictions         struct {
		Reason string `json:"reason"`
	} `json:"restrictions"`
	Type    string   `json:"type"`
	URI     string   `json:"uri"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		HREF  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Genres     []string `json:"genres"`
	HREF       string   `json:"href"`
	ID         string   `json:"id"`
	Images     []Image  `json:"images"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
	Type       string   `json:"type"`
	URI        string   `json:"uri"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
