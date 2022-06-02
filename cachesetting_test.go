package myrasec

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func preCacheListCacheSettings(url string, body string) *TestCache {

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(body)),
	}

	res, _ := decodeDefaultResponse(&resp, methods["listCacheSettings"])

	return &TestCache{
		Req: req,
		Res: res,
	}
}

func TestListCacheSettings(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheListCacheSettings(
			"https://apiv2.myracloud.com/domain/1/www.example.com/cache-settings",
			`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [{"id": 1, "path": "/index.html", "ttl": 300, "notFoundTtl": 300, "type": "exact", "enforce": false, "enabled": true}, {"id": 2, "path": "/index.php", "ttl": 300, "notFoundTtl": 300, "type": "exact", "enforce": false, "enabled": false}]}`,
		),
	})
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

			if v.Enabled != true {
				t.Errorf("Expected to get cache setting with Enabled [%t] but got [%t]", true, v.Enabled)
			}

		}

		if v.ID == 2 {
			if v.Path != "/index.php" {
				t.Errorf("Expected to get cache setting with Type [%s] but got [%s]", "/index.php", v.Path)
			}

			if v.Enabled != false {
				t.Errorf("Expected to get cache setting with Enabled [%t] but got [%t]", false, v.Enabled)
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
