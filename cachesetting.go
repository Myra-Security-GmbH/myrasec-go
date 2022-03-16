package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// CacheSetting ...
//
type CacheSetting struct {
	ID          int             `json:"id,omitempty"`
	Created     *types.DateTime `json:"created,omitempty"`
	Modified    *types.DateTime `json:"modified,omitempty"`
	Type        string          `json:"type"`
	Path        string          `json:"path"`
	TTL         int             `json:"ttl"`
	NotFoundTTL int             `json:"notFoundTtl"`
	Sort        int             `json:"sort,omitempty"`
	Enabled     bool            `json:"enabled"`
	Enforce     bool            `json:"enforce"`
}

//
// ListCacheSettings returns a slice containing all visible cache settings for a subdomain
//
func (api *API) ListCacheSettings(domainId int, subDomainName string, params map[string]string) ([]CacheSetting, error) {
	if _, ok := methods["listCacheSettings"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listCacheSettings")
	}

	definition := methods["listCacheSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []CacheSetting
	for _, v := range *result.(*[]CacheSetting) {
		records = append(records, v)
	}

	return records, nil
}

//
// CreateCacheSetting creates a new cache setting for the passed subdomain (name) using the MYRA API
//
func (api *API) CreateCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := methods["createCacheSetting"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "createCacheSetting")
	}

	definition := methods["createCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return result.(*CacheSetting), nil
}

//
// UpdateCacheSetting updates the passed cache setting using the MYRA API
//
func (api *API) UpdateCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := methods["updateCacheSetting"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "updateCacheSetting")
	}

	definition := methods["updateCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, setting.ID)

	result, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return result.(*CacheSetting), nil
}

//
// DeleteCacheSetting deletes the passed cache setting using the MYRA API
//
func (api *API) DeleteCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := methods["deleteCacheSetting"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteCacheSetting")
	}

	definition := methods["deleteCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, setting.ID)

	_, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}
