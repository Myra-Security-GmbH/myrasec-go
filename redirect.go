package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getRedirectMethods returns Redirect related API calls
func getRedirectMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getRedirect": {
			Name:               "getRedirect",
			Action:             "domain/%d/redirects/%s/%d",
			Method:             http.MethodGet,
			Result:             Redirect{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listRedirects": {
			Name:   "listRedirects",
			Action: "domain/%d/redirects/%s",
			Method: http.MethodGet,
			Result: []Redirect{},
		},
		"createRedirect": {
			Name:   "createRedirect",
			Action: "domain/%d/redirects/%s",
			Method: http.MethodPost,
			Result: Redirect{},
		},
		"updateRedirect": {
			Name:   "updateRedirect",
			Action: "domain/%d/redirects/%s/%d",
			Method: http.MethodPut,
			Result: Redirect{},
		},
		"deleteRedirect": {
			Name:   "deleteRedirect",
			Action: "domain/%d/redirects/%s/%d",
			Method: http.MethodDelete,
			Result: Redirect{},
		},
	}
}

// Redirect represents an HTTP redirection rule (Forwarding).
// It maps incoming requests from a source path to a destination URL based on specific matching criteria.
type Redirect struct {
	// ID is the unique identifier for the redirect rule.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the redirect. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the redirect was added.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Type determines the HTTP status code used for the redirection.
	// Valid values are 'permanent' (HTTP 301) and 'redirect' (HTTP 302).
	Type string `json:"type" jsonschema:"The HTTP redirect method. Valid values: 'permanent' (301 - browser cacheable) or 'redirect' (302 - temporary)."`

	// SubDomainName is the FQDN of the subdomain this redirect applies to.
	// This value is typically set via the URL context and is immutable on the object itself.
	SubDomainName string `json:"subDomainName" jsonschema:"The FQDN of the subdomain this redirect belongs to. Immutable; usually inferred from the URL parameter or set once during creation."`

	// Source is the incoming path or pattern to match against.
	// Depending on MatchingType, this can be an exact path, prefix or suffix.
	Source string `json:"source" jsonschema:"The source path or pattern to match against the incoming request URI."`

	// Destination is the target URL where the user will be sent.
	// Can be an absolute URL (https://...) or a relative path (/page).
	Destination string `json:"destination" jsonschema:"The target location. Can be a full URL (e.g., 'https://example.com') or a relative path (e.g., '/contact')."`

	// MatchingType defines how the Source field is interpreted.
	// Valid values are 'prefix', 'suffix' and 'exact'.
	MatchingType string `json:"matchingType" jsonschema:"Defines the matching strategy. Valid values: 'prefix' (starts with), 'suffix' (ends with), 'exact' (precise match)."`

	// Comment provides a descriptive note for the redirect rule.
	Comment string `json:"comment,omitempty" jsonschema:"A descriptive comment or note for this redirect rule."`

	// Sort defines the execution order of redirect rules.
	// Lower numbers are processed first (ascending order).
	Sort int `json:"sort,omitempty" jsonschema:"The priority/sort order (ascending). Lower numbers are processed first."`

	// Enabled controls whether the redirect rule is currently active.
	Enabled bool `json:"enabled" jsonschema:"Indicates if the redirect rule is currently active (enabled) or ignored."`

	// ExpertMode disables safety checks like loop detection.
	// Use with caution to prevent infinite redirect loops.
	ExpertMode bool `json:"expertMode,omitempty" jsonschema:"If true, disables automatic redirect loop detection. Use with caution."`
}

// GetRedirectContext returns a single redirect with/for the given identifier
func (api *API) GetRedirectContext(ctx context.Context, domainId int, subDomainName string, id int) (*Redirect, error) {
	if _, ok := api.methods["getRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getRedirect")
	}

	definition := api.methods["getRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, id)

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*Redirect)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// GetRedirect is equivalent to GetRedirectContext with context.Background().
//
// Deprecated: use GetRedirectContext.
func (api *API) GetRedirect(domainId int, subDomainName string, id int) (*Redirect, error) {
	return api.GetRedirectContext(context.Background(), domainId, subDomainName, id)
}

// ListRedirectsContext returns a slice containing all visible redirects for a subdomain
func (api *API) ListRedirectsContext(ctx context.Context, domainId int, subDomainName string, params map[string]string) ([]Redirect, error) {
	if _, ok := api.methods["listRedirects"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listRedirects")
	}

	definition := api.methods["listRedirects"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]Redirect)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListRedirects is equivalent to ListRedirectsContext with context.Background().
//
// Deprecated: use ListRedirectsContext.
func (api *API) ListRedirects(domainId int, subDomainName string, params map[string]string) ([]Redirect, error) {
	return api.ListRedirectsContext(context.Background(), domainId, subDomainName, params)
}

// CreateRedirectContext creates a new redirect for the passed subdomain (name) using the MYRA API
func (api *API) CreateRedirectContext(ctx context.Context, redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	if _, ok := api.methods["createRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createRedirect")
	}

	definition := api.methods["createRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(ctx, definition, redirect)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Redirect)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateRedirect is equivalent to CreateRedirectContext with context.Background().
//
// Deprecated: use CreateRedirectContext.
func (api *API) CreateRedirect(redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	return api.CreateRedirectContext(context.Background(), redirect, domainId, subDomainName)
}

// UpdateRedirectContext updates the passed redirect using the MYRA API
func (api *API) UpdateRedirectContext(ctx context.Context, redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	if _, ok := api.methods["updateRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateRedirect")
	}

	definition := api.methods["updateRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, redirect.ID)

	result, err := api.call(ctx, definition, redirect)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Redirect)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateRedirect is equivalent to UpdateRedirectContext with context.Background().
//
// Deprecated: use UpdateRedirectContext.
func (api *API) UpdateRedirect(redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	return api.UpdateRedirectContext(context.Background(), redirect, domainId, subDomainName)
}

// DeleteRedirectContext deletes the passed redirect using the MYRA API
func (api *API) DeleteRedirectContext(ctx context.Context, redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	if _, ok := api.methods["deleteRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteRedirect")
	}

	definition := api.methods["deleteRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, redirect.ID)

	_, err := api.call(ctx, definition, redirect)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}

// DeleteRedirect is equivalent to DeleteRedirectContext with context.Background().
//
// Deprecated: use DeleteRedirectContext.
func (api *API) DeleteRedirect(redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	return api.DeleteRedirectContext(context.Background(), redirect, domainId, subDomainName)
}
