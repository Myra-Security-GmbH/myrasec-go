package myrasec

import (
	"testing"
)

func TestListSettings(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/www.example.com/settings?flat=",
			`{
				"access_log": true,
				"antibot_post_flood": true,
				"antibot_post_flood_threshold": 10,
				"antibot_proof_of_work": true,
				"antibot_proof_of_work_threshold": 10,
				"balancing_method": "round_robin",
				"block_not_whitelisted": false,
				"block_tor_network": false,
				"cache_enabled": true,
				"cache_revalidate": false,
				"client_max_body_size": 10,
				"diffie_hellman_exchange": 2048,
				"disable_forwarded_for": true,
				"enable_origin_sni": false,
				"enforce_cache_ttl": true,
				"forwarded_for_replacement": "",
				"hsts": false,
				"hsts_include_subdomains": false,
				"hsts_max_age": 10,
				"hsts_preload": false,
				"http_origin_port": 80,
				"ignore_nocache": false,
				"image_optimization": false,
				"ipv6_active": true,
				"limit_allowed_http_method": ["GET", "HEAD"],
				"limit_tls_version": ["TLSv1.2", "TLSv1.3"],
				"log_format": "",
				"monitoring_alert_threshold": 10,
				"monitoring_contact_email": "test@example.com",
				"monitoring_send_alert": true,
				"myra_ssl_header": true,
				"next_upstream": [],
				"only_https": true,
				"origin_connection_header": "",
				"proxy_cache_bypass": "",
				"proxy_cache_stale": [],
				"proxy_connect_timeout": 10,
				"proxy_read_timeout": 10,
				"request_limit_block": "",
				"request_limit_level": 10,
				"request_limit_report": false,
				"request_limit_report_email": "",
				"rewrite": false,
				"source_protocol": "https",
				"spdy": true,
				"ssl_origin_port": 443,
				"ssl_client_certificate": [],
				"waf_enable": true,
				"waf_levels_enable": [],
				"waf_policy": "allow"
			}`,
			methods["listSettings"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	settings, err := api.ListSettings(1, "www.example.com", nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if settings.AccessLog == nil || *settings.AccessLog != true {
		t.Errorf("Expected to get Setting with AccessLog [%t] but got [%s]", true, boolPtrStr(settings.AccessLog))
	}

	if settings.AntibotPostFlood == nil || *settings.AntibotPostFlood != true {
		t.Errorf("Expected to get Setting with AntibotPostFlood [%t] but got [%s]", true, boolPtrStr(settings.AntibotPostFlood))
	}

	if settings.AntibotPostFloodThreshold != 10 {
		t.Errorf("Expected to get Setting with AntibotPostFloodThreshold [%d] but got [%d]", 10, settings.AntibotPostFloodThreshold)
	}

	if settings.AntibotProofOfWork == nil || *settings.AntibotProofOfWork != true {
		t.Errorf("Expected to get Setting with AntibotProofOfWork [%t] but got [%s]", true, boolPtrStr(settings.AntibotProofOfWork))
	}

	if settings.AntibotProofOfWorkThreshold != 10 {
		t.Errorf("Expected to get Setting with AntibotProofOfWorkThresh [%d] but got [%d]", 10, settings.AntibotProofOfWorkThreshold)
	}

	if settings.BalancingMethod != "round_robin" {
		t.Errorf("Expected to get Setting with BalancingMethod [%s] but got [%s]", "round_robin", settings.BalancingMethod)
	}

	if settings.BlockNotWhitelisted == nil || *settings.BlockNotWhitelisted != false {
		t.Errorf("Expected to get Setting with BlockNotWhitelisted [%t] but got [%s]", false, boolPtrStr(settings.BlockNotWhitelisted))
	}

	if settings.BlockTorNetwork == nil || *settings.BlockTorNetwork != false {
		t.Errorf("Expected to get Setting with BlockTorNetwork [%t] but got [%s]", false, boolPtrStr(settings.BlockTorNetwork))
	}

	if settings.CacheEnabled == nil || *settings.CacheEnabled != true {
		t.Errorf("Expected to get Setting CacheEnabled [%t] but got [%s]", true, boolPtrStr(settings.CacheEnabled))
	}

	if settings.CacheRevalidate == nil || *settings.CacheRevalidate != false {
		t.Errorf("Expected to get Setting with CacheRevalidate [%t] but got [%s]", false, boolPtrStr(settings.CacheRevalidate))
	}

	if settings.ClientMaxBodySize != 10 {
		t.Errorf("Expected to get Setting with ClientMaxBodySize [%d] but got [%d]", 10, settings.ClientMaxBodySize)
	}

	if settings.DiffieHellmanExchange != 2048 {
		t.Errorf("Expected to get Setting with DiffieHellmanExchange [%d] but got [%d]", 2048, settings.DiffieHellmanExchange)
	}

	if settings.DisableForwardFor == nil || *settings.DisableForwardFor != true {
		t.Errorf("Expected to get Setting with DisableForwardFor [%t] but got [%s]", true, boolPtrStr(settings.DisableForwardFor))
	}

	if settings.EnableOriginSNI == nil || *settings.EnableOriginSNI != false {
		t.Errorf("Expected to get Setting with EnableOriginSNI [%t] but got [%s]", false, boolPtrStr(settings.EnableOriginSNI))
	}

	if settings.EnforceCacheTTL == nil || *settings.EnforceCacheTTL != true {
		t.Errorf("Expected to get Setting with EnforceCacheTTL [%t] but got [%s]", true, boolPtrStr(settings.EnforceCacheTTL))
	}

	if settings.ForwardedForReplacement != "" {
		t.Errorf("Expected to get Setting with ForwardedForReplacement [%s] but got [%s]", "", settings.ForwardedForReplacement)
	}

	if settings.HSTS == nil || *settings.HSTS != false {
		t.Errorf("Expected to get Setting with HSTS [%t] but got [%s]", false, boolPtrStr(settings.HSTS))
	}

	if settings.HSTSIncludeSubdomains == nil || *settings.HSTSIncludeSubdomains != false {
		t.Errorf("Expected to get Setting with HSTSIncludeSubdomains [%t] but got [%s]", false, boolPtrStr(settings.HSTSIncludeSubdomains))
	}

	if settings.HSTSMaxAge != 10 {
		t.Errorf("Expected to get Setting with HSTSMaxAge [%d] but got [%d]", 10, settings.HSTSMaxAge)
	}

	if settings.HSTSPreload == nil || *settings.HSTSPreload != false {
		t.Errorf("Expected to get Setting with HSTSPreload [%t] but got [%s]", false, boolPtrStr(settings.HSTSPreload))
	}

	if settings.HTTPOriginPort != 80 {
		t.Errorf("Expected to get Setting with HTTPOriginPort [%d] but got [%d]", 80, settings.HTTPOriginPort)
	}

	if settings.IgnoreNoCache == nil || *settings.IgnoreNoCache != false {
		t.Errorf("Expected to get Setting with IgnoreNoCache [%t] but got [%s]", false, boolPtrStr(settings.IgnoreNoCache))
	}

	if settings.ImageOptimization == nil || *settings.ImageOptimization != false {
		t.Errorf("Expected to get Setting with ImageOptimization [%t] but got [%s]", false, boolPtrStr(settings.ImageOptimization))
	}

	if settings.IPv6Active == nil || *settings.IPv6Active != true {
		t.Errorf("Expected to get Setting with IPv6Active [%t] but got [%s]", true, boolPtrStr(settings.IPv6Active))
	}

	if settings.LogFormat != "" {
		t.Errorf("Expected to get Setting with LogFormat [%s] but got [%s]", "", settings.LogFormat)
	}

	if settings.MonitoringAlertThreshold != 10 {
		t.Errorf("Expected to get Setting with MonitoringAlertThreshold [%d] but got [%d]", 10, settings.MonitoringAlertThreshold)
	}

	if settings.MonitoringContactEMail != "test@example.com" {
		t.Errorf("Expected to get Setting with MonitoringContactEMail [%s] but got [%s]", "test@example.com", settings.MonitoringContactEMail)
	}

	if settings.MonitoringSendAlert == nil || *settings.MonitoringSendAlert != true {
		t.Errorf("Expected to get Setting with MonitoringSendAlert [%t] but got [%s]", true, boolPtrStr(settings.MonitoringSendAlert))
	}

	if settings.MyraSSLHeader == nil || *settings.MyraSSLHeader != true {
		t.Errorf("Expected to get Setting with MyraSSLHeader [%t] but got [%s]", true, boolPtrStr(settings.MyraSSLHeader))
	}

	if settings.OnlyHTTPS == nil || *settings.OnlyHTTPS != true {
		t.Errorf("Expected to get Setting with OnlyHTTPS [%t] but got [%s]", true, boolPtrStr(settings.OnlyHTTPS))
	}

	if settings.OriginConnectionHeader != "" {
		t.Errorf("Expected to get Setting with OriginConnectionHeader [%s] but got [%s]", "", settings.OriginConnectionHeader)
	}

	if settings.ProxyCacheBypass != "" {
		t.Errorf("Expected to get Setting with ProxyCacheBypass [%s] but got [%s]", "", settings.ProxyCacheBypass)
	}

	if settings.ProxyConnectTimeout != 10 {
		t.Errorf("Expected to get Setting with ProxyConnectTimeout [%d] but got [%d]", 10, settings.ProxyConnectTimeout)
	}

	if settings.ProxyReadTimeout != 10 {
		t.Errorf("Expected to get Setting with ProxyReadTimeout [%d] but got [%d]", 10, settings.ProxyReadTimeout)
	}

	if settings.RequestLimitBlock != "" {
		t.Errorf("Expected to get Setting with RequestLimitBlock [%s] but got [%s]", "", settings.RequestLimitBlock)
	}

	if settings.RequestLimitLevel != 10 {
		t.Errorf("Expected to get Setting with RequestLimitLevel [%d] but got [%d]", 10, settings.RequestLimitLevel)
	}

	if settings.RequestLimitReport == nil || *settings.RequestLimitReport != false {
		t.Errorf("Expected to get Setting with RequestLimitReport [%t] but got [%s]", false, boolPtrStr(settings.RequestLimitReport))
	}

	if settings.RequestLimitReportEMail != "" {
		t.Errorf("Expected to get Setting with RequestLimitReportEMail [%s] but got [%s]", "", settings.RequestLimitReportEMail)
	}

	if settings.Rewrite == nil || *settings.Rewrite != false {
		t.Errorf("Expected to get Setting with Rewrite [%t] but got [%s]", false, boolPtrStr(settings.Rewrite))
	}

	if settings.SourceProtocol != "https" {
		t.Errorf("Expected to get Setting with SourceProtocol [%s] but got [%s]", "https", settings.SourceProtocol)
	}

	if settings.Spdy == nil || *settings.Spdy != true {
		t.Errorf("Expected to get Setting with Spdy [%t] but got [%s]", true, boolPtrStr(settings.Spdy))
	}

	if settings.SSLOriginPort != 443 {
		t.Errorf("Expected to get Setting with SSLOriginPort [%d] but got [%d]", 443, settings.SSLOriginPort)
	}

	if settings.WAFEnable == nil || *settings.WAFEnable != true {
		t.Errorf("Expected to get Setting with WAFEnable [%t] but got [%s]", true, boolPtrStr(settings.WAFEnable))
	}

	if settings.WAFPolicy != "allow" {
		t.Errorf("Expected to get Setting with WAFPolicy [%s] but got [%s]", "allow", settings.WAFPolicy)
	}

}
