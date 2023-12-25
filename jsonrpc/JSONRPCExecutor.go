package jsonroc

import (
	"encoding/json"
	"fmt"
	"net/http"

	root "github.com/thcomp/GoLang_APIHandler"
)

type JSONRPCExecutor struct {
	ExecutorMap map[string](root.ExecuteHandler)
}

func (parser *JSONRPCExecutor) ParseRequest(req *http.Request) (ret interface{}, retErr error) {
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

	return
}

func (parser *JSONRPCExecutor) ParseResponse(res *http.Response) (ret interface{}, retErr error) {
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

	return
}

func (parser *JSONRPCExecutor) RegisterExecuteHandler(condMap map[string]interface{}, handler root.ExecuteHandler) (err error) {
	if methodInf, exist := condMap["method"]; exist {
		if method, assertionOK := methodInf.(string); assertionOK {
			if parser.ExecutorMap == nil {
				parser.ExecutorMap = map[string](root.ExecuteHandler){}
			}

			parser.ExecutorMap[method] = handler
		} else {
			err = fmt.Errorf("method format not string")
		}
	} else {
		err = fmt.Errorf("method not exist in condMap")
	}

	return err
}

func (parser *JSONRPCExecutor) Execute(req *http.Request, res http.ResponseWriter, parsedEntity interface{}) {
	if jsonReq, assertionOK := parsedEntity.(*JSONRPCRequest); assertionOK {
		if handler, exist := parser.ExecutorMap[jsonReq.Method]; exist {
			handler(req, res, jsonReq)
		}
	}
}
