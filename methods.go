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
	// VHost related API calls
	"listAllSubdomains": {
		Name:   "listAllSubdomains",
		Action: "subdomains",
		Method: http.MethodGet,
		Result: []VHost{},
	},
	"listSubdomainsForDomain": {
		Name:   "listSubdomainsForDomain",
		Action: "domain/%d/subdomains",
		Method: http.MethodGet,
		Result: []VHost{},
	},

	// Domain related API calls
	"getDomain": {
		Name:               "getDomain",
		Action:             "domains/%d",
		Method:             http.MethodGet,
		Result:             Domain{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	},
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
	"getDNSRecord": {
		Name:               "getDNSRecord",
		Action:             "domain/%d/dns-records/%d",
		Method:             http.MethodGet,
		Result:             DNSRecord{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	},
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
	"getRedirect": {
		Name:               "getRedirect",
		Action:             "domain/%d/redirects/%s/%d",
		Method:             http.MethodGet,
		Result:             Redirect{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	},
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
	"getIPFilter": {
		Name:               "getIPFilter",
		Action:             "domain/%d/ip-filters/%s/%d",
		Method:             http.MethodGet,
		Result:             IPFilter{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	},
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
	"getSSLCertificate": {
		Name:               "getSSLCertificate",
		Action:             "domain/%d/ssl/certificates/%d",
		Method:             http.MethodGet,
		Result:             SSLCertificate{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	},
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

	//Error Page related API calls
	"listErrorPages": {
		Name:   "listErrorPages",
		Action: "domain/%d/errorpages",
		Method: http.MethodGet,
		Result: []ErrorPage{},
	},
	"getErrorPage": {
		Name:               "getErrorPage",
		Action:             "domain/%d/errorpages/%d",
		Method:             http.MethodGet,
		Result:             ErrorPage{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	},
	"createErrorPage": {
		Name:               "createErrorPage",
		Action:             "domain/%d/errorpages",
		Method:             http.MethodPost,
		Result:             ErrorPage{},
		ResponseDecodeFunc: decodeErrorPageResponse,
	},
	"updateErrorPage": {
		Name:               "updateErrorPage",
		Action:             "domain/%d/errorpages",
		Method:             http.MethodPost,
		Result:             ErrorPage{},
		ResponseDecodeFunc: decodeErrorPageResponse,
	},
	"deleteErrorPage": {
		Name:   "deleteErrorPage",
		Action: "domain/%d/errorpages",
		Method: http.MethodDelete,
		Result: ErrorPage{},
	},

	//Maintenance related API calls
	"listMaintenances": {
		Name:   "listMaintenances",
		Action: "domain/%d/%s/maintenances",
		Method: http.MethodGet,
		Result: []Maintenance{},
	},
	"createMaintenance": {
		Name:   "createMaintenance",
		Action: "domain/%d/%s/maintenances",
		Method: http.MethodPost,
		Result: Maintenance{},
	},
	"updateMaintenance": {
		Name:   "updateMaintenance",
		Action: "domain/%d/%s/maintenances/%d",
		Method: http.MethodPut,
		Result: Maintenance{},
	},
	"deleteMaintenance": {
		Name:   "deleteMaintenance",
		Action: "domain/%d/%s/maintenances/%d",
		Method: http.MethodDelete,
		Result: Maintenance{},
	},

	//MaintenanceTemplate related API calls
	"listMaintenanceTemplates": {
		Name:   "listMaintenanceTemplates",
		Action: "domain/%d/maintenance-templates",
		Method: http.MethodGet,
		Result: []MaintenanceTemplate{},
	},
	"createMaintenanceTemplate": {
		Name:   "createMaintenanceTemplate",
		Action: "domain/%d/maintenance-templates",
		Method: http.MethodPost,
		Result: MaintenanceTemplate{},
	},
	"updateMaintenanceTemplate": {
		Name:   "updateMaintenanceTemplate",
		Action: "domain/%d/maintenance-templates/%d",
		Method: http.MethodPut,
		Result: MaintenanceTemplate{},
	},
	"deleteMaintenanceTemplate": {
		Name:   "deleteMaintenanceTemplate",
		Action: "domain/%d/maintenance-templates/%d",
		Method: http.MethodDelete,
		Result: MaintenanceTemplate{},
	},

	"clearCache": {
		Name:   "clearCache",
		Action: "domain/%d/cache-clear",
		Method: http.MethodPut,
		Result: CacheClear{},
	},

	//Tag related API calls
	"getTag": {
		Name:               "getTag",
		Action:             "tags/%d",
		Method:             http.MethodGet,
		Result:             Tag{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	},
	"listTags": {
		Name:   "listTags",
		Action: "tags",
		Method: http.MethodGet,
		Result: []Tag{},
	},
	"createTags": {
		Name:   "createTags",
		Action: "tags",
		Method: http.MethodPost,
		Result: Tag{},
	},
	"updateTags": {
		Name:   "updateTags",
		Action: "tags",
		Method: http.MethodPut,
		Result: Tag{},
	},
	"deleteTags": {
		Name:   "deleteTags",
		Action: "tags",
		Method: http.MethodDelete,
		Result: Tag{},
	},
}
