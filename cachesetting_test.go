package myrasec

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestListCacheSettings(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/www.example.com/cache-settings",
			`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [
				{"id": 1, "path": "/index.html", "ttl": 300, "notFoundTtl": 300, "type": "exact", "enforce": false, "enabled": true, "comment": "this is a comment"}, 
				{"id": 2, "path": "/index.php", "ttl": 300, "notFoundTtl": 300, "type": "exact", "enforce": false, "enabled": false, "comment": ""}
			]}`,
			"listCacheSettings",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	settings, err := api.ListCacheSettings(1, "www.example.com", nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(settings) != 2 {
		t.Errorf("Expected to get [%d] cache settings but got [%d]", 2, len(settings))
	}

	for k, v := range settings {
		if v.ID != k+1 {
			t.Errorf("Expected to get cache setting with ID [%d] but got [%d]", k+1, v.ID)
		}

		if v.ID == 1 {
			if v.Path != "/index.html" {
				t.Errorf("Expected to get cache setting with Type [%s] but got [%s]", "/index.html", v.Path)
			}

			if v.Enabled == nil || !*v.Enabled {
				t.Errorf("Expected to get cache setting with Enabled [%t] but got [%s]", true, boolPtrStr(v.Enabled))
			}

			if v.Comment != "this is a comment" {
				t.Errorf("Expected to get cache setting with Comment [%s] but got [%s]", "this is a comment", v.Comment)
			}

		}

		if v.ID == 2 {
			if v.Path != "/index.php" {
				t.Errorf("Expected to get cache setting with Type [%s] but got [%s]", "/index.php", v.Path)
			}

			if v.Enabled == nil || *v.Enabled {
				t.Errorf("Expected to get cache setting with Enabled [%t] but got [%s]", false, boolPtrStr(v.Enabled))
			}

			if v.Comment != "" {
				t.Errorf("Expected to get cache setting with Comment [%s] but got [%s]", "", v.Comment)
			}
		}

		if v.TTL != 300 {
			t.Errorf("Expected to get cache setting with TTL [%d] but got [%d]", 300, v.TTL)
		}

		if v.NotFoundTTL != 300 {
			t.Errorf("Expected to get cache setting with NotFoundTTL [%d] but got [%d]", 300, v.NotFoundTTL)
		}

		if v.Enforce != false {
			t.Errorf("Expected to get cache setting with Enforce [%t] but got [%t]", false, v.Enforce)
		}

		if v.Type != "exact" {
			t.Errorf("Expected to get cache setting with Type [%s] but got [%s]", "exact", v.Type)
		}
	}
}

func TestCacheSettingSendsDisabledFlags(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /domain/1/www.example.com/cache-settings": {Status: http.StatusOK, Body: `{"error": false, "targetObject": [{"id": 1, "path": "/index.html", "ttl": 300, "notFoundTtl": 300, "type": "exact", "enforce": false, "enabled": false}]}`},
	})

	_, err := api.CreateCacheSettingContext(context.Background(), &CacheSetting{Path: "/index.html", TTL: 300, NotFoundTTL: 300, Type: "exact", Enabled: BoolPtr(false)}, 1, "www.example.com")
	if err != nil {
		t.Fatalf("Expected not to get an error but got %v", err)
	}

	body := string(requests.last(t).Body)
	for _, want := range []string{`"enabled":false`, `"enforce":false`} {
		if !strings.Contains(body, want) {
			t.Errorf("Expected the request body to contain %s but got %s", want, body)
		}
	}
}

func TestCacheSettingOmitsUnsetEnabled(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /domain/1/www.example.com/cache-settings": {Status: http.StatusOK, Body: `{"error": false, "targetObject": [{"id": 1, "path": "/index.html", "ttl": 300, "notFoundTtl": 300, "type": "exact", "enforce": false, "enabled": true}]}`},
	})

	// Without an explicit value the API decides, a new setting is enabled by default.
	_, err := api.CreateCacheSettingContext(context.Background(), &CacheSetting{Path: "/index.html", TTL: 300, NotFoundTTL: 300, Type: "exact"}, 1, "www.example.com")
	if err != nil {
		t.Fatalf("Expected not to get an error but got %v", err)
	}

	if body := string(requests.last(t).Body); strings.Contains(body, `"enabled"`) {
		t.Errorf("Expected the request body not to contain the enabled attribute but got %s", body)
	}
}
