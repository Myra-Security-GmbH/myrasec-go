package myrasec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// getCacheClearMethods returns cache clear related API calls
func getCacheClearMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"clearCache": {
			Name:               "clearCache",
			Action:             "domain/%d/cache-clear",
			Method:             http.MethodPut,
			Result:             []CacheClear{},
			ResponseDecodeFunc: decodeCacheClearResponse,
		},
	}
}

// CacheClear defines the parameters for invalidating cached content.
// It allows targeting specific resources or performing bulk purges based on
// domain and path patterns.
type CacheClear struct {
	// FQDN specifies the Fully Qualified Domain Name (e.g., "www.example.com").
	// If left empty, the purge operation applies to the entire domain scope.
	FQDN string `json:"fqdn,omitempty" jsonschema:"The Fully Qualified Domain Name (e.g., 'www.example.com'). If omitted, the entire domain cache is cleared."`

	// Resource indicates the relative path to invalidate.
	// It supports wildcard characters (e.g., "*") for pattern matching.
	Resource string `json:"resource" jsonschema:"The relative path or pattern to purge (e.g., '/images/*' or '/index.html'). Supports wildcards and is relative to the website root."`

	// Recursive determines whether the purge applies to sub-directories.
	// If true, the operation extends to all resources nested under the target path.
	Recursive bool `json:"recursive" jsonschema:"Enables recursive purging. If true, the operation clears the target resource and all nested sub-resources."`
}

// ClearCache triggers a cache clear operation for the given domain.
func (api *API) ClearCache(cacheClear *CacheClear, domainId int) (*[]CacheClear, error) {
	if _, ok := methods["clearCache"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "clearCache")
	}

	definition := methods["clearCache"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, cacheClear)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]CacheClear)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// decodeCacheClearResponse - custom decode function for cache clear responses
func decodeCacheClearResponse(resp *http.Response, definition APIMethod) (any, error) {
	res, err := decodeBaseResponse(resp)
	if err != nil {
		return nil, err
	}
	result := res.Data

	tmp, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if definition.Result == nil {
		return tmp, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(tmp))
	retValue := reflect.New(reflect.TypeOf(definition.Result))
	returnResult := retValue.Interface()
	err = decoder.Decode(returnResult)

	return returnResult, err
}
