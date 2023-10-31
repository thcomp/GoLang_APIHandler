package APIHandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"

	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

var sAutoID int64
var sAutoIDMutex sync.Mutex

type JSONRPC struct {
	Version string      `json:"jsonrpc"`
	id      interface{} `json:"id"`
}

func newJSONRPC() *JSONRPC {
	return &JSONRPC{Version: "2.0"}
}

func newJSONRPCWithID(id interface{}) *JSONRPC {
	return &JSONRPC{Version: "2.0", id: id}
}

func (rpc *JSONRPC) IsIDNum() bool {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	return infHelper.IsNumber()
}

func (rpc *JSONRPC) IsIDString() bool {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	return infHelper.IsString()
}

func (rpc *JSONRPC) IDNum() (ret float64, isNum bool) {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	ret, isNum = infHelper.GetNumber()

	return
}

func (rpc *JSONRPC) IDString() (ret string, isNum bool) {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	ret, isNum = infHelper.GetString()

	return
}

func (rpc *JSONRPC) IDInterface() (ret interface{}) {
	return rpc.id
}

type JSONRPCRequest struct {
	JSONRPC
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func NewJSONRPCNotificationRequest(method string, params interface{}) (*JSONRPCRequest, error) {
	ret := &JSONRPCRequest{
		JSONRPC: *newJSONRPC(),
		Method:  method,
	}
	retErr := error(nil)

	if reader, assertionOK := params.(io.Reader); assertionOK {
		if paramsBytes, readErr := ioutil.ReadAll(reader); readErr == nil {
			ret.Params = paramsBytes
		} else {
			retErr = readErr
		}
	} else {
		ret.Params = params
	}

	return ret, retErr
}

func NewJSONRPCRequest(id interface{}, method string, params interface{}) (*JSONRPCRequest, error) {
	ret := &JSONRPCRequest{
		JSONRPC: *newJSONRPCWithID(id),
		Method:  method,
	}
	retErr := error(nil)

	if reader, assertionOK := params.(io.Reader); assertionOK {
		if paramsBytes, readErr := ioutil.ReadAll(reader); readErr == nil {
			ret.Params = paramsBytes
		} else {
			retErr = readErr
		}
	} else {
		ret.Params = params
	}

	return ret, retErr
}

func ParseJSONRequest(reader io.Reader) (*JSONRPCRequest, error) {
	ret := (*JSONRPCRequest)(nil)
	tempRet := map[string]interface{}{}
	retErr := json.NewDecoder(reader).Decode(&tempRet)

	if retErr == nil {
		ret = &JSONRPCRequest{}
		if valueInf, exist := tempRet["jsonrpc"]; exist {
			valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
			if valueInfHelper.IsString() {
				ret.JSONRPC.Version, _ = valueInfHelper.GetString()
			}
		}
		if valueInf, exist := tempRet["id"]; exist {
			ret.JSONRPC.id = valueInf
		}
		if valueInf, exist := tempRet["method"]; exist {
			valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
			if valueInfHelper.IsString() {
				ret.Method, _ = valueInfHelper.GetString()
			}
		}
		if valueInf, exist := tempRet["params"]; exist {
			ret.Params = valueInf
		}
	}

	return ret, retErr
}

type JSONRPCResponse struct {
	JSONRPC
	Result interface{}   `json:"result"`
	Error  *JSONRPCError `json:"error"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const JSONRPCParseError = -32700
const JSONRPCInvalidRequest = -32600
const JSONRPCMethodNotFound = -32601
const JSONRPCInvalidParams = -32602
const JSONRPCInternalError = -32603
const JSONRPCServerErrorMax = -32000
const JSONRPCServerErrorMin = -32099

func NewJSONRPCResponse(id interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: *newJSONRPCWithID(id),
	}
}

func NewJSONRPCResponseFromRequest(request *JSONRPCRequest) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: request.JSONRPC,
	}
}

func ParseJSONResponse(reader io.Reader) (*JSONRPCResponse, error) {
	ret := (*JSONRPCResponse)(nil)
	tempRet := map[string]interface{}{}
	retErr := json.NewDecoder(reader).Decode(&tempRet)

	if retErr == nil {
		ret = &JSONRPCResponse{}
		if valueInf, exist := tempRet["jsonrpc"]; exist {
			valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
			if valueInfHelper.IsString() {
				ret.JSONRPC.Version, _ = valueInfHelper.GetString()
			}
		}
		if valueInf, exist := tempRet["id"]; exist {
			ret.JSONRPC.id = valueInf
		}
		if valueInf, exist := tempRet["result"]; exist {
			ret.Result = valueInf
		}
		if valueInf, exist := tempRet["error"]; exist {
			if errorMap, assertionOK := valueInf.(map[string]interface{}); assertionOK {
				ret.Error = &JSONRPCError{}
				if valueInf, exist := errorMap["code"]; exist {
					valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
					if valueInfHelper.IsNumber() {
						tempValue, _ := valueInfHelper.GetNumber()
						ret.Error.Code = int(tempValue)
					}
				}
				if valueInf, exist := errorMap["message"]; exist {
					valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
					if valueInfHelper.IsString() {
						ret.Error.Message, _ = valueInfHelper.GetString()
					}
				}
				if valueInf, exist := errorMap["data"]; exist {
					ret.Error.Data = valueInf
				}
			}
		}
	}

	return ret, retErr
}

func NewJSONRPCError(code int, message string, data interface{}) *JSONRPCError {
	return &JSONRPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
