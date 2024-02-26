package multipart

import (
	"fmt"
	"net/http"
	"strings"

	root "github.com/thcomp/GoLang_APIHandler"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type MultipartFormExecutor struct {
	cacheEditorFactory ThcompUtility.CacheEditorFactory
	handler            root.ExecuteHandler
}

func (parser *MultipartFormExecutor) RegisterExecuteHandler(condMap map[string]string, handler root.ExecuteHandler) *MultipartFormExecutor {
	return parser
}

func (parser *MultipartFormExecutor) CacheEditorFactory(cacheEditorFactory ThcompUtility.CacheEditorFactory) {
	parser.cacheEditorFactory = cacheEditorFactory
}

func (parser *MultipartFormExecutor) ParseRequest(req *http.Request) (ret interface{}, retErr error) {
	formData := (*MultipartFormData)(nil)
	if multipartHelper, err := ThcompUtility.NewMultipartHelperFromHttpRequest(req); err == nil {
		formData = &MultipartFormData{helper: multipartHelper}
	} else {
		retErr = err
	}

	return formData, retErr
}

func (parser *MultipartFormExecutor) ParseResponse(res *http.Response) (ret interface{}, retErr error) {
	if contentTypeValue := res.Header.Get(`Content-type`); contentTypeValue != `` {
		lowerContentTypeValue := strings.ToLower(contentTypeValue)
		if strings.HasPrefix(lowerContentTypeValue, `multipart/form-data`) {
			boundaryText := (*string)(nil)
			partialTexts := strings.Split(contentTypeValue, `boundary=`)
			if len(partialTexts) >= 2 {
				for i := 1; i < len(partialTexts); i++ {
					prevPart := strings.TrimRight(partialTexts[i-1], " \t")
					if prevPart[len(prevPart)-1] == ';' {
						subPartialTexts := strings.Split(partialTexts[i], ";")
						boundaryText = &subPartialTexts[0]
						break
					}

				}
			}

			if boundaryText != nil {
				cacheEditorFactory := parser.cacheEditorFactory
				if parser.cacheEditorFactory == nil {
					cacheEditorFactory = ThcompUtility.NewMemoryCacheEditorFactory()
				}

				if multipartHelper, err := ThcompUtility.NewMultipartHelper(res.Body, *boundaryText, cacheEditorFactory); err == nil {
					ret = &MultipartFormData{
						helper: multipartHelper,
					}
				} else {
					retErr = err
				}
			} else {
				retErr = fmt.Errorf("not exist boundary text in content-type header: %s", contentTypeValue)
			}
		} else {
			retErr = fmt.Errorf("not multipart/form-data: %s", contentTypeValue)
		}
	} else {
		retErr = fmt.Errorf("not exist content-type header")
	}

	return
}

func (parser *MultipartFormExecutor) IsMultipartFormData(headers http.Header) (ret bool) {
	mimetype := headers.Get(`Content-type`)
	lowerMimetype := strings.ToLower(mimetype)
	if strings.HasPrefix(lowerMimetype, `multipart/form-data`) {
		ret = true
	}

	return
}

func (parser *MultipartFormExecutor) Execute(req *http.Request, res http.ResponseWriter, parsedEntity interface{}) {
	if multipartFormData, assertionOK := parsedEntity.(*MultipartFormData); assertionOK {
		defer multipartFormData.Close()

		parser.handler(req, res, multipartFormData)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
