package myrasec

import "net/http"

const (
	ParamPage     = "page"
	ParamPageSize = "pageSize"
	ParamSearch   = "search"
)

//
// APIMethod represents API call definitions used in the methods map
//
type APIMethod struct {
	BaseURL            string
	Name               string
	Action             string
	Method             string
	Result             interface{}
	ResponseDecodeFunc func(resp *http.Response, definition APIMethod) (interface{}, error)
}

//
// methods stores the available APIMethods
//
var methods = map[string]APIMethod{}

//
// initializeMethods ...
//
func initializeMethods() {
	for _, m := range []map[string]APIMethod{
		getBucketMethods(),
		getCacheClearMethods(),
		getCacheSettingMethods(),
		getDNSRecordMethods(),
		getDomainMethods(),
		getErrorPageMethods(),
		getIPFilterMethods(),
		getIPRangeMethods(),
		getMaintenanceMethods(),
		getMaintenanceTemplateMethods(),
		getRateLimitMethods(),
		getRedirectMethods(),
		getSettingsMethods(),
		getSSLMethods(),
		getTagCacheSettingMethods(),
		getTagRateLimitMethods(),
		getTagSettingsMethods(),
		getTagWAFRuleMethods(),
		getTagMethods(),
		getVHostMethods(),
		getWAFMethods(),
	} {
		for k, v := range m {
			methods[k] = v
		}
	}
}
