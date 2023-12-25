package APIHandler

type AuthorizedUser interface {
	ID() string
	DisplayName() string
	IsExpired() bool
}
