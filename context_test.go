package myrasec

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

type contextTestKey struct{}

func TestPrepareRequestCarriesContext(t *testing.T) {
	api, err := New("abc123", "123abc")
	if err != nil {
		t.Fatal("Unexpected error")
	}

	ctx := context.WithValue(context.Background(), contextTestKey{}, "value")

	for _, definition := range []APIMethod{
		api.methods["listDomains"],
		api.methods["createDomain"],
		api.methods["updateDomain"],
		api.methods["deleteDomain"],
	} {
		definition.Action = "domains"

		req, err := api.prepareRequest(ctx, definition, map[string]string{})
		if err != nil {
			t.Fatalf("Unexpected error for %s: %v", definition.Method, err)
		}

		if req.Context().Value(contextTestKey{}) != "value" {
			t.Errorf("Expected the %s request to carry the passed context", definition.Method)
		}
	}
}

func TestCallWithCancelledContext(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /domains": {Status: http.StatusOK, Body: `{"error": false, "data": []}`},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := api.ListDomainsContext(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled but got %v", err)
	}

	if len(requests.all()) != 0 {
		t.Error("Expected no request to be sent for a cancelled context")
	}
}

func TestRateLimiterRespectsDeadline(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /domains": {Status: http.StatusOK, Body: `{"error": false, "data": []}`},
	})

	// One token, the next one arrives in an hour.
	api.limiter = rate.NewLimiter(rate.Every(time.Hour), 1)

	if _, err := api.ListDomainsContext(context.Background(), nil); err != nil {
		t.Fatalf("Expected the first call to succeed but got %v", err)
	}

	// The deadline is generous for the test but far below the wait the limiter needs,
	// so the limiter has to fail without waiting.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	_, err := api.ListDomainsContext(ctx, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected context.DeadlineExceeded but got %v", err)
	}

	if !strings.Contains(err.Error(), ErrorMsgRateLimitReached) {
		t.Errorf("Expected the error to name the rate limit but got [%s]", err.Error())
	}

	if time.Since(start) > time.Second {
		t.Error("Expected the limiter to fail fast when the wait exceeds the deadline")
	}

	if len(requests.all()) != 1 {
		t.Errorf("Expected exactly one request but got %d", len(requests.all()))
	}
}

func TestRetrySleepRespectsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The server answers 500 and cancels the context while the client is about to
	// sleep before the retry, so the test does not depend on wall clock timing.
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusInternalServerError)
		cancel()
	})
	api.SetMaxRetries(3)
	api.SetRetrySleep(5)

	start := time.Now()
	_, err := api.ListDomainsContext(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled but got %v", err)
	}

	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Expected the retry sleep to be interrupted, the call took %s", elapsed)
	}

	if hits != 1 {
		t.Errorf("Expected one request before the interrupted retry but got %d", hits)
	}
}

func TestRetrySucceedsWhileContextIsAlive(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"error": false, "data": [{"id": 1, "name": "example.com"}]}`))
	})
	api.SetMaxRetries(2)
	api.SetRetrySleep(0)

	domains, err := api.ListDomainsContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected the retry to succeed but got %v", err)
	}

	if len(domains) != 1 || hits != 2 {
		t.Errorf("Expected one domain after two requests but got %d domains after %d requests", len(domains), hits)
	}
}

func TestCacheHitRespectsCancelledContext(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /domains": {Status: http.StatusOK, Body: `{"error": false, "data": [{"id": 1, "name": "example.com"}]}`},
	})
	api.EnableCaching()

	if _, err := api.ListDomainsContext(context.Background(), nil); err != nil {
		t.Fatalf("Expected the first call to succeed but got %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := api.ListDomainsContext(ctx, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled on a cache hit but got %v", err)
	}

	if len(requests.all()) != 1 {
		t.Errorf("Expected the cached response not to trigger a request, got %d requests", len(requests.all()))
	}
}

func TestContextPropagatesThroughInternalCalls(t *testing.T) {
	routes := map[string]testResponse{
		"GET /user/me":             {Status: http.StatusOK, Body: `{"error": false, "data": [{"objectType": "UserExtendedVO", "id": 12345, "login": "test@example.com"}]}`},
		"GET /user/12345/api-keys": {Status: http.StatusOK, Body: `{"error": false, "data": [{"objectType": "ApiKeyVO", "id": 1, "name": "key"}]}`},
	}

	api, requests := newTestAPI(t, routes)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// ListApiKeys resolves the user first. The cancelled context must stop that call as well.
	if _, err := api.ListApiKeysContext(ctx, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled but got %v", err)
	}

	if len(requests.all()) != 0 {
		t.Errorf("Expected no request but got %d", len(requests.all()))
	}

	keys, err := api.ListApiKeysContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected not to get an error but got %v", err)
	}

	if len(keys) != 1 || len(requests.all()) != 2 {
		t.Errorf("Expected one key from two requests but got %d keys from %d requests", len(keys), len(requests.all()))
	}
}

func TestDeprecatedMethodDelegatesToContextVariant(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /domains": {Status: http.StatusOK, Body: `{"error": false, "data": [{"id": 1, "name": "example.com"}]}`},
	})

	domains, err := api.ListDomains(nil)
	if err != nil {
		t.Fatalf("Expected not to get an error but got %v", err)
	}

	if len(domains) != 1 || domains[0].Name != "example.com" {
		t.Errorf("Expected the domain list, got %+v", domains)
	}

	if sent := requests.last(t); sent.Path != "/domains" {
		t.Errorf("Expected GET /domains but got %s", sent.Path)
	}
}
