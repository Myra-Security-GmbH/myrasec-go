package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// IPFilter ...
//
type IPFilter struct {
	ID         int             `json:"id,omitempty"`
	Created    *types.DateTime `json:"created,omitempty"`
	Modified   *types.DateTime `json:"modified,omitempty"`
	ExpireDate *types.DateTime `json:"expireDate,omitempty"`
	Value      string          `json:"value"`
	Type       string          `json:"type"`
	Comment    string          `json:"comment,omitempty"`
	Enabled    bool            `json:"enabled,omitempty"`
}

//
// GetIPFilter returns a single ip filter with/for the given identifier
//
func (api *API) GetIPFilter(domainId int, subDomainName string, id int) (*IPFilter, error) {
	if _, ok := methods["getIPFilter"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "getIPFilter")
	}

	definition := methods["getIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*IPFilter), nil
}

//
// ListIPFilters returns a slice containing all visible ip filters for a subdomain
//
func (api *API) ListIPFilters(domainId int, subDomainName string, params map[string]string) ([]IPFilter, error) {
	if _, ok := methods["listIPFilters"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listIPFilters")
	}

	definition := methods["listIPFilters"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []IPFilter
	for _, v := range *result.(*[]IPFilter) {
		records = append(records, v)
	}

	return records, nil
}

//
// CreateIPFilter creates a new ip filter for the passed subdomain (name) using the MYRA API
//
func (api *API) CreateIPFilter(filter *IPFilter, domainId int, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["createIPFilter"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "createIPFilter")
	}

	definition := methods["createIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	return result.(*IPFilter), nil
}

//
// UpdateIPFilter updates the passed ip filter using the MYRA API
//
func (api *API) UpdateIPFilter(filter *IPFilter, domainId int, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["updateIPFilter"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "updateIPFilter")
	}

	definition := methods["updateIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, filter.ID)

	result, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	return result.(*IPFilter), nil
}

//
// DeleteIPFilter deletes the passed ip filter using the MYRA API
//
func (api *API) DeleteIPFilter(filter *IPFilter, domainId int, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["deleteIPFilter"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteIPFilter")
	}

	definition := methods["deleteIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, filter.ID)

	_, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}
