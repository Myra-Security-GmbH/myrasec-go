package myrasec

import (
	"fmt"
	"net/http"
)

// getCacheClearMethods returns cache clear related API calls
func getCacheClearMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"clearCache": {
			Name:   "clearCache",
			Action: "domain/%d/cache-clear",
			Method: http.MethodPut,
			Result: CacheClear{},
		},
	}
}

// CacheClear ...
type CacheClear struct {
	FQDN      string `json:"fqdn"`
	Resource  string `json:"resource"`
	Recursive bool   `json:"recursive"`
}

// ClearCache ...
func (api *API) ClearCache(cacheClear *CacheClear, domainId int) (*CacheClear, error) {
	if _, ok := methods["clearCache"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "clearCache")
	}

	definition := methods["clearCache"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, cacheClear)
	if err != nil {
		return nil, err
	}
	return result.(*CacheClear), nil
}
