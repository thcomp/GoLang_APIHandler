package jsonrpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	root "github.com/thcomp/GoLang_APIHandler"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type JSONRPCExecutor struct {
	ExecutorMap map[string](root.ExecuteHandler)
}

const CondMapKeyMethod = "method"

func NewJSONRPCExecutor() *JSONRPCExecutor {
	return &JSONRPCExecutor{}
}

func (parser *JSONRPCExecutor) RegisterExecuteHandler(condMap map[string]interface{}, handler root.ExecuteHandler) *JSONRPCExecutor {
	if methodInf, exist := condMap[CondMapKeyMethod]; exist {
		if method, assertionOK := methodInf.(string); assertionOK {
			if parser.ExecutorMap == nil {
				parser.ExecutorMap = map[string](root.ExecuteHandler){}
			}

			parser.ExecutorMap[method] = handler
		} else {
			ThcompUtility.LogfE("%s format not string", CondMapKeyMethod)
		}
	} else {
		ThcompUtility.LogfE("%s not exist in condMap", CondMapKeyMethod)
	}

	return parser
}

func (parser *JSONRPCExecutor) ParseRequest(req *http.Request) (ret interface{}, retErr error) {
	if parser.IsJSON(req.Header) {
		jsonReq := JSONRPCRequest{}
		if parseErr := json.NewDecoder(req.Body).Decode(&jsonReq); parseErr == nil {
			if jsonReq.Version != "2.0" || jsonReq.Method == "" {
				retErr = fmt.Errorf("can not parse on JSONRPC Request")
			} else {
				ret = &jsonReq
			}
		} else {
			retErr = parseErr
		}
	} else {
		retErr = root.ErrUnsupportEntity
	}

	return
}

func (parser *JSONRPCExecutor) ParseResponse(res *http.Response) (ret interface{}, retErr error) {
	if parser.IsJSON(res.Header) {
		jsonRes := JSONRPCResponse{}
		if parseErr := json.NewDecoder(res.Body).Decode(&jsonRes); parseErr == nil {
			if jsonRes.Version != "2.0" {
				retErr = fmt.Errorf("can not parse on JSONRPC Request")
			} else {
				ret = &jsonRes
			}
		} else {
			retErr = parseErr
		}
	} else {
		retErr = root.ErrUnsupportEntity
	}

	return
}

func (parser *JSONRPCExecutor) IsJSON(headers http.Header) (ret bool) {
	mimetype := headers.Get("Content-type")
	lowerMimetype := strings.ToLower(mimetype)
	if strings.HasPrefix(lowerMimetype, "application/json") {
		ret = true
	}

	return
}

func (parser *JSONRPCExecutor) Execute(req *http.Request, res http.ResponseWriter, authUser root.AuthorizedUser, parsedEntity interface{}) {
	if jsonReq, assertionOK := parsedEntity.(*JSONRPCRequest); assertionOK {
		if handler, exist := parser.ExecutorMap[jsonReq.Method]; exist {
			handler(req, res, jsonReq, authUser)
		}
	}
}
