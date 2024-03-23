package spotify

const (
	AccountsEndpoint      = "https://accounts.spotify.com"
	ApiEndpoint           = "https://api.spotify.com"
	v1ApiEndpoint         = ApiEndpoint + "/v1"
	AuthorizeUserEndpoint = AccountsEndpoint + "/authorize"
	GenerateTokenEndpoint = AccountsEndpoint + "/api/token"
)

// v1 APIs
const (
	GetCurrentUserProfileEndpoint = v1ApiEndpoint + "/me"
)
