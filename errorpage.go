package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// getErrorPageMethods returns Error Page related API calls
//
func getErrorPageMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listErrorPages": {
			Name:   "listErrorPages",
			Action: "domain/%d/errorpages",
			Method: http.MethodGet,
			Result: []ErrorPage{},
		},
		"getErrorPage": {
			Name:               "getErrorPage",
			Action:             "domain/%d/errorpages/%d",
			Method:             http.MethodGet,
			Result:             ErrorPage{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"createErrorPage": {
			Name:               "createErrorPage",
			Action:             "domain/%d/errorpages",
			Method:             http.MethodPost,
			Result:             ErrorPage{},
			ResponseDecodeFunc: decodeErrorPageResponse,
		},
		"updateErrorPage": {
			Name:               "updateErrorPage",
			Action:             "domain/%d/errorpages",
			Method:             http.MethodPost,
			Result:             ErrorPage{},
			ResponseDecodeFunc: decodeErrorPageResponse,
		},
		"deleteErrorPage": {
			Name:   "deleteErrorPage",
			Action: "domain/%d/errorpages",
			Method: http.MethodDelete,
			Result: ErrorPage{},
		},
	}
}

//
// errorPageUpdate
//
type errorPageUpdate struct {
	ID          int                     `json:"id,omitempty"`
	PageContent string                  `json:"pageContent,omitempty"`
	Selection   map[string]map[int]bool `json:"selection,omitempty"`
	Created     *types.DateTime         `json:"created,omitempty"`
	Modified    *types.DateTime         `json:"modified,omitempty"`
}

//
// ErrorPage
//
type ErrorPage struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	ErrorCode     int             `json:"errorCode,omitempty"`
	Content       string          `json:"content,omitempty"`
	SubDomainName string          `json:"subDomainName,omitempty"`
}

//
// GetErrorPage returns a single error page with/for the given identifier
//
func (api *API) GetErrorPage(domainId int, pageId int) (*ErrorPage, error) {
	if _, ok := methods["getErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getErrorPage")
	}

	definition := methods["getErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, pageId)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*ErrorPage), nil
}

//
// ListErrorPages returns a slice containing all error pages
//
func (api *API) ListErrorPages(domainId int, params map[string]string) ([]ErrorPage, error) {
	if _, ok := methods["listErrorPages"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listErrorPages")
	}

	definition := methods["listErrorPages"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []ErrorPage
	records = append(records, *result.(*[]ErrorPage)...)

	return records, nil
}

//
// CreateErrorPage creates a new error page using the MYRA API
//
func (api *API) CreateErrorPage(errorPage *ErrorPage, domainId int) (*ErrorPage, error) {
	if _, ok := methods["createErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createErrorPage")
	}

	definition := methods["createErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	errorPageUpdate := convertErrorPageToErrorPageUpdate(errorPage)
	_, err := api.call(definition, errorPageUpdate)
	if err != nil {
		return nil, err
	}

	return errorPage, nil
}

//
// UpdateErrorPage updates the passed error page using the MYRA API
//
func (api *API) UpdateErrorPage(errorPage *ErrorPage, domainId int) (*ErrorPage, error) {
	if _, ok := methods["updateErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateErrorPage")
	}

	definition := methods["updateErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	errorPageUpdate := convertErrorPageToErrorPageUpdate(errorPage)
	_, err := api.call(definition, errorPageUpdate)
	if err != nil {
		return nil, err
	}

	return errorPage, nil
}

//
// DeleteErrorPage deletes the passed error page using the MYRA API
//
func (api *API) DeleteErrorPage(errorPage *ErrorPage, domainId int) (*ErrorPage, error) {
	if _, ok := methods["deleteErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteErrorPage")
	}

	definition := methods["deleteErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	errorPageUpdate := convertErrorPageToErrorPageUpdate(errorPage)
	_, err := api.call(definition, errorPageUpdate)
	if err != nil {
		return nil, err
	}
	return errorPage, nil
}

//
// decodeErrorPageResponse handles an empty response as it is returned by save error codes
//
func decodeErrorPageResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	return nil, nil
}

//
// convertErrorPageToErrorPageUpdate
//
func convertErrorPageToErrorPageUpdate(errorPage *ErrorPage) *errorPageUpdate {
	errorCode := map[int]bool{
		errorPage.ErrorCode: true,
	}
	selection := map[string]map[int]bool{
		errorPage.SubDomainName: errorCode,
	}

	return &errorPageUpdate{
		PageContent: errorPage.Content,
		Selection:   selection,
	}
}
