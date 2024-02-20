package APIHandler

import "net/http"

type Authorizer interface {
	Authorize(*http.Request) (*AuthorizedUser, error)
	Authenticate(http.ResponseWriter)
}
