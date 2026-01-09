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

// CacheClear ...
type CacheClear struct {
	FQDN      string `json:"fqdn,omitempty" jsonschema:"The specific name of the domain or subdomain. If the attribute is not set, the whole cache of the domain will be cleared."`
	Resource  string `json:"resource" jsonschema:"A specific path to be cache-cleared. Wildcards can be used. (The path is relative to the root of the website)."`
	Recursive bool   `json:"recursive" jsonschema:"A boolean value to define if the cache purge should be done recursively."`
}

// ClearCache ...
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

	return result.(*[]CacheClear), nil
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
	decoder.Decode(returnResult)

	return returnResult, err
}
