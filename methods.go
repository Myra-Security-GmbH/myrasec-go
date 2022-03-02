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
		Action: "domains",
		Method: http.MethodGet,
		Result: []Domain{},
	},
	"createDomain": {
		Name:   "createDomain",
		Action: "domains",
		Method: http.MethodPost,
		Result: Domain{},
	},
	"updateDomain": {
		Name:   "updateDomain",
		Action: "domains/%d",
		Method: http.MethodPut,
		Result: Domain{},
	},
	"deleteDomain": {
		Name:   "deleteDomain",
		Action: "domains/%d",
		Method: http.MethodDelete,
		Result: Domain{},
	},

	// DNS Record related API calls
	"listDNSRecords": {
		Name:   "listDNSRecords",
		Action: "domain/%d/dns-records",
		Method: http.MethodGet,
		Result: []DNSRecord{},
	},
	"createDNSRecord": {
		Name:   "createDNSRecord",
		Action: "domain/%d/dns-records",
		Method: http.MethodPost,
		Result: DNSRecord{},
	},
	"updateDNSRecord": {
		Name:   "updateDNSRecord",
		Action: "domain/%d/dns-records/%d",
		Method: http.MethodPut,
		Result: DNSRecord{},
	},
	"deleteDNSRecord": {
		Name:   "deleteDNSRecord",
		Action: "domain/%d/dns-records/%d",
		Method: http.MethodDelete,
		Result: DNSRecord{},
	},

	// Cache Setting related API calls
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

	// Redirect related API calls
	"listRedirects": {
		Name:   "listRedirects",
		Action: "domain/%d/redirects/%s",
		Method: http.MethodGet,
		Result: []Redirect{},
	},
	"createRedirect": {
		Name:   "createRedirect",
		Action: "domain/%d/redirects/%s",
		Method: http.MethodPost,
		Result: Redirect{},
	},
	"updateRedirect": {
		Name:   "updateRedirect",
		Action: "domain/%d/redirects/%s/%d",
		Method: http.MethodPut,
		Result: Redirect{},
	},
	"deleteRedirect": {
		Name:   "deleteRedirect",
		Action: "domain/%d/redirects/%s/%d",
		Method: http.MethodDelete,
		Result: Redirect{},
	},

	// Settings related API calls
	"listSettings": {
		Name:               "listSettings",
		Action:             "domain/%d/%s/settings?flat",
		Method:             http.MethodGet,
		Result:             Settings{},
		ResponseDecodeFunc: decodeSettingsResponse,
	},
	"updateSettings": {
		Name:   "updateSettings",
		Action: "domain/%d/%s/settings",
		Method: http.MethodPost,
		Result: Settings{},
	},

	//IP Filter related API calls
	"listIPFilters": {
		Name:   "listIPFilters",
		Action: "domain/%d/ip-filters/%s",
		Method: http.MethodGet,
		Result: []IPFilter{},
	},
	"createIPFilter": {
		Name:   "createIPFilter",
		Action: "domain/%d/ip-filters/%s",
		Method: http.MethodPost,
		Result: IPFilter{},
	},
	"updateIPFilter": {
		Name:   "updateIPFilter",
		Action: "domain/%d/ip-filters/%s/%d",
		Method: http.MethodPut,
		Result: IPFilter{},
	},
	"deleteIPFilter": {
		Name:   "deleteIPFilter",
		Action: "domain/%d/ip-filters/%s/%d",
		Method: http.MethodDelete,
		Result: IPFilter{},
	},

	//Rate Limit related API calls
	"listRateLimits": {
		Name:   "listRateLimits",
		Action: "domain/%d/%s/ratelimits",
		Method: http.MethodGet,
		Result: []RateLimit{},
	},
	"createRateLimit": {
		Name:   "createRateLimit",
		Action: "domain/%d/%s/ratelimits",
		Method: http.MethodPost,
		Result: RateLimit{},
	},
	"updateRateLimit": {
		Name:   "updateRateLimit",
		Action: "domain/%d/%s/ratelimits/%d",
		Method: http.MethodPut,
		Result: RateLimit{},
	},
	"deleteRateLimit": {
		Name:   "deleteRateLimit",
		Action: "domain/%d/%s/ratelimits/%d",
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
		Action: "domain/%d/waf-rules",
		Method: http.MethodGet,
		Result: []WAFRule{},
	},
	"fetchWAFRule": {
		Name:   "fetchWAFRule",
		Action: "domain/waf-rules/%d",
		Method: http.MethodGet,
		Result: []WAFRule{},
	},
	"createWAFRule": {
		Name:   "createWAFRule",
		Action: "domain/%d/%s/waf-rules",
		Method: http.MethodPost,
		Result: WAFRule{},
	},
	"updateWAFRule": {
		Name:   "updateWAFRule",
		Action: "domain/%d/%s/waf-rules/%d",
		Method: http.MethodPut,
		Result: WAFRule{},
	},
	"deleteWAFRule": {
		Name:   "deleteWAFRule",
		Action: "domain/waf-rules/%d",
		Method: http.MethodDelete,
		Result: WAFRule{},
	},

	//SSL certificate related API calls
	"listSSLCertificates": {
		Name:   "listSSLCertificates",
		Action: "domain/%d/ssl/certificates",
		Method: http.MethodGet,
		Result: []SSLCertificate{},
	},
	"createSSLCertificate": {
		Name:   "createSSLCertificate",
		Action: "domain/%d/certificates",
		Method: http.MethodPost,
		Result: SSLCertificate{},
	},
	"updateSSLCertificate": {
		Name:   "updateSSLCertificate",
		Action: "domain/%d/certificates/%d",
		Method: http.MethodPut,
		Result: SSLCertificate{},
	},
	"deleteSSLCertificate": {
		Name:   "deleteSSLCertificate",
		Action: "domain/%d/certificates/%d",
		Method: http.MethodPut,
		Result: SSLCertificate{},
	},

	//IP range related API calls
	"listIPRanges": {
		Name:   "listIPRanges",
		Action: "ip-ranges",
		Method: http.MethodGet,
		Result: []IPRange{},
	},
}
