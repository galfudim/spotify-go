package http

// HTTP request query params
const (
	ReqQueryParamCode              = "code"
	ReqQueryParamState             = "state"
	ReqQueryParamAuthorizationCode = "authorization_code"
)

// HTTP request body params
const (
	ReqBodyParamCode        = "code"
	ReqBodyParamRedirectURI = "redirect_uri"
	ReqBodyParamGrantType   = "grant_type"
)

// HTTP request headers - values
const ContentTypeApplicationFormUrlEncoded = "application/x-www-form-urlencoded"

// HTTP request headers
const (
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
)
