package authorizer

import (
	"fmt"
	"net/http"
	"strings"

	root "github.com/thcomp/GoLang_APIHandler"
	//ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type BasicAuthorizer struct {
	userAuthorizer func(authorizer *BasicAuthorizer, credentials string) *root.AuthorizedUser
	userAuthCache  interface{}
}

func (authorizer *BasicAuthorizer) RegisterUserAuthorizer(userAuthorizer func(authorizer *BasicAuthorizer, credentials string) *root.AuthorizedUser) {
	authorizer.userAuthorizer = userAuthorizer
}

func (authorizer *BasicAuthorizer) SetUserAuthCache(cache interface{}) {
	authorizer.userAuthCache = cache
}

func (authorizer *BasicAuthorizer) GetUserAuthCache() interface{} {
	return authorizer.userAuthCache
}

func (authorizer *BasicAuthorizer) Authorize(req *http.Request) (user *root.AuthorizedUser, err error) {
	if authorizer.userAuthorizer != nil {
		authHeader := req.Header.Get("Authorization")
		authHeader = strings.TrimLeft(authHeader, " \t")
		authHeaderParts := strings.Split(authHeader, " ")
		useNextPart := false

		for pos, authHeaderPart := range authHeaderParts {
			if pos == 0 {
				lowerAuthHeaderPart := strings.ToLower(authHeaderPart)
				if lowerAuthHeaderPart == "basic" {
					useNextPart = true
				} else {
					err = fmt.Errorf("not support authorize type: %s", authHeaderPart)
					break
				}
			} else if useNextPart {
				if authHeaderPart != "" {
					user = authorizer.userAuthorizer(authorizer, authHeaderPart)
					if user == nil {
						err = fmt.Errorf("not matched user")
					}
					break
				}
			}
		}
	} else {
		err = fmt.Errorf("not register user authorizer")
	}

	return
}
