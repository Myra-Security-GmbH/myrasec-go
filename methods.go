package myrasec

import "net/http"

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
}
