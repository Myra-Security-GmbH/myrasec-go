package myrasec

import (
	"context"
	"fmt"
	"net/http"
)

// getTagCacheSettingMethods returns Tag cache setting related API calls
func getTagCacheSettingMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listTagCacheSettings": {
			Name:   "listTagCacheSettings",
			Action: "tag/%d/cache-settings",
			Method: http.MethodGet,
			Result: []CacheSetting{},
		},
		"createTagCacheSetting": {
			Name:   "createTagCacheSetting",
			Action: "tag/%d/cache-settings",
			Method: http.MethodPost,
			Result: CacheSetting{},
		},
		"updateTagCacheSetting": {
			Name:   "updateTagCacheSetting",
			Action: "tag/%d/cache-settings/%d",
			Method: http.MethodPut,
			Result: CacheSetting{},
		},
		"deleteTagCacheSetting": {
			Name:   "deleteTagCacheSetting",
			Action: "tag/%d/cache-settings/%d",
			Method: http.MethodDelete,
			Result: CacheSetting{},
		},
	}
}

// ListTagCacheSettingsContext returns a slice containing all visible cache settings for a subdomain
func (api *API) ListTagCacheSettingsContext(ctx context.Context, tagId int, params map[string]string) ([]CacheSetting, error) {
	if _, ok := api.methods["listTagCacheSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagCacheSettings")
	}

	definition := api.methods["listTagCacheSettings"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]CacheSetting)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListTagCacheSettings is equivalent to ListTagCacheSettingsContext with context.Background().
//
// Deprecated: use ListTagCacheSettingsContext.
func (api *API) ListTagCacheSettings(tagId int, params map[string]string) ([]CacheSetting, error) {
	return api.ListTagCacheSettingsContext(context.Background(), tagId, params)
}

// CreateTagCacheSettingContext creates a new cache setting for the passed subdomain (name) using the MYRA API
func (api *API) CreateTagCacheSettingContext(ctx context.Context, setting *CacheSetting, tagId int) (*CacheSetting, error) {
	if _, ok := api.methods["createTagCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagCacheSetting")
	}

	definition := api.methods["createTagCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(ctx, definition, setting)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*CacheSetting)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateTagCacheSetting is equivalent to CreateTagCacheSettingContext with context.Background().
//
// Deprecated: use CreateTagCacheSettingContext.
func (api *API) CreateTagCacheSetting(setting *CacheSetting, tagId int) (*CacheSetting, error) {
	return api.CreateTagCacheSettingContext(context.Background(), setting, tagId)
}

// UpdateTagCacheSettingContext updates the passed cache setting using the MYRA API
func (api *API) UpdateTagCacheSettingContext(ctx context.Context, setting *CacheSetting, tagId int) (*CacheSetting, error) {
	if _, ok := api.methods["updateTagCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagCacheSetting")
	}

	definition := api.methods["updateTagCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, setting.ID)

	result, err := api.call(ctx, definition, setting)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*CacheSetting)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateTagCacheSetting is equivalent to UpdateTagCacheSettingContext with context.Background().
//
// Deprecated: use UpdateTagCacheSettingContext.
func (api *API) UpdateTagCacheSetting(setting *CacheSetting, tagId int) (*CacheSetting, error) {
	return api.UpdateTagCacheSettingContext(context.Background(), setting, tagId)
}

// DeleteTagCacheSettingContext deletes the passed cache setting using the MYRA API
func (api *API) DeleteTagCacheSettingContext(ctx context.Context, setting *CacheSetting, tagId int) (*CacheSetting, error) {
	if _, ok := api.methods["deleteTagCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTagCacheSetting")
	}

	definition := api.methods["deleteTagCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, setting.ID)

	_, err := api.call(ctx, definition, setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

// DeleteTagCacheSetting is equivalent to DeleteTagCacheSettingContext with context.Background().
//
// Deprecated: use DeleteTagCacheSettingContext.
func (api *API) DeleteTagCacheSetting(setting *CacheSetting, tagId int) (*CacheSetting, error) {
	return api.DeleteTagCacheSettingContext(context.Background(), setting, tagId)
}
