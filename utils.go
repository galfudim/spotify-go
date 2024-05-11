package main

const (
	ClientId     = "SPOTIFY_CLIENT_ID"
	ClientSecret = "SPOTIFY_CLIENT_SECRET"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type CommonResponse struct {
	HREF     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

type SavedTracksResponse struct {
	CommonResponse
	Items []TrackDTO `json:"items"`
}

type SimplifiedPlaylistsResponse struct {
	CommonResponse
	Items []PlaylistDTO `json:"items"`
}

type PlaylistTracksResponse struct {
	CommonResponse
	Items []PlaylistTrackItem `json:"items"`
}

type TrackDTO struct {
	AddedAt string `json:"added_at"`
	Track   Track  `json:"track"`
}

type PlaylistDTO struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	ExternalURLs  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	HREF       string        `json:"href"`
	ID         string        `json:"id"`
	Images     []Image       `json:"images"`
	Name       string        `json:"name"`
	User       User          `json:"owner"`
	Public     bool          `json:"public"`
	SnapshotID string        `json:"snapshot_id"`
	Tracks     TrackMetadata `json:"tracks"`
	Type       string        `json:"type"`
	URI        string        `json:"uri"`
}

type PlaylistTrackItem struct {
	AddedAt string `json:"added_at"`
	User    User   `json:"added_by"`
	Local   bool   `json:"is_local"`
	Track   Track  `json:"track"`
}

type TrackMetadata struct {
	HREF  string `json:"href"`
	Total int    `json:"total"`
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
	ExternalURLs ExternalURLs `json:"external_urls"`
	HREF         string       `json:"href"`
	ID           string       `json:"id"`
	IsPlayable   bool         `json:"is_playable"`
	LinkedFrom   struct{}     `json:"linked_from"`
	Name         string       `json:"name"`
	Popularity   int          `json:"popularity"`
	PreviewURL   string       `json:"preview_url"`
	TrackNumber  int          `json:"track_number"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
	IsLocal      bool         `json:"is_local"`
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
	ExternalURLs ExternalURLs `json:"external_urls"`
	Followers    struct {
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

type User struct {
	Name         string       `json:"display_name"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	Followers    struct {
		HREF  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	HREF string `json:"href"`
	ID   string `json:"id"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type ExternalURLs struct {
	Spotify string `json:"spotify"`
}
