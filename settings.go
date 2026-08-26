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

// Settings represents the comprehensive configuration for a domain or subdomain.
// It controls caching, WAF, security features and upstream balancing behavior.
type Settings struct {
	// AccessLog controls the generation of separate access logs.
	// If enabled, logs from each edge node are aggregated and available via SFTP.
	AccessLog *bool `json:"access_log,omitempty" jsonschema:"Enables separate access logging. If true, logs are saved and downloadable via SFTP from 'custlogs.myracloud.com'."`

	// AntibotPostFlood enables the detection of POST flood attacks.
	// Uses a JavaScript-based puzzle to verify the client.
	AntibotPostFlood *bool `json:"antibot_post_flood,omitempty" jsonschema:"Enables detection of POST floods using a JavaScript Proof-of-Work puzzle."`

	// AntibotPostFloodThreshold sets the trigger frequency for the POST flood puzzle.
	// Higher values mean the puzzle is presented less frequently.
	AntibotPostFloodThreshold int `json:"antibot_post_flood_threshold,omitempty" jsonschema:"Determines the puzzle frequency. Higher values decrease the likelihood/frequency of the puzzle challenge."`

	// AntibotProofOfWork enables general bot detection via JS puzzles.
	AntibotProofOfWork *bool `json:"antibot_proof_of_work,omitempty" jsonschema:"Enables validation of legitimate clients using a JavaScript Proof-of-Work puzzle."`

	// AntibotProofOfWorkThreshold sets the trigger frequency for the general PoW puzzle.
	AntibotProofOfWorkThreshold int `json:"antibot_proof_of_work_threshold,omitempty" jsonschema:"Determines the puzzle frequency for general bot detection. Higher values decrease the likelihood/frequency of the challenge."`

	// BalancingMethod defines the strategy for distributing requests to upstream servers.
	// Valid values: 'round-robin', 'ip_hash', 'least_conn'.
	BalancingMethod string `json:"balancing_method,omitempty" jsonschema:"The load balancing strategy. Valid values: 'round-robin' (default), 'ip_hash' (sticky IP), 'least_conn' (lowest active connections)."`

	// BlockNotWhitelisted blocks all IPs not explicitly whitelisted in IP filters.
	BlockNotWhitelisted *bool `json:"block_not_whitelisted,omitempty" jsonschema:"Security toggle: If true, blocks ALL IPs that are not explicitly whitelisted in the IP Filter settings."`

	// BlockTorNetwork blocks traffic originating from known Tor exit nodes.
	BlockTorNetwork *bool `json:"block_tor_network,omitempty" jsonschema:"Security toggle: If true, blocks all traffic originating from the Tor anonymity network."`

	// CacheEnabled toggles the caching engine.
	// Requires defined Cache Settings objects to function effectively.
	CacheEnabled *bool `json:"cache_enabled,omitempty" jsonschema:"Master switch for caching. If true, you must also define specific Cache Setting objects for caching to occur."`

	// CacheRevalidate forces revalidation of expired cache items.
	// Uses 'If-Modified-Since' and 'If-None-Match' headers.
	CacheRevalidate *bool `json:"cache_revalidate,omitempty" jsonschema:"If true, expired cache items are revalidated with the origin using conditional HTTP headers (If-Modified-Since/If-None-Match)."`

	// CDN is a deprecated setting.
	// It has no effect and should not be used.
	// Kept as plain bool to avoid breaking consumers on a field that does nothing.
	CDN bool `json:"cdn,omitempty" jsonschema:"Deprecated setting. Has no effect."`

	// ClientMaxBodySize sets the maximum allowed size of the request body.
	// Matches the 'Content-Length' header.
	ClientMaxBodySize int `json:"client_max_body_size,omitempty" jsonschema:"Maximum allowed size of the client request body (in bytes). Requests exceeding this limit are rejected."`

	// CookieName is the name of the cookie used for stickiness.
	// Only used when BalancingMethod is set to 'cookie_based' (custom).
	CookieName string `json:"cookie_name,omitempty" jsonschema:"The specific cookie name to use for session stickiness. Only relevant if balancing_method is set to 'cookie_based'."`

	// DiffieHellmanExchange defines the bit size of DH parameters.
	// Note: Java 6/7 clients do not support >1024 bits.
	DiffieHellmanExchange int `json:"diffie_hellman_exchange,omitempty" jsonschema:"The size of Diffie-Hellman parameters in bits. Standard is 2048. Use 1024 only if legacy Java 6/7 support is required."`

	// DisableForwardFor disables the automatic injection/replacement of the Forwarded-For header.
	DisableForwardFor *bool `json:"disable_forwarded_for,omitempty" jsonschema:"If true, disables the automatic replacement/injection of the 'X-Forwarded-For' header."`

	// EnableOriginSNI allows SNI (Server Name Indication) when connecting to the origin.
	EnableOriginSNI *bool `json:"enable_origin_sni,omitempty" jsonschema:"Enables SNI (Server Name Indication) for upstream SSL handshakes. Required if the origin serves multiple certificates on one IP."`

	// EnforceCacheTTL overrides origin cache headers with Myra settings.
	EnforceCacheTTL *bool `json:"enforce_cache_ttl,omitempty" jsonschema:"If true, ignores the origin's Cache-Control headers and enforces the TTL configured in Myra settings."`

	// ForwardedForReplacement allows setting a custom name for the client IP header.
	ForwardedForReplacement string `json:"forwarded_for_replacement,omitempty" jsonschema:"Allows defining a custom header name to transport the original client IP (replacing standard X-Forwarded-For)."`

	// HSTS enables Strict-Transport-Security.
	// Forces browsers to interact with the domain only via HTTPS.
	HSTS *bool `json:"hsts,omitempty" jsonschema:"Enables HTTP Strict Transport Security (HSTS). Forces browsers to use HTTPS only."`

	// HSTSIncludeSubdomains extends HSTS protection to all subdomains.
	HSTSIncludeSubdomains *bool `json:"hsts_include_subdomains,omitempty" jsonschema:"If true, the HSTS policy applies to all subdomains as well."`

	// HSTSMaxAge defines the duration (in seconds) the HSTS header is valid.
	HSTSMaxAge int `json:"hsts_max_age,omitempty" jsonschema:"The duration (in seconds) for which the browser should remember to force HTTPS."`

	// HSTSPreload allows the domain to be submitted to the global HSTS preload list.
	HSTSPreload *bool `json:"hsts_preload,omitempty" jsonschema:"If true, allows the domain to be included in the browser hardcoded HSTS preload list (requires valid HTTPS setup)."`

	// HTTPOriginPort sets the port for plain HTTP upstream connections.
	HTTPOriginPort int `json:"http_origin_port,omitempty" jsonschema:"The TCP port used to connect to the origin server via plain HTTP (usually 80)."`

	// IgnoreNoCache forces caching even if the origin sends 'no-cache' headers.
	IgnoreNoCache *bool `json:"ignore_nocache,omitempty" jsonschema:"If true, the system ignores 'Cache-Control: private/no-store/no-cache' headers from the origin and caches content anyway."`

	// ImageOptimization enables lossless compression for JPEG and PNGs.
	ImageOptimization *bool `json:"image_optimization,omitempty" jsonschema:"Enables automatic lossless compression/optimization of JPEG and PNG images."`

	// IPLock prevents accidental IP address changes via the API/GUI.
	// Only available at the general domain level.
	IPLock *bool `json:"ip_lock,omitempty" jsonschema:"Protective lock. If true, prevents changes to the domain's IP configuration. Only available on domain level."`

	// IPv6Active enables IPv6 connectivity for the domain.
	IPv6Active *bool `json:"ipv6_active,omitempty" jsonschema:"Enables IPv6 access for clients. IPv6 traffic is translated to IPv4 if the origin is IPv4-only."`

	// LimitAllowedHTTPMethod restricts the HTTP methods accepted by the edge.
	// E.g., ["GET", "POST"].
	LimitAllowedHTTPMethod []string `json:"limit_allowed_http_method,omitempty" jsonschema:"List of allowed HTTP methods (e.g., ['GET', 'POST']). All other methods will be blocked (405 Method Not Allowed)."`

	// LimitTLSVersion restricts the allowed TLS protocol versions.
	// E.g., ["TLSv1.2", "TLSv1.3"].
	LimitTLSVersion []string `json:"limit_tls_version,omitempty" jsonschema:"List of allowed TLS versions (e.g., ['TLSv1.2', 'TLSv1.3']). Older versions will be rejected."`

	// LogFormat specifies a custom log line format.
	LogFormat string `json:"log_format,omitempty" jsonschema:"Defines a custom structure for log entries."`

	// MonitoringAlertThreshold sets the error rate (errors/minute) that triggers an alert.
	MonitoringAlertThreshold int `json:"monitoring_alert_threshold,omitempty" jsonschema:"The threshold of errors per minute required to trigger a monitoring email report."`

	// MonitoringContactEMail is a space-separated list of alert recipients.
	MonitoringContactEMail string `json:"monitoring_contact_email,omitempty" jsonschema:"Space-separated list of email addresses to receive monitoring alerts."`

	// MonitoringSendAlert enables upstream error reporting.
	MonitoringSendAlert *bool `json:"monitoring_send_alert,omitempty" jsonschema:"Enables sending of email alerts when upstream errors exceed the defined threshold."`

	// MyraSSLHeader injects 'X-Myra-SSL' to indicate a secure connection to the origin.
	MyraSSLHeader *bool `json:"myra_ssl_header,omitempty" jsonschema:"If true, adds the 'X-Myra-SSL' header to requests forwarded to the origin to indicate the client used HTTPS."`

	// MyraSSLCertificate lists certificates to use for upstream authentication.
	MyraSSLCertificate []string `json:"myra_ssl_certificate,omitempty" jsonschema:"List of SSL Certificates (PEM chain) used for client authentication against the origin server."`

	// MyraSSLCertificateKey lists private keys for the upstream certificates.
	MyraSSLCertificateKey []string `json:"myra_ssl_certificate_key,omitempty" jsonschema:"List of private keys corresponding to the MyraSSLCertificate."`

	// NextUpstream defines conditions to try the next server in the pool.
	// Values: error, timeout, invalid_header, http_500, http_502, etc. 'off' disables it.
	NextUpstream []string `json:"next_upstream,omitempty" jsonschema:"Conditions under which the request is retried on the next upstream server. Examples: 'error', 'timeout', 'http_500'. Use 'off' to disable."`

	// OnlyHTTPS forces all traffic to the origin to use HTTPS.
	OnlyHTTPS *bool `json:"only_https,omitempty" jsonschema:"If true, all requests to the origin are sent via HTTPS, even if the client connected via HTTP."`

	// OriginConnectionHeader defines the 'Connection' header sent to the upstream.
	OriginConnectionHeader string `json:"origin_connection_header,omitempty" jsonschema:"Sets the value of the 'Connection' header sent to the origin (e.g., 'keep-alive' or 'close')."`

	// ProxyCacheBypass defines a cookie name that forces a cache miss.
	ProxyCacheBypass string `json:"proxy_cache_bypass,omitempty" jsonschema:"Name of a cookie. If this cookie is present (and not 0/empty), the cache is bypassed."`

	// ProxyCacheStale defines when to serve stale content on upstream errors.
	// Values: error, timeout, updating, http_500, etc.
	ProxyCacheStale []string `json:"proxy_cache_stale,omitempty" jsonschema:"Conditions under which expired (stale) cache content is delivered if the origin fails. Examples: 'error', 'timeout', 'updating'."`

	// ProxyConnectTimeout is the timeout (seconds) for connecting to the upstream.
	// Max: 60s.
	ProxyConnectTimeout int `json:"proxy_connect_timeout,omitempty" jsonschema:"Timeout in seconds for establishing a TCP connection to the origin. Maximum: 60 seconds."`

	// ProxyHostHeader sets a custom 'Host' header for upstream requests.
	// Default: current subdomain.
	ProxyHostHeader *string `json:"host_header,omitempty" jsonschema:"Overrides the 'Host' header sent to the origin. If null/empty, defaults to the request's subdomain."`

	// ProxyReadTimeout is the timeout (seconds) for reading the upstream response.
	// Applies between two successive read operations.
	ProxyReadTimeout int `json:"proxy_read_timeout,omitempty" jsonschema:"Timeout in seconds for reading the response from the origin (between two successive read operations)."`

	// RequestLimitBlock enables CAPTCHA challenges for rate-limited IPs.
	RequestLimitBlock string `json:"request_limit_block,omitempty" jsonschema:"Controls behavior when limit is reached. If set, users must solve a CAPTCHA to unblock their IP."`

	// RequestLimitLevel sets the max requests per minute per IP.
	RequestLimitLevel int `json:"request_limit_level,omitempty" jsonschema:"Rate limit threshold: Maximum requests allowed per IP per minute. Exceeding this blocks the IP."`

	// RequestLimitReport enables email reporting for rate limits.
	RequestLimitReport *bool `json:"request_limit_report,omitempty" jsonschema:"If true, sends email reports containing IPs that exceeded the request limit."`

	// RequestLimitReportEMail is a space-separated list of rate-limit report recipients.
	RequestLimitReportEMail string `json:"request_limit_report_email,omitempty" jsonschema:"Space-separated list of email addresses to receive request limit reports."`

	// Rewrite enables automated JavaScript optimization (bundling/deferring).
	Rewrite *bool `json:"rewrite,omitempty" jsonschema:"Enables automatic JavaScript optimization (bundling and deferred execution) to improve page load times."`

	// SourceProtocol defines the protocol scheme for upstream connections.
	// Values: 'same' (match client), 'http', 'https'.
	SourceProtocol string `json:"source_protocol,omitempty" jsonschema:"Protocol policy for origin connections. Valid values: 'same' (match client protocol), 'http' (force plain), 'https' (force SSL)."`

	// Spdy enables HTTP/2.
	// Note: Requires HTTPS to be active.
	Spdy *bool `json:"spdy,omitempty" jsonschema:"Enables the HTTP/2 protocol (formerly SPDY). Note: Requires active HTTPS."`

	// SSLClientVerify enables mTLS client certificate verification.
	SSLClientVerify string `json:"ssl_client_verify,omitempty" jsonschema:"Controls Mutual TLS (mTLS). Enables verification of client certificates against trusted CAs."`

	// SSLClientCertificate is a list of trusted CA certificates (PEM) for mTLS.
	SSLClientCertificate []string `json:"ssl_client_certificate,omitempty" jsonschema:"List of trusted CA certificates (PEM format) used to verify client certificates (mTLS)."`

	// SSLClientHeaderVerification is the header name containing the verification result.
	SSLClientHeaderVerification string `json:"ssl_client_header_verification,omitempty" jsonschema:"Header name that will contain the SSL verification status (e.g., 'SUCCESS', 'FAILED') forwarded to the origin."`

	// SSLClientHeaderFingerprint is the header name containing the client cert fingerprint.
	SSLClientHeaderFingerprint string `json:"ssl_client_header_fingerprint,omitempty" jsonschema:"Header name that will contain the SHA fingerprint of the client certificate."`

	// SSLOriginPort sets the port for SSL upstream connections.
	SSLOriginPort int `json:"ssl_origin_port,omitempty" jsonschema:"The TCP port used to connect to the origin server via HTTPS (usually 443)."`

	// WAFEnable toggles the Web Application Firewall.
	WAFEnable *bool `json:"waf_enable,omitempty" jsonschema:"Master switch: Enables or disables the Web Application Firewall (WAF) for this domain."`

	// WAFLevelsEnable selects the WAF rule sets to apply.
	// E.g., ["wafrules_sql", "wafrules_xss"].
	WAFLevelsEnable []string `json:"waf_levels_enable,omitempty" jsonschema:"List of active WAF rule sets (e.g., ['wafrules_sql', 'wafrules_xss'])."`

	// WAFPolicy defines the default action if a rule matches.
	// Values: 'block', 'allow', 'log'.
	WAFPolicy string `json:"waf_policy,omitempty" jsonschema:"Default action when a WAF rule is triggered. Valid values: 'block', 'allow', 'log'."`
}

// ListSettings returns a Setting struct containing the settings for the passed subdomain
func (api *API) ListSettings(domainId int, subDomainName string, params map[string]string) (*Settings, error) {
	if _, ok := api.methods["listSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSettings")
	}

	definition := api.methods["listSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*Settings)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// ListSettingsFull returns a Setting struct containing the full hierarchie of the settings
func (api *API) ListSettingsFull(domainId int, subDomainName string, params map[string]string) (any, error) {
	if _, ok := api.methods["listSettingsFull"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSettings")
	}

	definition := api.methods["listSettingsFull"]
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
	if _, ok := api.methods["updateSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSettings")
	}

	definition := api.methods["updateSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, settings)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Settings)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateSettingsPartial updates the passed settings using the MYRA API
func (api *API) UpdateSettingsPartial(settings map[string]any, domainId int, subDomainName string) (any, error) {
	if _, ok := api.methods["updateSettingsPartial"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSettingsPartial")
	}

	definition := api.methods["updateSettingsPartial"]
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
