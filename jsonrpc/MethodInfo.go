package jsonrpc

import (
	root "github.com/thcomp/GoLang_APIHandler"
)

type methodInfo struct {
	executeHandler root.ExecuteHandler
	authorizer     root.Authorizer
}
