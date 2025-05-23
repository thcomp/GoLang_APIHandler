package APIHandler

import (
	"errors"
)

type apiInfo struct {
	method     *string
	executor   Executor
	authorizer Authorizer
}

func (info *apiInfo) Method() *string {
	return info.method
}

func (info *apiInfo) IsAuthorizeByHttpHeader() bool {
	return info.authorizer != nil && info.authorizer.AuthorizeBy() == ByHttpHeader
}

func (info *apiInfo) IsAuthorizeByHttpEntity() bool {
	return info.authorizer != nil && info.authorizer.AuthorizeBy() == ByHttpEntity
}

var ErrUnsupportEntity = errors.New("unsupport entity")
