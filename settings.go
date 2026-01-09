package myrasec

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// getSettingsMethods returns Settings related API calls
func getSettingsMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listSettings": {
			Name:               "listSettings",
			Action:             "domain/%d/%s/settings?flat",
			Method:             http.MethodGet,
			Result:             Settings{},
			ResponseDecodeFunc: decodeSettingsResponse,
		},
		"listSettingsFull": {
			Name:               "listSettingsFull",
			Action:             "domain/%d/%s/settings",
			Method:             http.MethodGet,
			Result:             map[string]any{},
			ResponseDecodeFunc: decodeSettingsResponseFull,
		},
		"updateSettings": {
			Name:   "updateSettings",
			Action: "domain/%d/%s/settings",
			Method: http.MethodPost,
			Result: Settings{},
		},
		"updateSettingsPartial": {
			Name:   "updateSettingsPartial",
			Action: "domain/%d/%s/settings",
			Method: http.MethodPost,
			Result: map[string]any{},
		},
	}
}

// Settings ...
type Settings struct {
	AccessLog                   bool     `json:"access_log,omitempty" jsonschema:"Activate separated access log. A access log from each Myra node delivering your website will be saved. You can download the access log files via sftp from custlogs.myracloud.com."`
	AntibotPostFlood            bool     `json:"antibot_post_flood,omitempty" jsonschema:"Detection of POST floods by using a JavaScript based puzzle."`
	AntibotPostFloodThreshold   int      `json:"antibot_post_flood_threshold,omitempty" jsonschema:"This parameter determines the frequency how often the puzzle has to be solved. The higher the value the less likely the puzzle needs to be solved."`
	AntibotProofOfWork          bool     `json:"antibot_proof_of_work,omitempty" jsonschema:"Detection of valid clients by using a JavaScript based puzzle."`
	AntibotProofOfWorkThreshold int      `json:"antibot_proof_of_work_threshold,omitempty" jsonschema:"This parameter determines the frequency how often the puzzle has to be solved. The higher the value the less likely the puzzle needs to be solved."`
	BalancingMethod             string   `json:"balancing_method,omitempty" jsonschema:"Specifies with which method requests are balanced between upstream servers. The default behavior is the round-robin balancing. The value ip_hash will cause Myra to forward the same client IP always to the same upstream server. The value least_conn will cause Myra to forward the request to the upstream server with least connections"`
	BlockNotWhitelisted         bool     `json:"block_not_whitelisted,omitempty" jsonschema:"Block all IPs which are not whitelisted in the IP filter settings."`
	BlockTorNetwork             bool     `json:"block_tor_network,omitempty" jsonschema:"Block traffic from the TOR network."`
	CacheEnabled                bool     `json:"cache_enabled,omitempty" jsonschema:"Turn caching on or off. If you enable the cache, you also have to define the objects to be cached in the cache settings."`
	CacheRevalidate             bool     `json:"cache_revalidate,omitempty" jsonschema:"If enabled, expired cache items will be requested with the additional HTTP header If-Modified-Since and If-None-Match."`
	CDN                         bool     `json:"cdn,omitempty" jsonschema:"This setting is deprecated and has no effect anymore."`
	ClientMaxBodySize           int      `json:"client_max_body_size,omitempty" jsonschema:"Sets the maximum allowed size of the client request body, specified in the Content-Length request header field."`
	CookieName                  string   `json:"cookie_name,omitempty" jsonschema:"Specifies the cookie name when balancing_method is cookie_based."`
	DiffieHellmanExchange       int      `json:"diffie_hellman_exchange,omitempty" jsonschema:"Defines the size of the Diffie-Hellman key exchange parameters in bits. Please, note that Java 6 and 7 do not support Diffie-Hellman parameters larger than 1024 bits. If your server expects to receive connections from java 6 clients and wants to enable PFS, it must provide a DHE parameter of 1024 bits."`
	DisableForwardFor           bool     `json:"disable_forwarded_for,omitempty" jsonschema:"Disable the forwarded for replacement."`
	EnableOriginSNI             bool     `json:"enable_origin_sni,omitempty" jsonschema:"Enable or disable origin SNI."`
	EnforceCacheTTL             bool     `json:"enforce_cache_ttl,omitempty" jsonschema:"Enforce using given cache TTL settings instead of origin cache information. This will set the Cache-Control header max-age to the given TTL."`
	ForwardedForReplacement     string   `json:"forwarded_for_replacement,omitempty" jsonschema:"Set your own X-Forwarded-For header."`
	HSTS                        bool     `json:"hsts,omitempty" jsonschema:"Enable HSTS protection for a domain. This will tell browsers to use secure https connections only when interacting with your domain."`
	HSTSIncludeSubdomains       bool     `json:"hsts_include_subdomains,omitempty" jsonschema:"This will extend the HSTS protection for all subdomains."`
	HSTSMaxAge                  int      `json:"hsts_max_age,omitempty" jsonschema:"Specified how long the HSTS header is valid before the browser has to revalidate."`
	HSTSPreload                 bool     `json:"hsts_preload,omitempty" jsonschema:"Allow the domain to be added to the HSTS preload list used by all major browsers (https://hstspreload.appspot.com/)."`
	HTTPOriginPort              int      `json:"http_origin_port,omitempty" jsonschema:"Allows to set a port for communication with origin via HTTP."`
	IgnoreNoCache               bool     `json:"ignore_nocache,omitempty" jsonschema:"If activated, no-cache headers (Cache-Control: [private|no-store|no-cache]) will be ignored."`
	ImageOptimization           bool     `json:"image_optimization,omitempty" jsonschema:"Activate lossless optimization of JPEG and PNG images (recommended setting)."`
	IPLock                      bool     `json:"ip_lock,omitempty" jsonschema:"Prevent accidental IP address changes if activated. This setting is only available on domain level (general domain settings)."`
	IPv6Active                  bool     `json:"ipv6_active,omitempty" jsonschema:"Allow connections via IPv6 to your systems. IPv4 connections will be forwarded in any case."`
	LimitAllowedHTTPMethod      []string `json:"limit_allowed_http_method,omitempty" jsonschema:"Not selected HTTP methods will be blocked."`
	LimitTLSVersion             []string `json:"limit_tls_version,omitempty" jsonschema:"Only selected TLS versions will be used."`
	LogFormat                   string   `json:"log_format,omitempty" jsonschema:"Use a different log format."`
	MonitoringAlertThreshold    int      `json:"monitoring_alert_threshold,omitempty" jsonschema:"Errors per minute that must occur until a report is sent."`
	MonitoringContactEMail      string   `json:"monitoring_contact_email,omitempty" jsonschema:"Email addresses, to which monitoring emails should be send. Multiple addresses are separated with a space."`
	MonitoringSendAlert         bool     `json:"monitoring_send_alert,omitempty" jsonschema:"Enables / disables the upstream error reporting."`
	MyraSSLHeader               bool     `json:"myra_ssl_header,omitempty" jsonschema:"Activate the X-Myra-SSL Header, which indicates if a request was received via SSL."`
	MyraSSLCertificate          []string `json:"myra_ssl_certificate,omitempty" jsonschema:"An SSL Certificate (and chain) to be used to make requests on the origin."`
	MyraSSLCertificateKey       []string `json:"myra_ssl_certificate_key,omitempty" jsonschema:"The private key for the MyraSSLCertificate."`
	NextUpstream                []string `json:"next_upstream,omitempty" jsonschema:"Specify in which case the current upstream should be marked as down. The values can be arbitrary combined, expect the value off."`
	OnlyHTTPS                   bool     `json:"only_https,omitempty" jsonschema:"If activated, Myra will forward all requests to the origin using HTTPS regardless of the used protocol of the originating request."`
	OriginConnectionHeader      string   `json:"origin_connection_header,omitempty" jsonschema:"Sets the Connection header, which is transmitted to the origin with a request."`
	ProxyCacheBypass            string   `json:"proxy_cache_bypass,omitempty" jsonschema:"Defines the name of the cookie which forces Myra to deliver the response not from cache. The values of the cookie must be not empty or equal to 0 to enable bypassing."`
	ProxyCacheStale             []string `json:"proxy_cache_stale,omitempty" jsonschema:"Determines in which cases a stale cached response can be used when an error occurs during communication with your server. The values can be arbitrary combined, expect the value off."`
	ProxyConnectTimeout         int      `json:"proxy_connect_timeout,omitempty" jsonschema:"Defines a timeout in seconds for establishing a connection with the origin server. The timeout cannot be greater than 60 seconds."`
	ProxyHostHeader             *string  `json:"host_header,omitempty" jsonschema:"Set your own Proxy Host header. The default value is the current subdomain."`
	ProxyReadTimeout            int      `json:"proxy_read_timeout,omitempty" jsonschema:"Defines a timeout in seconds for reading a response from the proxied server. The timeout is set only between two successive read operations, not for the transmission of the whole response."`
	RequestLimitBlock           string   `json:"request_limit_block,omitempty" jsonschema:"If activated, the user has to solve a CAPTCHA after exceeding the configured request limit."`
	RequestLimitLevel           int      `json:"request_limit_level,omitempty" jsonschema:"Define how many requests are allowed from an IP per minute. If this limit is reached, the IP will be blocked. If request_limit_block is enabled, the user can solve a CAPTCHA to unblock his IP address."`
	RequestLimitReport          bool     `json:"request_limit_report,omitempty" jsonschema:"If activated, an email will be send containing blocked ip addresses that exceeded the configured request limit."`
	RequestLimitReportEMail     string   `json:"request_limit_report_email,omitempty" jsonschema:"Email addresses, to which request limit emails should be send. Multiple addresses are separated with a space."`
	Rewrite                     bool     `json:"rewrite,omitempty" jsonschema:"Enable automated JavaScript optimization. All JavaScript is collected and executed at the end of the page. This significantly decreases the DOM content loaded time. If not all JavaScript files should be collected you can set the value to regex and specify the regex to use while matching filenames in the option rewrite_regex."`
	SourceProtocol              string   `json:"source_protocol,omitempty" jsonschema:"Define which protocol should be used when passing a request to your servers. The value same will ensure that the same protocol is used as in the originating request to Myra. The http and https value will force Myra to always use the specified protocol when connecting."`
	Spdy                        bool     `json:"spdy,omitempty" jsonschema:"Activate the high performance HTTP/2 protocol. Please note that you have to enable HTTPS for Myra to get HTTP/2 enabled."`
	SSLClientVerify             string   `json:"ssl_client_verify,omitempty" jsonschema:"Enables verification of client certificates."`
	SSLClientCertificate        []string `json:"ssl_client_certificate,omitempty" jsonschema:"Specifies files with trusted CA certificates in the PEM format used to verify client certificates."`
	SSLClientHeaderVerification string   `json:"ssl_client_header_verification,omitempty" jsonschema:"The name of the header, which contains the SSL verification status."`
	SSLClientHeaderFingerprint  string   `json:"ssl_client_header_fingerprint,omitempty" jsonschema:"Contains the fingerprint of the certificate, the client used to authenticate itself."`
	SSLOriginPort               int      `json:"ssl_origin_port,omitempty" jsonschema:"Allows to set a port for communication with origin via SSL."`
	WAFEnable                   bool     `json:"waf_enable,omitempty" jsonschema:"Enables or disables the Web Application Firewall."`
	WAFLevelsEnable             []string `json:"waf_levels_enable,omitempty" jsonschema:"Level of applied Web Application Firewall rules."`
	WAFPolicy                   string   `json:"waf_policy,omitempty" jsonschema:"Default policy for the Web Application Firewall in case of rule error."`
}

// ListSettings returns a Setting struct containing the settings for the passed subdomain
func (api *API) ListSettings(domainId int, subDomainName string, params map[string]string) (*Settings, error) {
	if _, ok := methods["listSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSettings")
	}

	definition := methods["listSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return result.(*Settings), nil
}

// ListSettingsFull returns a Setting struct containing the full hierarchie of the settings
func (api *API) ListSettingsFull(domainId int, subDomainName string, params map[string]string) (any, error) {
	if _, ok := methods["listSettingsFull"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSettings")
	}

	definition := methods["listSettingsFull"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateSettings updates the passed settings using the MYRA API
// Deprecated: this method uses myra-api settings in a wrong way, please use UpdateSettingsPartial instead
func (api *API) UpdateSettings(settings *Settings, domainId int, subDomainName string) (*Settings, error) {
	if _, ok := methods["updateSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSettings")
	}

	definition := methods["updateSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, settings)
	if err != nil {
		return nil, err
	}
	return result.(*Settings), nil
}

// UpdateSettingsPartial updates the passed settings using the MYRA API
func (api *API) UpdateSettingsPartial(settings map[string]any, domainId int, subDomainName string) (any, error) {
	if _, ok := methods["updateSettingsPartial"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSettingsPartial")
	}

	definition := methods["updateSettingsPartial"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, settings)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// decodeSettingsResponse - custom decode function for settings response. Used in the ListSettings action.
func decodeSettingsResponse(resp *http.Response, definition APIMethod) (any, error) {
	var res Settings
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// decodeSettingsResponseFull - custom decode function for full settings response. Used in the ListSettingsFull action.
func decodeSettingsResponseFull(resp *http.Response, definition APIMethod) (any, error) {
	var res map[string]any
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
