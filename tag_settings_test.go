package myrasec

import (
	"testing"
)

func TestListTagSettings(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tag/1/settings",
			`{
				"settings": {
					"access_log": true,
					"antibot_post_flood": true,
					"antibot_proof_of_work": true,
					"block_not_whitelisted": false,
					"block_tor_network": false,
					"cache_enabled": true,
					"cache_revalidate": false,
					"cdn": true,
					"enable_origin_sni": false,
					"hsts": false,
					"hsts_include_subdomains": false,
					"hsts_preload": false,
					"ignore_nocache": false,
					"image_optimization": false,
					"ipv6_active": true,
					"monitoring_send_alert": true,
					"myra_ssl_header": true,
					"only_https": true,
					"request_limit_report": false,
					"rewrite": false,
					"spdy": true,
					"waf_enable": true
				}
			}`,
			methods["listTagSettings"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	settings, err := api.ListTagSettings(1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if settings.AccessLog == nil || *settings.AccessLog != true {
		t.Errorf("Expected to get Setting with AccessLog [%t] but got [%v]", true, settings.AccessLog)
	}

	if settings.AntibotPostFlood == nil || *settings.AntibotPostFlood != true {
		t.Errorf("Expected to get Setting with AntibotPostFlood [%t] but got [%v]", true, settings.AntibotPostFlood)
	}

	if settings.AntibotProofOfWork == nil || *settings.AntibotProofOfWork != true {
		t.Errorf("Expected to get Setting with AntibotProofOfWork [%t] but got [%v]", true, settings.AntibotProofOfWork)
	}

	if settings.BlockNotWhitelisted == nil || *settings.BlockNotWhitelisted != false {
		t.Errorf("Expected to get Setting with BlockNotWhitelisted [%t] but got [%v]", false, settings.BlockNotWhitelisted)
	}

	if settings.BlockTorNetwork == nil || *settings.BlockTorNetwork != false {
		t.Errorf("Expected to get Setting with BlockTorNetwork [%t] but got [%v]", false, settings.BlockTorNetwork)
	}

	if settings.CacheEnabled == nil || *settings.CacheEnabled != true {
		t.Errorf("Expected to get Setting CacheEnabled [%t] but got [%v]", true, settings.CacheEnabled)
	}

	if settings.CacheRevalidate == nil || *settings.CacheRevalidate != false {
		t.Errorf("Expected to get Setting with CacheRevalidate [%t] but got [%v]", false, settings.CacheRevalidate)
	}

	if settings.EnableOriginSNI == nil || *settings.EnableOriginSNI != false {
		t.Errorf("Expected to get Setting with EnableOriginSNI [%t] but got [%v]", false, settings.EnableOriginSNI)
	}

	if settings.HSTS == nil || *settings.HSTS != false {
		t.Errorf("Expected to get Setting with HSTS [%t] but got [%v]", false, settings.HSTS)
	}

	if settings.HSTSIncludeSubdomains == nil || *settings.HSTSIncludeSubdomains != false {
		t.Errorf("Expected to get Setting with HSTSIncludeSubdomains [%t] but got [%v]", false, settings.HSTSIncludeSubdomains)
	}

	if settings.HSTSPreload == nil || *settings.HSTSPreload != false {
		t.Errorf("Expected to get Setting with HSTSPreload [%t] but got [%v]", false, settings.HSTSPreload)
	}

	if settings.IgnoreNoCache == nil || *settings.IgnoreNoCache != false {
		t.Errorf("Expected to get Setting with IgnoreNoCache [%t] but got [%v]", false, settings.IgnoreNoCache)
	}

	if settings.ImageOptimization == nil || *settings.ImageOptimization != false {
		t.Errorf("Expected to get Setting with ImageOptimization [%t] but got [%v]", false, settings.ImageOptimization)
	}

	if settings.IPv6Active == nil || *settings.IPv6Active != true {
		t.Errorf("Expected to get Setting with IPv6Active [%t] but got [%v]", true, settings.IPv6Active)
	}

	if settings.MonitoringSendAlert == nil || *settings.MonitoringSendAlert != true {
		t.Errorf("Expected to get Setting with MonitoringSendAlert [%t] but got [%v]", true, settings.MonitoringSendAlert)
	}

	if settings.MyraSSLHeader == nil || *settings.MyraSSLHeader != true {
		t.Errorf("Expected to get Setting with MyraSSLHeader [%t] but got [%v]", true, settings.MyraSSLHeader)
	}

	if settings.OnlyHTTPS == nil || *settings.OnlyHTTPS != true {
		t.Errorf("Expected to get Setting with OnlyHTTPS [%t] but got [%v]", true, settings.OnlyHTTPS)
	}

	if settings.RequestLimitReport == nil || *settings.RequestLimitReport != false {
		t.Errorf("Expected to get Setting with RequestLimitReport [%t] but got [%v]", false, settings.RequestLimitReport)
	}

	if settings.Rewrite == nil || *settings.Rewrite != false {
		t.Errorf("Expected to get Setting with Rewrite [%t] but got [%v]", false, settings.Rewrite)
	}

	if settings.Spdy == nil || *settings.Spdy != true {
		t.Errorf("Expected to get Setting with Spdy [%t] but got [%v]", true, settings.Spdy)
	}

	if settings.WAFEnable == nil || *settings.WAFEnable != true {
		t.Errorf("Expected to get Setting with WAFEnable [%t] but got [%v]", true, settings.WAFEnable)
	}
}
