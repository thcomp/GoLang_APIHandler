package urlenc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	root "github.com/thcomp/GoLang_APIHandler"
)

type sExecutorInfo struct {
	value   string
	handler root.ExecuteHandler
}
type URLEncExecutor struct {
	ExecutorMap map[string](*sExecutorInfo)
}

func (parser *URLEncExecutor) RegisterExecuteHandler(condMap map[string]string, handler root.ExecuteHandler) *URLEncExecutor {
	if len(condMap) > 0 {
		if parser.ExecutorMap == nil {
			parser.ExecutorMap = map[string](*sExecutorInfo){}
		}

		for key, value := range condMap {
			parser.ExecutorMap[key] = &sExecutorInfo{
				value:   value,
				handler: handler,
			}
		}
	}

	return parser
}

func (parser *URLEncExecutor) parseEntity(reader io.Reader) (ret *URLEncData, retErr error) {
	urlEncData := URLEncData{}
	if valueBytes, readErr := ioutil.ReadAll(reader); readErr == nil {
		if queryValues, parseErr := url.ParseQuery(string(valueBytes)); parseErr == nil {
			urlEncData.queryValues = &queryValues
			ret = &urlEncData
		} else {
			retErr = parseErr
		}
	} else {
		retErr = readErr
	}

	return
}

func (parser *URLEncExecutor) ParseRequest(req *http.Request) (ret interface{}, retErr error) {
	reader := io.Reader(nil)

	if contentTypeValue := req.Header.Get(`Content-type`); contentTypeValue != `` {
		contentTypeValue = strings.ToLower(contentTypeValue)
		if strings.HasPrefix(contentTypeValue, `application/x-www-form-urlencoded`) {
			if originalData, readErr := ioutil.ReadAll(req.Body); readErr == nil {
				reader = bytes.NewReader(originalData)
			} else {
				retErr = readErr
			}
		} else {
			if len(req.URL.RawQuery) > 0 {
				reader = bytes.NewReader([]byte(req.URL.RawQuery))
			} else {
				retErr = fmt.Errorf("no url encoded value in request")
			}
		}
	} else {
		if len(req.URL.RawQuery) > 0 {
			reader = bytes.NewReader([]byte(req.URL.RawQuery))
		} else {
			retErr = fmt.Errorf("no url encoded value in request")
		}
	}

	if reader != nil {
		return parser.parseEntity(reader)
	} else {
		return nil, retErr
	}
}

func (parser *URLEncExecutor) ParseResponse(res *http.Response) (ret interface{}, retErr error) {
	reader := io.Reader(nil)

	if contentTypeValue := res.Header.Get(`Content-type`); contentTypeValue != `` {
		contentTypeValue = strings.ToLower(contentTypeValue)
		if strings.HasPrefix(contentTypeValue, `application/x-www-form-urlencoded`) {
			if originalData, readErr := ioutil.ReadAll(res.Body); readErr == nil {
				reader = bytes.NewReader(originalData)
			} else {
				retErr = readErr
			}
		}
	}

	if reader != nil {
		return parser.parseEntity(reader)
	} else {
		return nil, retErr
	}
}

func (parser *URLEncExecutor) Execute(req *http.Request, res http.ResponseWriter, authUser root.AuthorizedUser, parsedEntity interface{}) {
	if urlEncData, assertionOK := parsedEntity.(*URLEncData); assertionOK {
		for queryKey, queryValues := range *urlEncData.queryValues {
			if executorInfo, exist := parser.ExecutorMap[queryKey]; exist {
				matched := false

				for _, queryValue := range queryValues {
					if executorInfo.value == queryValue {
						executorInfo.handler(req, res, parsedEntity, authUser)
						matched = true
						break
					}
				}

				if matched {
					break
				}
			}
		}
	}
}
