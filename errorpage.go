package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getErrorPageMethods returns Error Page related API calls
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

// errorPageUpdate
type errorPageUpdate struct {
	ID          int                     `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete an error page, the ID does not help because it is not used."`
	PageContent string                  `json:"pageContent,omitempty"`
	Selection   map[string]map[int]bool `json:"selection,omitempty"`
	Created     *types.DateTime         `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new error page object. This value is only informational so it is not neccessary to add this attribute to any API call."`
	Modified    *types.DateTime         `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updateing the most recent version and not overwriting other changes, you always have to add the modified timestamp for updates and deletes. This value is always a date type with an Identifies the version of the object. To ensure that you are updateing the most recent version and not overwriting other changes, you always have to add the modified timestamp for updates and deletes. This value is always a date type with an ISO 8601 format.ISO 8601 format."`
}

// ErrorPage
type ErrorPage struct {
	ID            int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete an error page, the ID does not help because it is not used."`
	Created       *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new error page object. This value is only informational so it is not neccessary to add this attribute to any API call."`
	Modified      *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updateing the most recent version and not overwriting other changes, you always have to add the modified timestamp for updates and deletes. This value is always a date type with an ISO 8601 format."`
	ErrorCode     int             `json:"errorCode,omitempty" jsonschema:"ErrorCode represents the Http Code for this error page. The ErrorCode should never be changed, because it is used as identifier additionally with the SubDomainName."`
	Content       string          `json:"content,omitempty" jsonschema:"The Content is the HTML code for the error page."`
	SubDomainName string          `json:"subDomainName,omitempty" jsonschema:"The configured error page is available for this subdomain. The SubDomainName should never be changed, becuase it is used as identifier additionally with the ErrorCode."`
}

// GetErrorPage returns a single error page with/for the given identifier
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

// ListErrorPages returns a slice containing all error pages
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

	return *result.(*[]ErrorPage), nil
}

// CreateErrorPage creates a new error page using the MYRA API
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

// UpdateErrorPage updates the passed error page using the MYRA API
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

// DeleteErrorPage deletes the passed error page using the MYRA API
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

// decodeErrorPageResponse handles an empty response as it is returned by save error codes
func decodeErrorPageResponse(resp *http.Response, definition APIMethod) (any, error) {
	return nil, nil
}

// convertErrorPageToErrorPageUpdate
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
