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
		"updateSettings": {
			Name:   "updateSettings",
			Action: "domain/%d/%s/settings",
			Method: http.MethodPost,
			Result: Settings{},
		},
	}
}

// Settings ...
type Settings struct {
	AccessLog                   bool     `json:"access_log"`
	AntibotPostFlood            bool     `json:"antibot_post_flood"`
	AntibotPostFloodThreshold   int      `json:"antibot_post_flood_threshold,omitempty"`
	AntibotProofOfWork          bool     `json:"antibot_proof_of_work"`
	AntibotProofOfWorkThreshold int      `json:"antibot_proof_of_work_threshold,omitempty"`
	BalancingMethod             string   `json:"balancing_method,omitempty"`
	BlockNotWhitelisted         bool     `json:"block_not_whitelisted"`
	BlockTorNetwork             bool     `json:"block_tor_network"`
	CacheEnabled                bool     `json:"cache_enabled"`
	CacheRevalidate             bool     `json:"cache_revalidate"`
	CDN                         bool     `json:"cdn"`
	ClientMaxBodySize           int      `json:"client_max_body_size,omitempty"`
	DiffieHellmanExchange       int      `json:"diffie_hellman_exchange,omitempty"`
	EnableOriginSNI             bool     `json:"enable_origin_sni"`
	ForwardedForReplacement     string   `json:"forwarded_for_replacement,omitempty"`
	HSTS                        bool     `json:"hsts"`
	HSTSIncludeSubdomains       bool     `json:"hsts_include_subdomains"`
	HSTSMaxAge                  int      `json:"hsts_max_age,omitempty"`
	HSTSPreload                 bool     `json:"hsts_preload"`
	HTTPOriginPort              int      `json:"http_origin_port,omitempty"`
	IgnoreNoCache               bool     `json:"ignore_nocache"`
	ImageOptimization           bool     `json:"image_optimization"`
	IPv6Active                  bool     `json:"ipv6_active"`
	LimitAllowedHTTPMethod      []string `json:"limit_allowed_http_method,omitempty"`
	LimitTLSVersion             []string `json:"limit_tls_version,omitempty"`
	LogFormat                   string   `json:"log_format,omitempty"`
	MonitoringAlertThreshold    int      `json:"monitoring_alert_threshold,omitempty"`
	MonitoringContactEMail      string   `json:"monitoring_contact_email,omitempty"`
	MonitoringSendAlert         bool     `json:"monitoring_send_alert"`
	MyraSSLHeader               bool     `json:"myra_ssl_header"`
	NextUpstream                []string `json:"next_upstream,omitempty"`
	OnlyHTTPS                   bool     `json:"only_https"`
	OriginConnectionHeader      string   `json:"origin_connection_header,omitempty"`
	ProxyCacheBypass            string   `json:"proxy_cache_bypass,omitempty"`
	ProxyCacheStale             []string `json:"proxy_cache_stale,omitempty"`
	ProxyConnectTimeout         int      `json:"proxy_connect_timeout,omitempty"`
	ProxyReadTimeout            int      `json:"proxy_read_timeout,omitempty"`
	RequestLimitBlock           string   `json:"request_limit_block,omitempty"`
	RequestLimitLevel           int      `json:"request_limit_level,omitempty"`
	RequestLimitReport          bool     `json:"request_limit_report"`
	RequestLimitReportEMail     string   `json:"request_limit_report_email,omitempty"`
	Rewrite                     bool     `json:"rewrite"`
	SourceProtocol              string   `json:"source_protocol,omitempty"`
	Spdy                        bool     `json:"spdy"`
	SSLOriginPort               int      `json:"ssl_origin_port,omitempty"`
	WAFEnable                   bool     `json:"waf_enable"`
	WAFLevelsEnable             []string `json:"waf_levels_enable,omitempty"`
	WAFPolicy                   string   `json:"waf_policy,omitempty"`
	ProxyHostHeader             *string  `json:"host_header"`
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

// UpdateSettings updates the passed settings using the MYRA API
func (api *API) UpdateSettings(settings *Settings, domainId int, subDomainName string) (*Settings, error) {
	if _, ok := methods["updateSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createSettings")
	}

	definition := methods["updateSettings"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, settings)
	if err != nil {
		return nil, err
	}
	return result.(*Settings), nil
}

// decodeSettingsResponse - custom decode function for settings response. Used in the ListSettings action.
func decodeSettingsResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var res Settings
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
