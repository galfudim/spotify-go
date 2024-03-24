package spotify

const (
	AccountsEndpoint      = "https://accounts.spotify.com"
	ApiEndpoint           = "https://api.spotify.com"
	v1ApiEndpoint         = ApiEndpoint + "/v1"
	AuthorizeUserEndpoint = AccountsEndpoint + "/authorize"
	GenerateTokenEndpoint = AccountsEndpoint + "/api/token"
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
