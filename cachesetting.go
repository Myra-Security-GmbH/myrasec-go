package myrasec

import (
	"context"
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

// CacheSetting represents a single rule definition for caching behavior.
// It determines how requests are matched (via path patterns) and how long
// responses are retained (TTL).
type CacheSetting struct {
	// ID is the unique identifier for the cache setting.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the cache setting. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the setting was established.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). This is a server-managed, read-only value."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Type defines the matching strategy for the path.
	// Valid options are "prefix", "suffix" and "exact".
	Type string `json:"type" jsonschema:"The strategy used to match the request path. Allowed values: 'prefix', 'suffix', 'exact'."`

	// Path is the pattern used to identify requests for this rule.
	// It supports regular expressions but must not contain start ('^') or end ('$') anchors,
	// as these are implied by the chosen Type.
	Path string `json:"path" jsonschema:"The URI path or regex pattern to match. Note: Do not use start ('^') or end ('$') anchors, as they are automatically handled based on the selected 'type'."`

	// TTL (Time To Live) defines the cache duration in seconds.
	TTL int `json:"ttl" jsonschema:"The Time To Live (TTL) in seconds. Defines the lifespan of the cached response."`

	// NotFoundTTL defines the cache duration for 404 responses.
	NotFoundTTL int `json:"notFoundTtl" jsonschema:"The duration in seconds to cache 404 (Not Found) responses."`

	// Sort controls the priority of the rule execution.
	// Lower numbers typically indicate higher priority when sorting is active.
	Sort int `json:"sort,omitempty" jsonschema:"The execution priority order of the cache rule."`

	// Enabled determines if the rule is currently active.
	Enabled bool `json:"enabled,omitempty" jsonschema:"Determines whether this cache setting is active."`

	// Enforce overrides the Origin's Cache-Control headers.
	// If true, the backend uses the defined TTL regardless of the Origin's instructions.
	Enforce bool `json:"enforce,omitempty" jsonschema:"If true, overrides the Origin server's Cache-Control headers and enforces the configured TTL."`

	// Comment is an optional text to describe the purpose of this rule.
	Comment string `json:"comment,omitempty" jsonschema:"An optional comment or description for this cache setting."`
}

// ListCacheSettingsContext returns a slice containing all visible cache settings for a subdomain
func (api *API) ListCacheSettingsContext(ctx context.Context, domainId int, subDomainName string, params map[string]string) ([]CacheSetting, error) {
	if _, ok := api.methods["listCacheSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listCacheSettings")
	}

	definition := api.methods["listCacheSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

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

// ListCacheSettings is equivalent to ListCacheSettingsContext with context.Background().
//
// Deprecated: use ListCacheSettingsContext.
func (api *API) ListCacheSettings(domainId int, subDomainName string, params map[string]string) ([]CacheSetting, error) {
	return api.ListCacheSettingsContext(context.Background(), domainId, subDomainName, params)
}

// CreateCacheSettingContext creates a new cache setting for the passed subdomain (name) using the MYRA API
func (api *API) CreateCacheSettingContext(ctx context.Context, setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := api.methods["createCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createCacheSetting")
	}

	definition := api.methods["createCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

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

// CreateCacheSetting is equivalent to CreateCacheSettingContext with context.Background().
//
// Deprecated: use CreateCacheSettingContext.
func (api *API) CreateCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	return api.CreateCacheSettingContext(context.Background(), setting, domainId, subDomainName)
}

// UpdateCacheSettingContext updates the passed cache setting using the MYRA API
func (api *API) UpdateCacheSettingContext(ctx context.Context, setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := api.methods["updateCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateCacheSetting")
	}

	definition := api.methods["updateCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, setting.ID)

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

// UpdateCacheSetting is equivalent to UpdateCacheSettingContext with context.Background().
//
// Deprecated: use UpdateCacheSettingContext.
func (api *API) UpdateCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	return api.UpdateCacheSettingContext(context.Background(), setting, domainId, subDomainName)
}

// DeleteCacheSettingContext deletes the passed cache setting using the MYRA API
func (api *API) DeleteCacheSettingContext(ctx context.Context, setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	if _, ok := api.methods["deleteCacheSetting"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteCacheSetting")
	}

	definition := api.methods["deleteCacheSetting"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, setting.ID)

	_, err := api.call(ctx, definition, setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

// DeleteCacheSetting is equivalent to DeleteCacheSettingContext with context.Background().
//
// Deprecated: use DeleteCacheSettingContext.
func (api *API) DeleteCacheSetting(setting *CacheSetting, domainId int, subDomainName string) (*CacheSetting, error) {
	return api.DeleteCacheSettingContext(context.Background(), setting, domainId, subDomainName)
}
