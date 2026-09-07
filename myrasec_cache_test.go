package myrasec

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestIsExpired(t *testing.T) {
	expired := responseCache{
		Key:     "1",
		Created: time.Now().Unix(),
		Expire:  time.Now().Unix() - 1,
		Body:    nil,
	}

	if !expired.isExpired() {
		t.Errorf("Expected that the cache is expired")
	}

	valid := responseCache{
		Key:     "1",
		Created: time.Now().Unix(),
		Expire:  time.Now().Add(time.Second * 1).Unix(),
		Body:    nil,
	}

	if valid.isExpired() {
		t.Errorf("Expected that the cache is valid/not expired")
	}
}

func TestInCache(t *testing.T) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://apiv2.myracloud.com/domains", nil)

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

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
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

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
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

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	if !isCachable(req) {
		t.Errorf("Expected the request to be cachable as it is [%s]", req.Method)
	}

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodDelete, "https://apiv2.myracloud.com/domains", nil)
	if isCachable(req) {
		t.Errorf("Expected the request not to be cachable as it is [%s]", req.Method)
	}

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodPost, "https://apiv2.myracloud.com/domains", nil)
	if isCachable(req) {
		t.Errorf("Expected the request not to be cachable as it is [%s]", req.Method)
	}

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodPut, "https://apiv2.myracloud.com/domains", nil)
	if isCachable(req) {
		t.Errorf("Expected the request not to be cachable as it is [%s]", req.Method)
	}
}

func TestRemoveFromCache(t *testing.T) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
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

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
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

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	sha := BuildSHA256(req.URL.String())
	key := BuildCacheKey(req)

	if sha != key {
		t.Errorf("Expected to have key and sha the same value")
	}
}

// cacheTestServer answers every route of a DNS record with a canned body and records
// the methods it received.
func cacheTestServer(t *testing.T) (*API, func() []string) {
	t.Helper()

	return cacheTestServerWithHandler(t, nil)
}

// cacheTestServerWithHandler is cacheTestServer with a handler that runs after the request
// was recorded. The canned body is skipped when the handler returns true.
func cacheTestServerWithHandler(t *testing.T, handler func(w http.ResponseWriter, r *http.Request) bool) (*API, func() []string) {
	t.Helper()

	var mu sync.Mutex
	var received []string

	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		received = append(received, r.Method+" "+r.URL.Path)
		mu.Unlock()

		if handler != nil && handler(w, r) {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			_, _ = w.Write([]byte(`{"error":false,"data":[{"id":5,"name":"www.example.com.","value":"1.1.1.1","recordType":"A","ttl":300}]}`))
		default:
			_, _ = w.Write([]byte(`{"error":false,"targetObject":[{"id":5,"name":"www.example.com.","value":"2.2.2.2","recordType":"A","ttl":300}]}`))
		}
	})
	api.EnableCaching()
	api.SetCachingTTL(60)

	return api, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), received...)
	}
}

// TestCacheServesRepeatedGET is the positive case: the second GET of the same URL is
// answered from the cache.
func TestCacheServesRepeatedGET(t *testing.T) {
	api, received := cacheTestServer(t)
	ctx := context.Background()

	for range 2 {
		if _, err := api.GetDNSRecordContext(ctx, 1, 5); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	if got := received(); len(got) != 1 {
		t.Errorf("Expected the server to receive one GET but got %v", got)
	}
}

// TestCacheNeverServesWrites is the regression test for writes on a URL that was read
// before: PUT and DELETE share the URL of the GET and must reach the API nonetheless.
func TestCacheNeverServesWrites(t *testing.T) {
	api, received := cacheTestServer(t)
	ctx := context.Background()

	rec, err := api.GetDNSRecordContext(ctx, 1, 5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	updated, err := api.UpdateDNSRecordContext(ctx, rec, 1)
	if err != nil {
		t.Fatalf("Unexpected error updating: %v", err)
	}
	if updated.Value != "2.2.2.2" {
		t.Errorf("Expected the update to return the API response but got value [%s]", updated.Value)
	}

	if _, err := api.DeleteDNSRecordContext(ctx, rec, 1); err != nil {
		t.Fatalf("Unexpected error deleting: %v", err)
	}

	want := []string{"GET /domain/1/dns-records/5", "PUT /domain/1/dns-records/5", "DELETE /domain/1/dns-records/5"}
	if got := received(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected the server to receive %v but got %v", want, got)
	}
}

// TestCacheNeverServesCreate covers the POST on a collection URL that was listed before.
func TestCacheNeverServesCreate(t *testing.T) {
	api, received := cacheTestServer(t)
	ctx := context.Background()

	if _, err := api.ListDNSRecordsContext(ctx, 1, map[string]string{}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	created, err := api.CreateDNSRecordContext(ctx, &DNSRecord{Name: "www.example.com.", Value: "2.2.2.2", RecordType: "A", TTL: 300}, 1)
	if err != nil {
		t.Fatalf("Unexpected error creating: %v", err)
	}
	if created.ID != 5 {
		t.Errorf("Expected the created record from the API response but got ID [%d]", created.ID)
	}

	want := []string{"GET /domain/1/dns-records", "POST /domain/1/dns-records"}
	if got := received(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected the server to receive %v but got %v", want, got)
	}
}

// TestWriteFlushesCache checks that a read after a write is not served from the cache.
func TestWriteFlushesCache(t *testing.T) {
	api, received := cacheTestServer(t)
	ctx := context.Background()

	rec, err := api.GetDNSRecordContext(ctx, 1, 5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := api.UpdateDNSRecordContext(ctx, rec, 1); err != nil {
		t.Fatalf("Unexpected error updating: %v", err)
	}
	if _, err := api.GetDNSRecordContext(ctx, 1, 5); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	want := []string{"GET /domain/1/dns-records/5", "PUT /domain/1/dns-records/5", "GET /domain/1/dns-records/5"}
	if got := received(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected the server to receive %v but got %v", want, got)
	}
}

// TestRejectedWriteFlushesCache checks that a write the API answered with an error
// flushes the cache like a successful one.
func TestRejectedWriteFlushesCache(t *testing.T) {
	api, received := cacheTestServerWithHandler(t, func(w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodPut {
			return false
		}

		w.WriteHeader(http.StatusBadRequest)
		return true
	})
	ctx := context.Background()

	rec, err := api.GetDNSRecordContext(ctx, 1, 5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := api.UpdateDNSRecordContext(ctx, rec, 1); err == nil {
		t.Fatal("Expected the rejected update to return an error")
	}
	if _, err := api.GetDNSRecordContext(ctx, 1, 5); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	want := []string{"GET /domain/1/dns-records/5", "PUT /domain/1/dns-records/5", "GET /domain/1/dns-records/5"}
	if got := received(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected the server to receive %v but got %v", want, got)
	}
}

// TestWriteWithoutAnswerFlushesCache covers a write whose answer never arrives: the API
// can have applied it nonetheless, so the cache is flushed.
func TestWriteWithoutAnswerFlushesCache(t *testing.T) {
	api, received := cacheTestServerWithHandler(t, func(w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodPut {
			return false
		}

		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Errorf("Unexpected error hijacking the connection: %v", err)
			return true
		}
		_ = conn.Close()

		return true
	})
	ctx := context.Background()

	rec, err := api.GetDNSRecordContext(ctx, 1, 5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := api.UpdateDNSRecordContext(ctx, rec, 1); err == nil {
		t.Fatal("Expected the update without an answer to return an error")
	}
	if _, err := api.GetDNSRecordContext(ctx, 1, 5); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	want := []string{"GET /domain/1/dns-records/5", "PUT /domain/1/dns-records/5", "GET /domain/1/dns-records/5"}
	if got := received(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected the server to receive %v but got %v", want, got)
	}
}

// TestWriteWhileCachingDisabledFlushesCache checks that an entry cached before
// DisableCaching is not served after EnableCaching when a write happened in between.
func TestWriteWhileCachingDisabledFlushesCache(t *testing.T) {
	api, received := cacheTestServer(t)
	ctx := context.Background()

	rec, err := api.GetDNSRecordContext(ctx, 1, 5)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	api.DisableCaching()
	if _, err := api.UpdateDNSRecordContext(ctx, rec, 1); err != nil {
		t.Fatalf("Unexpected error updating: %v", err)
	}
	api.EnableCaching()

	if _, err := api.GetDNSRecordContext(ctx, 1, 5); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	want := []string{"GET /domain/1/dns-records/5", "PUT /domain/1/dns-records/5", "GET /domain/1/dns-records/5"}
	if got := received(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected the server to receive %v but got %v", want, got)
	}
}

// TestReadInFlightDuringWriteIsNotCached covers a GET that was sent before a write and is
// answered after it. Its response holds the state from before the write and must not be
// cached, the flush of the write happened before it arrived.
func TestReadInFlightDuringWriteIsNotCached(t *testing.T) {
	readSent := make(chan struct{})
	releaseRead := make(chan struct{})
	var firstRead sync.Once

	api, received := cacheTestServerWithHandler(t, func(w http.ResponseWriter, r *http.Request) bool {
		if r.Method == http.MethodGet {
			firstRead.Do(func() {
				close(readSent)
				<-releaseRead
			})
		}

		return false
	})
	ctx := context.Background()

	readDone := make(chan error, 1)
	go func() {
		_, err := api.GetDNSRecordContext(ctx, 1, 5)
		readDone <- err
	}()

	<-readSent
	_, err := api.UpdateDNSRecordContext(ctx, &DNSRecord{ID: 5, Name: "www.example.com.", Value: "2.2.2.2", RecordType: "A", TTL: 300}, 1)
	close(releaseRead)
	if err != nil {
		t.Fatalf("Unexpected error updating: %v", err)
	}
	if err := <-readDone; err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, err := api.GetDNSRecordContext(ctx, 1, 5); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	want := []string{"GET /domain/1/dns-records/5", "PUT /domain/1/dns-records/5", "GET /domain/1/dns-records/5"}
	if got := received(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected the server to receive %v but got %v", want, got)
	}
}
