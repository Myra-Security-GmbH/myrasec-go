package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getIPFilterMethods returns IP Filter related API calls
func getIPFilterMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getIPFilter": {
			Name:               "getIPFilter",
			Action:             "domain/%d/ip-filters/%s/%d",
			Method:             http.MethodGet,
			Result:             IPFilter{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listIPFilters": {
			Name:   "listIPFilters",
			Action: "domain/%d/ip-filters/%s",
			Method: http.MethodGet,
			Result: []IPFilter{},
		},
		"createIPFilter": {
			Name:   "createIPFilter",
			Action: "domain/%d/ip-filters/%s",
			Method: http.MethodPost,
			Result: IPFilter{},
		},
		"updateIPFilter": {
			Name:   "updateIPFilter",
			Action: "domain/%d/ip-filters/%s/%d",
			Method: http.MethodPut,
			Result: IPFilter{},
		},
		"deleteIPFilter": {
			Name:   "deleteIPFilter",
			Action: "domain/%d/ip-filters/%s/%d",
			Method: http.MethodDelete,
			Result: IPFilter{},
		},
	}
}

// IPFilter represents an access control rule for a specific subdomain.
// It allows blocking or whitelisting of traffic based on IP addresses or CIDR ranges.
type IPFilter struct {
	// ID is the unique identifier for the IP filter.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the IP filter. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the filter was added.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// ExpireDate schedules the automatic deactivation of the filter.
	// If nil, the filter remains active until manually disabled or deleted.
	ExpireDate *types.DateTime `json:"expireDate,omitempty" jsonschema:"The timestamp (ISO 8601) when the filter will be automatically deactivated. If null, the filter remains active indefinitely."`

	// Value is the IP address or CIDR range to match.
	// Supports IPv4 and IPv6. Note: IPv6 ranges are restricted to /128.
	Value string `json:"value" jsonschema:"The IP address or CIDR notation (e.g., '1.2.3.4' or '10.0.0.0/24'). Supports IPv4 and IPv6. Note: IPv6 is limited to /128 subnets."`

	// Type defines the action to take when the filter matches.
	// Valid values are BLACKLIST, WHITELIST or WHITELIST_REQUEST_LIMITER.
	Type string `json:"type" jsonschema:"The action type for the filter. Valid values: 'BLACKLIST', 'WHITELIST', 'WHITELIST_REQUEST_LIMITER'."`

	// Comment provides a descriptive note for the filter rule.
	Comment string `json:"comment,omitempty" jsonschema:"A descriptive comment or note for this IP filter."`

	// Enabled controls whether the filter is currently active.
	Enabled bool `json:"enabled" jsonschema:"Indicates if the IP filter rule is currently active (enabled) or inactive."`

	// SubDomainName is the FQDN of the subdomain this filter applies to.
	// This value is typically set via the URL context and is immutable on the object itself.
	SubDomainName string `json:"subDomainName" jsonschema:"The FQDN of the subdomain this filter belongs to. Immutable; usually inferred from the URL parameter or set once during creation."`
}

// GetIPFilter returns a single ip filter with/for the given identifier
func (api *API) GetIPFilter(domainId int, subDomainName string, id int) (*IPFilter, error) {
	if _, ok := methods["getIPFilter"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getIPFilter")
	}

	definition := methods["getIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*IPFilter)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// ListIPFilters returns a slice containing all visible ip filters for a subdomain
func (api *API) ListIPFilters(domainId int, subDomainName string, params map[string]string) ([]IPFilter, error) {
	if _, ok := methods["listIPFilters"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listIPFilters")
	}

	definition := methods["listIPFilters"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]IPFilter)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateIPFilter creates a new ip filter for the passed subdomain (name) using the MYRA API
func (api *API) CreateIPFilter(filter *IPFilter, domainId int, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["createIPFilter"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createIPFilter")
	}

	definition := methods["createIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*IPFilter)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateIPFilter updates the passed ip filter using the MYRA API
func (api *API) UpdateIPFilter(filter *IPFilter, domainId int, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["updateIPFilter"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateIPFilter")
	}

	definition := methods["updateIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, filter.ID)

	result, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*IPFilter)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteIPFilter deletes the passed ip filter using the MYRA API
func (api *API) DeleteIPFilter(filter *IPFilter, domainId int, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["deleteIPFilter"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteIPFilter")
	}

	definition := methods["deleteIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, filter.ID)

	_, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}
