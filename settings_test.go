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
				"cdn": true,
				"client_max_body_size": 10,
				"diffie_hellman_exchange": 2048,
				"enable_origin_sni": false,
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

	if settings.AccessLog != true {
		t.Errorf("Expected to get Setting with AccessLog [%t] but got [%t]", true, settings.AccessLog)
	}

	if settings.AntibotPostFlood != true {
		t.Errorf("Expected to get Setting with AntibotPostFlood [%t] but got [%t]", true, settings.AntibotPostFlood)
	}

	if settings.AntibotPostFloodThreshold != 10 {
		t.Errorf("Expected to get Setting with AntibotPostFloodThreshold [%d] but got [%d]", 10, settings.AntibotPostFloodThreshold)
	}

	if settings.AntibotProofOfWork != true {
		t.Errorf("Expected to get Setting with AntibotProofOfWork [%t] but got [%t]", true, settings.AntibotProofOfWork)
	}

	if settings.AntibotProofOfWorkThreshold != 10 {
		t.Errorf("Expected to get Setting with AntibotProofOfWorkThresh [%d] but got [%d]", 10, settings.AntibotProofOfWorkThreshold)
	}

	if settings.BalancingMethod != "round_robin" {
		t.Errorf("Expected to get Setting with BalancingMethod [%s] but got [%s]", "round_robin", settings.BalancingMethod)
	}

	if settings.BlockNotWhitelisted != false {
		t.Errorf("Expected to get Setting with BlockNotWhitelisted [%t] but got [%t]", false, settings.BlockNotWhitelisted)
	}

	if settings.BlockTorNetwork != false {
		t.Errorf("Expected to get Setting with BlockTorNetwork [%t] but got [%t]", false, settings.BlockTorNetwork)
	}

	if settings.CacheEnabled != true {
		t.Errorf("Expected to get Setting CacheEnabled [%t] but got [%t]with ", true, settings.CacheEnabled)
	}

	if settings.CacheRevalidate != false {
		t.Errorf("Expected to get Setting with CacheRevalidate [%t] but got [%t]", false, settings.CacheRevalidate)
	}

	if settings.CDN != true {
		t.Errorf("Expected to get Setting with CDN [%t] but got [%t]", true, settings.CDN)
	}

	if settings.ClientMaxBodySize != 10 {
		t.Errorf("Expected to get Setting with ClientMaxBodySize [%d] but got [%d]", 10, settings.ClientMaxBodySize)
	}

	if settings.DiffieHellmanExchange != 2048 {
		t.Errorf("Expected to get Setting with DiffieHellmanExchange [%d] but got [%d]", 2048, settings.DiffieHellmanExchange)
	}

	if settings.EnableOriginSNI != false {
		t.Errorf("Expected to get Setting with EnableOriginSNI [%t] but got [%t]", false, settings.EnableOriginSNI)
	}

	if settings.ForwardedForReplacement != "" {
		t.Errorf("Expected to get Setting with ForwardedForReplacement [%s] but got [%s]", "", settings.ForwardedForReplacement)
	}

	if settings.HSTS != false {
		t.Errorf("Expected to get Setting with HSTS [%t] but got [%t]", false, settings.HSTS)
	}

	if settings.HSTSIncludeSubdomains != false {
		t.Errorf("Expected to get Setting with HSTSIncludeSubdomains [%t] but got [%t]", false, settings.HSTSIncludeSubdomains)
	}

	if settings.HSTSMaxAge != 10 {
		t.Errorf("Expected to get Setting with HSTSMaxAge [%d] but got [%d]", 10, settings.HSTSMaxAge)
	}

	if settings.HSTSPreload != false {
		t.Errorf("Expected to get Setting with HSTSPreload [%t] but got [%t]", false, settings.HSTSPreload)
	}

	if settings.HTTPOriginPort != 80 {
		t.Errorf("Expected to get Setting with HTTPOriginPort [%d] but got [%d]", 80, settings.HTTPOriginPort)
	}

	if settings.IgnoreNoCache != false {
		t.Errorf("Expected to get Setting with IgnoreNoCache [%t] but got [%t]", false, settings.IgnoreNoCache)
	}

	if settings.ImageOptimization != false {
		t.Errorf("Expected to get Setting with ImageOptimization [%t] but got [%t]", false, settings.ImageOptimization)
	}

	if settings.IPv6Active != true {
		t.Errorf("Expected to get Setting with IPv6Active [%t] but got [%t]", true, settings.IPv6Active)
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

	if settings.MonitoringSendAlert != true {
		t.Errorf("Expected to get Setting with MonitoringSendAlert [%t] but got [%t]", true, settings.MonitoringSendAlert)
	}

	if settings.MyraSSLHeader != true {
		t.Errorf("Expected to get Setting with MyraSSLHeader [%t] but got [%t]", true, settings.MyraSSLHeader)
	}

	if settings.OnlyHTTPS != true {
		t.Errorf("Expected to get Setting with OnlyHTTPS [%t] but got [%t]", true, settings.OnlyHTTPS)
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

	if settings.RequestLimitReport != false {
		t.Errorf("Expected to get Setting with RequestLimitReport [%t] but got [%t]", false, settings.RequestLimitReport)
	}

	if settings.RequestLimitReportEMail != "" {
		t.Errorf("Expected to get Setting with RequestLimitReportEMail [%s] but got [%s]", "", settings.RequestLimitReportEMail)
	}

	if settings.Rewrite != false {
		t.Errorf("Expected to get Setting with Rewrite [%t] but got [%t]", false, settings.Rewrite)
	}

	if settings.SourceProtocol != "https" {
		t.Errorf("Expected to get Setting with SourceProtocol [%s] but got [%s]", "https", settings.SourceProtocol)
	}

	if settings.Spdy != true {
		t.Errorf("Expected to get Setting with Spdy [%t] but got [%t]", true, settings.Spdy)
	}

	if settings.SSLOriginPort != 443 {
		t.Errorf("Expected to get Setting with SSLOriginPort [%d] but got [%d]", 443, settings.SSLOriginPort)
	}

	if settings.WAFEnable != true {
		t.Errorf("Expected to get Setting with WAFEnable [%t] but got [%t]", true, settings.WAFEnable)
	}

	if settings.WAFPolicy != "allow" {
		t.Errorf("Expected to get Setting with WAFPolicy [%s] but got [%s]", "allow", settings.WAFPolicy)
	}

}
