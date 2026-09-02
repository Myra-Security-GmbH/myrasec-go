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
	ID          int                     `json:"id,omitempty"`
	PageContent string                  `json:"pageContent,omitempty"`
	Selection   map[string]map[int]bool `json:"selection,omitempty"`
	Created     *types.DateTime         `json:"created,omitempty"`
	Modified    *types.DateTime         `json:"modified,omitempty"`
}

// ErrorPage represents a custom HTML error page configuration (e.g., 404, 500).
// Unlike other objects, it is uniquely identified by the combination of SubDomainName
// and ErrorCode, rather than the numeric ID.
type ErrorPage struct {
	// ID is the internal system identifier for the object.
	// This value is read-only and ignored for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"System-assigned numeric identifier. Read-only. This value is ignored during updates or deletes; use ErrorCode and SubDomainName to identify the resource instead."`

	// Created indicates when the error page was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only. Informational only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// ErrorCode represents the HTTP Status Code (e.g., 404, 500).
	// This value is part of the composite unique key and is immutable once created.
	ErrorCode int `json:"errorCode,omitempty" jsonschema:"The HTTP status code (e.g., 404 or 500). Part of the composite unique key. Immutable after creation."`

	// Content contains the raw HTML code to be rendered.
	Content string `json:"content,omitempty" jsonschema:"The raw HTML content to be displayed for this error page."`

	// SubDomainName is the FQDN for which this error page is configured.
	// This value is part of the composite unique key and is immutable once created.
	SubDomainName string `json:"subDomainName,omitempty" jsonschema:"The target subdomain FQDN (e.g., 'shop.example.com'). Part of the composite unique key. Immutable after creation."`
}

// GetErrorPage returns a single error page with/for the given identifier
func (api *API) GetErrorPage(domainId int, pageId int) (*ErrorPage, error) {
	if _, ok := api.methods["getErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getErrorPage")
	}

	definition := api.methods["getErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, pageId)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*ErrorPage)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// ListErrorPages returns a slice containing all error pages
func (api *API) ListErrorPages(domainId int, params map[string]string) ([]ErrorPage, error) {
	if _, ok := api.methods["listErrorPages"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listErrorPages")
	}

	definition := api.methods["listErrorPages"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]ErrorPage)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateErrorPage creates a new error page using the MYRA API
func (api *API) CreateErrorPage(errorPage *ErrorPage, domainId int) (*ErrorPage, error) {
	if _, ok := api.methods["createErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createErrorPage")
	}

	definition := api.methods["createErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	errorPageUpdate := convertErrorPageToErrorPageUpdate(errorPage)
	res, err := api.call(definition, errorPageUpdate)
	if err != nil {
		return nil, err
	}

	if err := enrichErrorPageIdentity(errorPage, res); err != nil {
		return nil, err
	}

	return errorPage, nil
}

// UpdateErrorPage updates the passed error page using the MYRA API
func (api *API) UpdateErrorPage(errorPage *ErrorPage, domainId int) (*ErrorPage, error) {
	if _, ok := api.methods["updateErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateErrorPage")
	}

	definition := api.methods["updateErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	errorPageUpdate := convertErrorPageToErrorPageUpdate(errorPage)
	res, err := api.call(definition, errorPageUpdate)
	if err != nil {
		return nil, err
	}

	if err := enrichErrorPageIdentity(errorPage, res); err != nil {
		return nil, err
	}

	return errorPage, nil
}

// DeleteErrorPage deletes the passed error page using the MYRA API
func (api *API) DeleteErrorPage(errorPage *ErrorPage, domainId int) (*ErrorPage, error) {
	if _, ok := api.methods["deleteErrorPage"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteErrorPage")
	}

	definition := api.methods["deleteErrorPage"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	errorPageUpdate := convertErrorPageToErrorPageUpdate(errorPage)
	_, err := api.call(definition, errorPageUpdate)
	if err != nil {
		return nil, err
	}
	return errorPage, nil
}

// decodeErrorPageResponse decodes the save (create/update) response of an error page.
// A server carrying the companion change echoes the persisted page(s) in the "data"
// envelope, so the first element is decoded into an *ErrorPage (carrying the
// server-assigned id, created and modified).
//
// A server without that change answers with an empty data array ({"data": []}); in
// that case a nil result is returned so the caller keeps the object it sent. The save
// endpoint is a bulk operation, but this client only ever selects a single
// subdomain/errorCode pair, so the response holds at most one element.
func decodeErrorPageResponse(resp *http.Response, definition APIMethod) (any, error) {
	res, err := decodeBaseResponse(resp)
	if err != nil {
		return nil, err
	}

	if res == nil || len(res.Data) == 0 {
		return nil, nil
	}

	return prepareSingleElementResult(*res, definition)
}

// enrichErrorPageIdentity copies the server-assigned identity fields (id, created,
// modified) from a decoded save response onto the error page the caller sent, leaving
// the caller's content/subdomain/code intact. When the save response was empty (res is
// nil, e.g. a server without the companion change), the caller's object is left
// unchanged. An unexpected result type is reported as an error, matching the rest of
// the package.
func enrichErrorPageIdentity(errorPage *ErrorPage, res any) error {
	if res == nil {
		return nil
	}

	created, ok := res.(*ErrorPage)
	if !ok {
		return fmt.Errorf("unexpected result type %T", res)
	}

	errorPage.ID = created.ID
	errorPage.Created = created.Created
	errorPage.Modified = created.Modified

	return nil
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
