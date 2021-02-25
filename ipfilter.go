package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/pkg/types"
)

//
// IPFilter ...
//
type IPFilter struct {
	ID         int             `json:"id,omitempty"`
	Created    *types.DateTime `json:"created,omitempty"`
	Modified   *types.DateTime `json:"modified,omitempty"`
	Value      string          `json:"value"`
	Type       string          `json:"type"`
	ExpireDate *types.DateTime `json:"expireDate,omitempty"`
	Enabled    bool            `json:"enabled,omitempty"`
}

//
// ListIPFilters returns a slice containing all visible ip filters for a subdomain
//
func (api *API) ListIPFilters(subDomainName string, params map[string]string) ([]IPFilter, error) {
	if _, ok := methods["listIPFilters"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listIPFilters")
	}
	definition := methods["listIPFilters"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName, 1)

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
func (api *API) CreateIPFilter(filter *IPFilter, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["createIPFilter"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "createIPFilter")
	}
	definition := methods["createIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	return result.(*IPFilter), nil
}

//
// UpdateIPFilter updates the passed ip filter using the MYRA API
//
func (api *API) UpdateIPFilter(filter *IPFilter, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["updateIPFilter"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "updateIPFilter")
	}
	definition := methods["updateIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	return result.(*IPFilter), nil
}

//
// DeleteIPFilter deletes the passed ip filter using the MYRA API
//
func (api *API) DeleteIPFilter(filter *IPFilter, subDomainName string) (*IPFilter, error) {
	if _, ok := methods["deleteIPFilter"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteIPFilter")
	}
	definition := methods["deleteIPFilter"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, filter)
	if err != nil {
		return nil, err
	}
	return result.(*IPFilter), nil
}
