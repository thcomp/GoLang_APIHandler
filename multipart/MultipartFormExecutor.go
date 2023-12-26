package multipart

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type MultipartFormExecutor struct {
}

func (parser *MultipartFormExecutor) ParseRequest(req *http.Request) (ret interface{}, retErr error) {
	formData := (*&MultipartFormData)(nil)
	if multipartHelper, err := ThcompUtility.NewMultipartHelperFromHttpRequest(req); err == nil {
		formData = &MultipartFormData{helper: multipartHelper}
	} else {
		retErr = err
	}

	return formData, retErr
}

func (parser *MultipartFormExecutor) ParseResponse(res *http.Response) (ret interface{}, retErr error) {
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

func (parser *MultipartFormExecutor) RegisterExecuteHandler(condMap map[string]string, handler root.ExecuteHandler) (err error) {
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

	return err
}

func (parser *MultipartFormExecutor) Execute(req *http.Request, res http.ResponseWriter, parsedEntity interface{}) {
	if urlEncData, assertionOK := parsedEntity.(*URLEncData); assertionOK {
		for queryKey, queryValues := range *urlEncData.queryValues {
			if executorInfo, exist := parser.ExecutorMap[queryKey]; exist {
				matched := false

				for _, queryValue := range queryValues {
					if executorInfo.value == queryValue {
						executorInfo.handler(req, res, parsedEntity)
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
