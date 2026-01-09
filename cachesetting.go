package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getCacheSettingMethods returns Cache Setting related API calls
func getCacheSettingMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listCacheSettings": {
			Name:   "listCacheSettings",
			Action: "domain/%d/%s/cache-settings",
			Method: http.MethodGet,
			Result: []CacheSetting{},
		},
		"createCacheSetting": {
			Name:   "createCacheSetting",
			Action: "domain/%d/%s/cache-settings",
			Method: http.MethodPost,
			Result: CacheSetting{},
		},
		"updateCacheSetting": {
			Name:   "updateCacheSetting",
			Action: "domain/%d/%s/cache-settings/%d",
			Method: http.MethodPut,
			Result: CacheSetting{},
		},
		"deleteCacheSetting": {
			Name:   "deleteCacheSetting",
			Action: "domain/%d/%s/cache-settings/%d",
			Method: http.MethodDelete,
			Result: CacheSetting{},
		},
	}
}

// CacheSetting ...
type CacheSetting struct {
	ID          int             `json:"id,omitempty" jsonschema:"Id is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a cache setting it is necessary to add this attribute to your object."`
	Created     *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new cache setting object. This value is only informational so it is not necessary to add this an attribute to any API call."`
	Modified    *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add the modified timestamp for updates and deletes. This value is always a date type with an ISO 8601 format."`
	Type        string          `json:"type" jsonschema:"This will tell the server how to match the given path against a request. Available options are ’prefix’, ’suffix’ and ’exact’."`
	Path        string          `json:"path" jsonschema:"A request will be matched against this path to decide if this request is cacheable or not. It is possible to write a regexp in this attribute. It is not allowed to use start ’ˆ’ or end ’$’ regexp characters as it they are generated depending on the given type."`
	TTL         int             `json:"ttl" jsonschema:"Time to live limits the lifespan of a cached response for the given path. This is a numeric value and is given in seconds. Special case is ’like origin server’, which uses the TTL returned by your origin server."`
	NotFoundTTL int             `json:"notFoundTtl" jsonschema:"How long an object will be cached. Origin responses with the HTTP codes 404 will be cached."`
	Sort        int             `json:"sort,omitempty" jsonschema:"The order in which the cache rules take action as long as the cache sorting is activated."`
	Enabled     bool            `json:"enabled,omitempty" jsonschema:"Define wether this cache setting is enabled or not."`
	Enforce     bool            `json:"enforce,omitempty" jsonschema:"Enforce cache TTL allows you to set the cache TTL (Cache Control: max-age) in the backend regardless of the response sent from your Origin."`
	Comment     string          `json:"comment,omitempty" jsonschema:"A comment to describe this cache setting."`
}

// ListCacheSettings returns a slice containing all visible cache settings for a subdomain
func (api *API) ListCacheSettings(domainId int, subDomainName string, params map[string]string) ([]CacheSetting, error) {
	if _, ok := methods["listCacheSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listCacheSettings")
	}

	definition := methods["listCacheSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]CacheSetting), nil
}

// CreateCacheSetting creates a new cache setting for the passed subdomain (name) using the MYRA API
func (api *API) CreateCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := methods["createCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createCacheSetting")
	}

	definition := methods["createCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return result.(*CacheSetting), nil
}

// UpdateCacheSetting updates the passed cache setting using the MYRA API
func (api *API) UpdateCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := methods["updateCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateCacheSetting")
	}

	definition := methods["updateCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, setting.ID)

	result, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return result.(*CacheSetting), nil
}

// DeleteCacheSetting deletes the passed cache setting using the MYRA API
func (api *API) DeleteCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := methods["deleteCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteCacheSetting")
	}

	definition := methods["deleteCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, setting.ID)

	_, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}
