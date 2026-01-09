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

// IPFilter ...
type IPFilter struct {
	ID            int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete an IP filter it is necessary to add this attribute to your object."`
	Created       *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new cache setting object. This value is only informational so it is not necessary to add this an attribute to any API call."`
	Modified      *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add the modified timestamp for updates and deletes. This value is always a date type with an ISO 8601 format."`
	ExpireDate    *types.DateTime `json:"expireDate,omitempty" jsonschema:"Expire date schedules the deaktivation of the IP filter. If none is set, the filter will be active until manual deactivation."`
	Value         string          `json:"value" jsonschema:"The value of an IP filter rule can contain a single IP address or a CIDR notation. IPv4 and IPv6 are both supported. An IP filter for IPv6 can only contain a /128 subnet."`
	Type          string          `json:"type" jsonschema:"This specifies how the IP filter is applied. Valid values are BLACKLIST, WHITELIST or WHITELIST_REQUEST_LIMITER."`
	Comment       string          `json:"comment,omitempty" jsonschema:"A comment to describe this IP filter."`
	Enabled       bool            `json:"enabled" jsonschema:"Enable or disable an IP filter."`
	SubDomainName string          `json:"subDomainName" jsonschema:"Identifies the subdomain via a FQDN (Full Qualified Domain Name) where this IP filter belongs to. This value cannot be changed through the object’s attribute as it is set via URL parameter."`
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

	return result.(*IPFilter), nil
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

	return *result.(*[]IPFilter), nil
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
	return result.(*IPFilter), nil
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
	return result.(*IPFilter), nil
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
