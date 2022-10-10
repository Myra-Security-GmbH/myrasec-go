package myrasec

import "testing"

func TestInitializeMethods(t *testing.T) {
	if len(methods) <= 0 {
		t.Error("Expected to have some APIMethods defined in the methods variable")
	}

	for k := range getCacheClearMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getCacheSettingMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getDNSRecordMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getDomainMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getErrorPageMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getIPFilterMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getIPRangeMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getMaintenanceMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getMaintenanceTemplateMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getRateLimitMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getRedirectMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getSettingsMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getSSLMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getTagCacheSettingMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getTagRateLimitMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getTagSettingsMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getTagWAFRuleMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getTagMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getVHostMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getWAFMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getBucketMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}

	for k := range getFileMethods() {
		if _, ok := methods[k]; !ok {
			t.Errorf("Expected to find [%s] in the methods variable", k)
		}
	}
}
