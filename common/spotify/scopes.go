package spotify

type Scope string

const (
	UserReadPrivateScope Scope = "user-read-private"
	UserEmailReadScope   Scope = "user-read-email"

	UserLibraryReadScope Scope = "user-library-read"
)

func GetAllScopes() []Scope {
	return []Scope{UserReadPrivateScope, UserEmailReadScope, UserLibraryReadScope}
}
