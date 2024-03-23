package http

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
