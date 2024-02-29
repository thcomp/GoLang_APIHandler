package authorizer

import (
	"fmt"
	"net/http"
	"strings"

	root "github.com/thcomp/GoLang_APIHandler"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type BasicAuthorizer struct {
	userAuthorizer func(authorizer *BasicAuthorizer, credentials string) root.AuthorizedUser
	userAuthCache  interface{}
	authParams     map[string]string
}

func (authorizer *BasicAuthorizer) RegisterUserAuthorizer(userAuthorizer func(authorizer *BasicAuthorizer, credentials string) root.AuthorizedUser) {
	authorizer.userAuthorizer = userAuthorizer
}

func (authorizer *BasicAuthorizer) SetUserAuthCache(cache interface{}) {
	authorizer.userAuthCache = cache
}

func (authorizer *BasicAuthorizer) GetUserAuthCache() interface{} {
	return authorizer.userAuthCache
}

func (authorizer *BasicAuthorizer) SetRealm(realm string) {
	if authorizer.authParams == nil {
		authorizer.authParams = map[string]string{}
	}

	authorizer.authParams["realm"] = realm
}

func (authorizer *BasicAuthorizer) SetCharset(charset string) {
	if authorizer.authParams == nil {
		authorizer.authParams = map[string]string{}
	}

	authorizer.authParams["charset"] = charset
}

func (authorizer *BasicAuthorizer) AuthorizeBy() root.AuthorizeBy {
	return root.ByHttpHeader
}

func (authorizer *BasicAuthorizer) Authorize(req *http.Request) (user root.AuthorizedUser, err error) {
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

func (authorizer *BasicAuthorizer) Authenticate(res http.ResponseWriter) {
	builder := ThcompUtility.StringBuilder{}

	if authorizer.authParams != nil && len(authorizer.authParams) > 0 {
		paramCount := 0
		if realm, exist := authorizer.authParams["realm"]; exist {
			builder.Appendf("realm=\"%s\"", realm)
			paramCount++
		}

		if charset, exist := authorizer.authParams["charset"]; exist {
			format := "charset=\"%s\""
			if paramCount > 0 {
				format = ", " + format
			}
			builder.Appendf(format, charset)
			paramCount++
		}
	}

	res.Header().Set("Www-Autheticate", fmt.Sprintf("basic %s", builder.String()))
}
