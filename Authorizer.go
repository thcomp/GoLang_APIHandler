package APIHandler

import "net/http"

type AuthorizeBy int

const (
	ByHttpHeader AuthorizeBy = iota
	ByHttpEntity
)

type Authorizer interface {
	AuthorizeBy() AuthorizeBy
	Authorize(*http.Request) (*AuthorizedUser, error)
	Authenticate(http.ResponseWriter)
}
