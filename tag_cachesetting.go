package myrasec

import (
	"fmt"
)

//
// ListTagCacheSettings returns a slice containing all visible cache settings for a subdomain
//
func (api *API) ListTagCacheSettings(tagId int, params map[string]string) ([]CacheSetting, error) {
	if _, ok := methods["listTagCacheSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagCacheSettings")
	}

	definition := methods["listTagCacheSettings"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []CacheSetting
	records = append(records, *result.(*[]CacheSetting)...)

	return records, nil
}

//
// CreateTagCacheSetting creates a new cache setting for the passed subdomain (name) using the MYRA API
//
func (api *API) CreateTagCacheSetting(setting *CacheSetting, tagId int) (*CacheSetting, error) {
	if _, ok := methods["createTagCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagCacheSetting")
	}

	definition := methods["createTagCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return result.(*CacheSetting), nil
}

//
// UpdateTagCacheSetting updates the passed cache setting using the MYRA API
//
func (api *API) UpdateTagCacheSetting(setting *CacheSetting, tagId int) (*CacheSetting, error) {
	if _, ok := methods["updateTagCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagCacheSetting")
	}

	definition := methods["updateTagCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, setting.ID)

	result, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return result.(*CacheSetting), nil
}

//
// DeleteTagCacheSetting deletes the passed cache setting using the MYRA API
//
func (api *API) DeleteTagCacheSetting(setting *CacheSetting, tagId int) (*CacheSetting, error) {
	if _, ok := methods["deleteTagCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTagCacheSetting")
	}

	definition := methods["deleteTagCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, setting.ID)

	_, err := api.call(definition, setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}
