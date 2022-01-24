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
	Name               string
	Action             string
	Method             string
	Result             interface{}
	ResponseDecodeFunc func(resp *http.Response, definition APIMethod) (interface{}, error)
}

var methods = map[string]APIMethod{
	// Domain related API calls
	"listDomains": {
		Name:   "listDomains",
		Action: "domains/%d",
		Method: http.MethodGet,
		Result: []Domain{},
	},
	"createDomain": {
		Name:   "createDomain",
		Action: "domains",
		Method: http.MethodPut,
		Result: Domain{},
	},
	"updateDomain": {
		Name:   "updateDomain",
		Action: "domains",
		Method: http.MethodPost,
		Result: Domain{},
	},
	"deleteDomain": {
		Name:   "deleteDomain",
		Action: "domains",
		Method: http.MethodDelete,
		Result: Domain{},
	},

	// DNS Record related API calls
	"listDNSRecords": {
		Name:   "listDNSRecords",
		Action: "dnsRecords/%s/%d",
		Method: http.MethodGet,
		Result: []DNSRecord{},
	},
	"createDNSRecord": {
		Name:   "createDNSRecord",
		Action: "dnsRecords/%s",
		Method: http.MethodPut,
		Result: DNSRecord{},
	},
	"updateDNSRecord": {
		Name:   "updateDNSRecord",
		Action: "dnsRecords/%s",
		Method: http.MethodPost,
		Result: DNSRecord{},
	},
	"deleteDNSRecord": {
		Name:   "deleteDNSRecord",
		Action: "dnsRecords/%s",
		Method: http.MethodDelete,
		Result: DNSRecord{},
	},

	// Cache Setting related API calls
	"listCacheSettings": {
		Name:   "listCacheSettings",
		Action: "cacheSettings/%s/%d",
		Method: http.MethodGet,
		Result: []CacheSetting{},
	},
	"createCacheSetting": {
		Name:   "createCacheSetting",
		Action: "cacheSettings/%s",
		Method: http.MethodPut,
		Result: CacheSetting{},
	},
	"updateCacheSetting": {
		Name:   "updateCacheSetting",
		Action: "cacheSettings/%s",
		Method: http.MethodPost,
		Result: CacheSetting{},
	},
	"deleteCacheSetting": {
		Name:   "deleteCacheSetting",
		Action: "cacheSettings/%s",
		Method: http.MethodDelete,
		Result: CacheSetting{},
	},

	// Redirect related API calls
	"listRedirects": {
		Name:   "listRedirects",
		Action: "redirects/%s/%d",
		Method: http.MethodGet,
		Result: []Redirect{},
	},
	"createRedirect": {
		Name:   "createRedirect",
		Action: "redirects/%s",
		Method: http.MethodPut,
		Result: Redirect{},
	},
	"updateRedirect": {
		Name:   "updateRedirect",
		Action: "redirects/%s",
		Method: http.MethodPost,
		Result: Redirect{},
	},
	"deleteRedirect": {
		Name:   "deleteRedirect",
		Action: "redirects/%s",
		Method: http.MethodDelete,
		Result: Redirect{},
	},

	// Settings related API calls
	"listSettings": {
		Name:               "listSettings",
		Action:             "subdomainSetting/%s?flat",
		Method:             http.MethodGet,
		Result:             Settings{},
		ResponseDecodeFunc: decodeSettingsResponse,
	},
	"updateSettings": {
		Name:   "updateSettings",
		Action: "subdomainSetting/%s",
		Method: http.MethodPost,
		Result: Settings{},
	},

	//IP Filter related API calls
	"listIPFilters": {
		Name:   "listIPFilters",
		Action: "ipfilter/%s/%d",
		Method: http.MethodGet,
		Result: []IPFilter{},
	},
	"createIPFilter": {
		Name:   "createIPFilter",
		Action: "ipfilter/%s",
		Method: http.MethodPut,
		Result: IPFilter{},
	},
	"updateIPFilter": {
		Name:   "updateIPFilter",
		Action: "ipfilter/%s",
		Method: http.MethodPost,
		Result: IPFilter{},
	},
	"deleteIPFilter": {
		Name:   "deleteIPFilter",
		Action: "ipfilter/%s",
		Method: http.MethodDelete,
		Result: IPFilter{},
	},

	//Rate Limit related API calls
	"listRateLimits": {
		Name:   "listRateLimits",
		Action: "ratelimit/%s/%d",
		Method: http.MethodGet,
		Result: []RateLimit{},
	},
	"createRateLimit": {
		Name:   "createRateLimit",
		Action: "ratelimit",
		Method: http.MethodPut,
		Result: RateLimit{},
	},
	"updateRateLimit": {
		Name:   "updateRateLimit",
		Action: "ratelimit",
		Method: http.MethodPost,
		Result: RateLimit{},
	},
	"deleteRateLimit": {
		Name:   "deleteRateLimit",
		Action: "ratelimit",
		Method: http.MethodDelete,
		Result: RateLimit{},
	},

	//WAF related API calls
	"listWAFConditions": {
		Name:   "listWAFConditions",
		Action: "waf/conditions",
		Method: http.MethodGet,
		Result: []WAFCondition{},
	},
	"listWAFActions": {
		Name:   "listWAFActions",
		Action: "waf/actions",
		Method: http.MethodGet,
		Result: []WAFAction{},
	},
	"listWAFRules": {
		Name:   "listWAFRules",
		Action: "waf/rules/%s/%d",
		Method: http.MethodGet,
		Result: []WAFRule{},
	},
	"fetchWAFRule": {
		Name:   "fetchWAFRule",
		Action: "waf/rule/%d",
		Method: http.MethodGet,
		Result: []WAFRule{},
	},
	"createWAFRule": {
		Name:   "createWAFRule",
		Action: "waf/rule",
		Method: http.MethodPut,
		Result: WAFRule{},
	},
	"updateWAFRule": {
		Name:   "updateWAFRule",
		Action: "waf/rule",
		Method: http.MethodPost,
		Result: WAFRule{},
	},
	"deleteWAFRule": {
		Name:   "deleteWAFRule",
		Action: "waf/rule",
		Method: http.MethodDelete,
		Result: WAFRule{},
	},
}
