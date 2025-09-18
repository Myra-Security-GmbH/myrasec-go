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
	AccessLog                   bool     `json:"access_log,omitempty"`
	AntibotPostFlood            bool     `json:"antibot_post_flood,omitempty"`
	AntibotPostFloodThreshold   int      `json:"antibot_post_flood_threshold,omitempty"`
	AntibotProofOfWork          bool     `json:"antibot_proof_of_work,omitempty"`
	AntibotProofOfWorkThreshold int      `json:"antibot_proof_of_work_threshold,omitempty"`
	BalancingMethod             string   `json:"balancing_method,omitempty"`
	BlockNotWhitelisted         bool     `json:"block_not_whitelisted,omitempty"`
	BlockTorNetwork             bool     `json:"block_tor_network,omitempty"`
	CacheEnabled                bool     `json:"cache_enabled,omitempty"`
	CacheRevalidate             bool     `json:"cache_revalidate,omitempty"`
	CDN                         bool     `json:"cdn,omitempty"`
	ClientMaxBodySize           int      `json:"client_max_body_size,omitempty"`
	CookieName                  string   `json:"cookie_name,omitempty"`
	DiffieHellmanExchange       int      `json:"diffie_hellman_exchange,omitempty"`
	DisableForwardFor           bool     `json:"disable_forwarded_for,omitempty"`
	EnableOriginSNI             bool     `json:"enable_origin_sni,omitempty"`
	EnforceCacheTTL             bool     `json:"enforce_cache_ttl,omitempty"`
	ForwardedForReplacement     string   `json:"forwarded_for_replacement,omitempty"`
	HSTS                        bool     `json:"hsts,omitempty"`
	HSTSIncludeSubdomains       bool     `json:"hsts_include_subdomains,omitempty"`
	HSTSMaxAge                  int      `json:"hsts_max_age,omitempty"`
	HSTSPreload                 bool     `json:"hsts_preload,omitempty"`
	HTTPOriginPort              int      `json:"http_origin_port,omitempty"`
	IgnoreNoCache               bool     `json:"ignore_nocache,omitempty"`
	ImageOptimization           bool     `json:"image_optimization,omitempty"`
	IPLock                      bool     `json:"ip_lock,omitempty"`
	IPv6Active                  bool     `json:"ipv6_active,omitempty"`
	LimitAllowedHTTPMethod      []string `json:"limit_allowed_http_method,omitempty"`
	LimitTLSVersion             []string `json:"limit_tls_version,omitempty"`
	LogFormat                   string   `json:"log_format,omitempty"`
	MonitoringAlertThreshold    int      `json:"monitoring_alert_threshold,omitempty"`
	MonitoringContactEMail      string   `json:"monitoring_contact_email,omitempty"`
	MonitoringSendAlert         bool     `json:"monitoring_send_alert,omitempty"`
	MyraSSLHeader               bool     `json:"myra_ssl_header,omitempty"`
	MyraSSLCertificate          []string `json:"myra_ssl_certificate,omitempty"`
	MyraSSLCertificateKey       []string `json:"myra_ssl_certificate_key,omitempty"`
	NextUpstream                []string `json:"next_upstream,omitempty"`
	OnlyHTTPS                   bool     `json:"only_https,omitempty"`
	OriginConnectionHeader      string   `json:"origin_connection_header,omitempty"`
	ProxyCacheBypass            string   `json:"proxy_cache_bypass,omitempty"`
	ProxyCacheStale             []string `json:"proxy_cache_stale,omitempty"`
	ProxyConnectTimeout         int      `json:"proxy_connect_timeout,omitempty"`
	ProxyHostHeader             *string  `json:"host_header,omitempty"`
	ProxyReadTimeout            int      `json:"proxy_read_timeout,omitempty"`
	RequestLimitBlock           string   `json:"request_limit_block,omitempty"`
	RequestLimitLevel           int      `json:"request_limit_level,omitempty"`
	RequestLimitReport          bool     `json:"request_limit_report,omitempty"`
	RequestLimitReportEMail     string   `json:"request_limit_report_email,omitempty"`
	Rewrite                     bool     `json:"rewrite,omitempty"`
	SourceProtocol              string   `json:"source_protocol,omitempty"`
	Spdy                        bool     `json:"spdy,omitempty"`
	SSLClientVerify             string   `json:"ssl_client_verify,omitempty"`
	SSLClientCertificate        []string `json:"ssl_client_certificate,omitempty"`
	SSLClientHeaderVerification string   `json:"ssl_client_header_verification,omitempty"`
	SSLClientHeaderFingerprint  string   `json:"ssl_client_header_fingerprint,omitempty"`
	SSLOriginPort               int      `json:"ssl_origin_port,omitempty"`
	WAFEnable                   bool     `json:"waf_enable,omitempty"`
	WAFLevelsEnable             []string `json:"waf_levels_enable,omitempty"`
	WAFPolicy                   string   `json:"waf_policy,omitempty"`
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
