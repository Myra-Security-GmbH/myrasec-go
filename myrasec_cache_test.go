package myrasec

import (
	"net/http"
	"testing"
	"time"
)

func TestIsExpired(t *testing.T) {
	expired := responseCache{
		Key:     "1",
		Created: time.Now().Unix(),
		Expire:  time.Now().Unix() - 1,
		Request: nil,
		Body:    nil,
	}

	if !expired.isExpired() {
		t.Errorf("Expected that the cache is expired")
	}

	valid := responseCache{
		Key:     "1",
		Created: time.Now().Unix(),
		Expire:  time.Now().Add(time.Second * 1).Unix(),
		Request: nil,
		Body:    nil,
	}

	if valid.isExpired() {
		t.Errorf("Expected that the cache is valid/not expired")
	}
}

func TestInCache(t *testing.T) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)

	if api.inCache(req) {
		t.Errorf("Expected not to find the passed request in the cache")
	}

	api.cacheResponse(req, "CONTENT")

	if !api.inCache(req) {
		t.Errorf("Expected to find the passed request in the cache")
	}

	for k := range api.cache {
		api.RemoveFromCache(k)
	}

	if api.inCache(req) {
		t.Errorf("Expected not to find the passed request in the cache (cleared cache)")
	}
}

func TestFromCache(t *testing.T) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	api.cacheResponse(req, "CONTENT")

	v := api.fromCache(req)
	if v != "CONTENT" {
		t.Errorf("Expected to get [%s] but got [%s]", "CONTENT", v)
	}

	for k := range api.cache {
		api.RemoveFromCache(k)
	}

	api.cacheTTL = -10
	api.cacheResponse(req, "CONTENT")

	v = api.fromCache(req)
	if v != nil {
		t.Errorf("Expected not to get a expired cache result")
	}

	if len(api.cache) > 0 {
		t.Errorf("Expected not to have any element in the cache")
	}
}

func TestCacheResponse(t *testing.T) {
	api, _ := New("abc123", "123abc")

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	api.cacheResponse(req, "CONTENT")

	if len(api.cache) > 0 {
		t.Errorf("Expected not to have anything in the cache as caching is not enabled")
	}

	api.EnableCaching()
	api.cacheResponse(req, "CONTENT")

	if len(api.cache) != 1 {
		t.Errorf("Expected to have one single element in the cache but got %d", len(api.cache))
	}

}

func TestIsCachable(t *testing.T) {
	var req *http.Request

	req, _ = http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	if !isCachable(req) {
		t.Errorf("Expected the request to be cachable as it is [%s]", req.Method)
	}

	req, _ = http.NewRequest(http.MethodDelete, "https://apiv2.myracloud.com/domains", nil)
	if isCachable(req) {
		t.Errorf("Expected the request not to be cachable as it is [%s]", req.Method)
	}

	req, _ = http.NewRequest(http.MethodPost, "https://apiv2.myracloud.com/domains", nil)
	if isCachable(req) {
		t.Errorf("Expected the request not to be cachable as it is [%s]", req.Method)
	}

	req, _ = http.NewRequest(http.MethodPut, "https://apiv2.myracloud.com/domains", nil)
	if isCachable(req) {
		t.Errorf("Expected the request not to be cachable as it is [%s]", req.Method)
	}

}

func TestRemoveFromCache(t *testing.T) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	api.cacheResponse(req, "CONTENT")

	if len(api.cache) != 1 {
		t.Errorf("Expected to have a single element in the cache")
	}

	for k := range api.cache {
		api.RemoveFromCache(k)
	}

	if len(api.cache) != 0 {
		t.Errorf("Expected not to have any element in the cache")
	}

}

func TestPruneCache(t *testing.T) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	api.cacheResponse(req, "CONTENT")

	if len(api.cache) != 1 {
		t.Errorf("Expected to have a single element in the cache")
	}

	api.PruneCache()

	if len(api.cache) != 0 {
		t.Errorf("Expected not to have any element in the cache")
	}
}

func TestBuildCacheKey(t *testing.T) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	sha := BuildSHA256(req.URL.String())
	key := BuildCacheKey(req)

	if sha != key {
		t.Errorf("Expected to have key and sha the same value")
	}
}
