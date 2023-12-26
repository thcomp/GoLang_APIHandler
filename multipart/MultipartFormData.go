package multipart

import (
	"fmt"

	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

var sErrNotExistFormData error = fmt.Errorf("not exist form data")

type MultipartFormData struct {
	helper *ThcompUtility.MultipartHelper
}

func (formData *MultipartFormData) Close() error {
	if formData.helper != nil {
		return formData.helper.Close()
	} else {
		return sErrNotExistFormData
	}
}

func (formData *MultipartFormData) Count() int {
	if formData.helper != nil {
		return formData.helper.Count()
	} else {
		return -1
	}
}

func (formData *MultipartFormData) GetByIndex(index int) (*ThcompUtility.FormData, error) {
	if formData.helper != nil {
		return formData.helper.GetByIndex(index)
	} else {
		return nil, sErrNotExistFormData
	}
}

func (formData *MultipartFormData) GetByName(partName string) (*ThcompUtility.FormData, error) {
	if formData.helper != nil {
		return formData.helper.GetByName(partName)
	} else {
		return nil, sErrNotExistFormData
	}
}

func (formData *MultipartFormData) PartNames() []string {
	if formData.helper != nil {
		return formData.helper.PartNames()
	} else {
		return nil
	}
}
