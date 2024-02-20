package APIHandler

import (
	"errors"
	"net/http"
)

type Executor interface {
	ParseRequestBody(req *http.Request) (interface{}, error)
	ParseResponseBody(res *http.Response) (interface{}, error)
	//RegisterExecuteHandler(condMap map[string]interface{}, handler ExecuteHandler) (err error)
	Execute(req *http.Request, res http.ResponseWriter, parsedEntity interface{})
}

type apiInfo struct {
	executor  Executor
	authrizer Authorizer
}

var ErrUnsupportEntity = errors.New("unsupport entity")
