package APIHandler

import (
	"net/http"
	"strings"

	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type ExecuteHandler func(req *http.Request, res http.ResponseWriter, parsedEntity interface{})

type APIManager struct {
	apiMap    map[string](*apiInfo)
	executors []Executor
}

var sAPIManager APIManager = APIManager{
	apiMap:    map[string]*apiInfo{},
	executors: []Executor{},
}

func CreateLocalAPIManager() *APIManager {
	ret := APIManager{}

	for _, parser := range sAPIManager.executors {
		ret.RegisterParser(parser)
	}

	return &ret
}

func RegisterDefaultAPI(executor Executor, params ...interface{}) {
	tempApiInfo := &apiInfo{
		executor: executor,
	}

	if len(params) > 0 {
		for _, paramInf := range params {
			if paramAuthHandler, assertionOK := paramInf.(Authorizer); assertionOK {
				tempApiInfo.authrizer = paramAuthHandler
			}
		}
	}

	sAPIManager.apiMap["/"] = tempApiInfo
}

func RegisterAPI(path string, executor Executor, params ...interface{}) {
	absPath := path

	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}

	tempApiInfo := &apiInfo{
		executor: executor,
	}

	if len(params) > 0 {
		for _, paramInf := range params {
			if paramAuthHandler, assertionOK := paramInf.(Authorizer); assertionOK {
				tempApiInfo.authrizer = paramAuthHandler
			}
		}
	}

	sAPIManager.apiMap[absPath] = tempApiInfo
}

func RegisterParser(executor Executor) {
	sAPIManager.executors = append(sAPIManager.executors, executor)
}

func ExecuteRequest(req *http.Request, res http.ResponseWriter) {
	sAPIManager.ExecuteRequest(req, res)
}

func (manager *APIManager) RegisterDefaultAPI(executor Executor, params ...interface{}) *APIManager {
	tempApiInfo := &apiInfo{
		executor: executor,
	}

	if len(params) > 0 {
		for _, paramInf := range params {
			if paramAuthHandler, assertionOK := paramInf.(Authorizer); assertionOK {
				tempApiInfo.authrizer = paramAuthHandler
			}
		}
	}

	manager.apiMap["/"] = tempApiInfo
	return manager
}

func (manager *APIManager) RegisterAPI(path string, executor Executor, params ...interface{}) *APIManager {
	absPath := path

	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}

	tempApiInfo := &apiInfo{
		executor: executor,
	}

	if len(params) > 0 {
		for _, paramInf := range params {
			if paramAuthHandler, assertionOK := paramInf.(Authorizer); assertionOK {
				tempApiInfo.authrizer = paramAuthHandler
			}
		}
	}

	manager.apiMap[absPath] = tempApiInfo
	return manager
}

func (manager *APIManager) RegisterParser(executor Executor) *APIManager {
	manager.executors = append(manager.executors, executor)
	return manager
}

func (manager *APIManager) ExecuteRequest(req *http.Request, res http.ResponseWriter) {
	path := req.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	executorApiInfo := (*apiInfo)(nil)

	if apiInfo, exist := manager.apiMap[path]; exist {
		executorApiInfo = apiInfo
	} else {
		if apiInfo, exist := manager.apiMap["/"]; exist {
			executorApiInfo = apiInfo
		}
	}

	if executorApiInfo != nil {
		if parsedEntity, parseErr := executorApiInfo.executor.ParseRequestBody(req); parseErr == nil {
			executorApiInfo.executor.Execute(req, res, parsedEntity)
		} else {
			ThcompUtility.LogfE("fail to parse request entity: %v", parseErr)
			res.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		ThcompUtility.LogfE("not register executor: %s", path)
		res.WriteHeader(http.StatusBadRequest)
	}
}
