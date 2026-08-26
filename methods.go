package myrasec

import (
	"maps"
	"net/http"
)

const (
	ParamPage     = "page"
	ParamPageSize = "pageSize"
	ParamSearch   = "search"
)

// APIMethod represents API call definitions used in the methods map
type APIMethod struct {
	BaseURL            string
	Name               string
	Action             string
	Method             string
	Result             any
	AdditionalHeaders  map[string]string
	ResponseDecodeFunc func(resp *http.Response, definition APIMethod) (any, error)
}

// initializeMethods populates the methods registry with all available API method definitions.
func initializeMethods() map[string]APIMethod {
	methods := map[string]APIMethod{}

	for _, m := range []map[string]APIMethod{
		getAPIKeyMethods(),
		getCacheClearMethods(),
		getCacheSettingMethods(),
		getDNSRecordMethods(),
		getDomainMethods(),
		getErrorPageMethods(),
		getIPFilterMethods(),
		getIPRangeMethods(),
		getMaintenanceMethods(),
		getMaintenanceTemplateMethods(),
		getPermissionMethods(),
		getRedirectMethods(),
		getSettingsMethods(),
		getSslConfigurationMethods(),
		getSSLMethods(),
		getStatisticsMethods(),
		getTagCacheSettingMethods(),
		getTagSettingsMethods(),
		getTagWAFRuleMethods(),
		getTagMethods(),
		getUserGroupMethods(),
		getUserMethods(),
		getVHostMethods(),
		getWAFMethods(),
		getWaitingRoomMethods(),
		getTagInformationMethods(),
		getZoneConfigMethods(),
	} {
		maps.Copy(methods, m)
	}

	return methods
}
