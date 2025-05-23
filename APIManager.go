package APIHandler

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type ExecuteHandler func(req *http.Request, res http.ResponseWriter, parsedEntity interface{}, authUser AuthorizedUser)

type APIManager struct {
	apiMap map[string](*apiInfo)
}

var sAPIManager APIManager = APIManager{
	apiMap: map[string]*apiInfo{},
}

func CreateLocalAPIManager() *APIManager {
	return &APIManager{}
}

func RegisterDefaultAPI(executor Executor, params ...interface{}) {
	tempApiInfo := &apiInfo{
		executor: executor,
	}

	if len(params) > 0 {
		for _, paramInf := range params {
			if paramAuthHandler, assertionOK := paramInf.(Authorizer); assertionOK {
				tempApiInfo.authorizer = paramAuthHandler
			}
		}
	}

	sAPIManager.apiMap["/"] = tempApiInfo
}

func RegisterAPI(method *string, path string, executor Executor, params ...interface{}) {
	absPath := path

	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}

	tempApiInfo := &apiInfo{
		method:   method,
		executor: executor,
	}

	if len(params) > 0 {
		for _, paramInf := range params {
			if paramAuthHandler, assertionOK := paramInf.(Authorizer); assertionOK {
				tempApiInfo.authorizer = paramAuthHandler
			}
		}
	}

	if method != nil {
		absPath = absPath + `_` + strings.ToLower(*method)
	}
	sAPIManager.apiMap[absPath] = tempApiInfo
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
				tempApiInfo.authorizer = paramAuthHandler
			}
		}
	}

	manager.apiMap["/"] = tempApiInfo
	return manager
}

func (manager *APIManager) RegisterAPI(method *string, path string, executor Executor, params ...interface{}) *APIManager {
	absPath := path

	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}

	tempApiInfo := &apiInfo{
		method:   method,
		executor: executor,
	}

	if len(params) > 0 {
		for _, paramInf := range params {
			if paramAuthHandler, assertionOK := paramInf.(Authorizer); assertionOK {
				tempApiInfo.authorizer = paramAuthHandler
			}
		}
	}

	if method != nil {
		absPath = absPath + `_` + strings.ToLower(*method)
	}
	manager.apiMap[absPath] = tempApiInfo
	return manager
}

func (manager *APIManager) ExecuteRequest(req *http.Request, res http.ResponseWriter) {
	path := req.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	executorApiInfo := (*apiInfo)(nil)

	pathWithMethod := path + `_` + strings.ToLower(req.Method)

	if apiInfo, exist := manager.apiMap[pathWithMethod]; exist {
		executorApiInfo = apiInfo
	} else if apiInfo, exist := manager.apiMap[path]; exist {
		executorApiInfo = apiInfo
	} else {
		if apiInfo, exist := manager.apiMap["/"]; exist {
			executorApiInfo = apiInfo
		}
	}

	if executorApiInfo != nil {
		authorized := (*bool)(nil)
		authorizedUser := (AuthorizedUser)(nil)

		if executorApiInfo.IsAuthorizeByHttpHeader() {
			if tempAuthorizedUser, authErr := executorApiInfo.authorizer.Authorize(req); authErr == nil {
				authorizedUser = tempAuthorizedUser
				tempAuthorized := true
				authorized = &tempAuthorized
			} else {
				tempAuthorized := false
				authorized = &tempAuthorized
			}
		}

		if authorized == nil || (*authorized) {
			entityReader := (*ThcompUtility.NopCloser)(nil)
			if entity, readErr := ioutil.ReadAll(req.Body); readErr == nil {
				entityReader = ThcompUtility.NewNopCloser(bytes.NewReader(entity))
				req.Body = entityReader
				entityReader.Seek(0, io.SeekStart)

				if authorized == nil && executorApiInfo.authorizer.AuthorizeBy() == ByHttpEntity {
					if tempAuthorizedUser, authErr := executorApiInfo.authorizer.Authorize(req); authErr == nil {
						authorizedUser = tempAuthorizedUser
						tempAuthorized := true
						authorized = &tempAuthorized
					} else {
						tempAuthorized := false
						authorized = &tempAuthorized
					}
				}

				if authorized == nil || (*authorized) {
					if parsedEntity, parseErr := executorApiInfo.executor.ParseRequestBody(req); parseErr == nil {
						entityReader.Seek(0, io.SeekStart)

						executorApiInfo.executor.Execute(req, res, authorizedUser, parsedEntity)
					} else {
						ThcompUtility.LogfE("fail to parse request entity: %v", parseErr)
						res.WriteHeader(http.StatusInternalServerError)
					}
				} else {
					executorApiInfo.authorizer.Authenticate(res)
				}
			} else {
				res.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			executorApiInfo.authorizer.Authenticate(res)
		}
	} else {
		ThcompUtility.LogfE("not register executor: %s", path)
		res.WriteHeader(http.StatusNotFound)
	}
}
