package myrasec

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const domainsResponse = `{"error": false, "data": [{"id": 1, "name": "example.com"}]}`

func TestCheckRedirect(t *testing.T) {
	first := func(raw string) *http.Request {
		return &http.Request{URL: mustParseURL(t, raw)}
	}

	tests := []struct {
		name    string
		from    string
		to      string
		wantErr string
	}{
		{"same host https", "https://api.example.test/domains", "https://api.example.test/v2/domains", ""},
		{"same host different case", "https://api.example.test/domains", "https://API.example.test/domains", ""},
		{"downgrade to http", "https://api.example.test/domains", "http://api.example.test/domains", "downgrade"},
		{"other host", "https://api.example.test/domains", "https://evil.example.test/domains", "different host"},
		{"other port", "https://api.example.test/domains", "https://api.example.test:8443/domains", "different host"},
		{"http to http same host", "http://127.0.0.1:8080/domains", "http://127.0.0.1:8080/other", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkRedirect(first(tt.to), []*http.Request{first(tt.from)})
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Expected the redirect to be allowed but got %v", err)
				}

				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Expected an error containing %q but got %v", tt.wantErr, err)
			}
		})
	}
}

func TestCheckRedirectLimit(t *testing.T) {
	via := make([]*http.Request, maxRedirects)
	for i := range via {
		via[i] = &http.Request{URL: mustParseURL(t, "https://api.example.test/domains")}
	}

	err := checkRedirect(&http.Request{URL: mustParseURL(t, "https://api.example.test/domains")}, via)
	if err == nil || !strings.Contains(err.Error(), "redirects") {
		t.Fatalf("Expected the redirect limit to stop the chain but got %v", err)
	}
}

func TestRedirectDowngradeIsRefused(t *testing.T) {
	var hits int
	var target string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}))
	t.Cleanup(server.Close)
	// The same host and port the client connected to, plain HTTP instead of HTTPS.
	target = "http://" + strings.TrimPrefix(server.URL, "https://") + "/domains"

	api := newTLSTestAPI(t, server)

	_, err := api.ListDomainsContext(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "downgrade") {
		t.Fatalf("Expected the downgrade to be refused but got %v", err)
	}

	if hits != 1 {
		t.Errorf("Expected a single request before the refused redirect but got %d", hits)
	}
}

func TestRedirectToOtherHostIsRefused(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "localhost" is another host name than the 127.0.0.1 the client connected to.
		http.Redirect(w, r, "https://localhost/domains", http.StatusFound)
	}))
	t.Cleanup(server.Close)

	api := newTLSTestAPI(t, server)

	_, err := api.ListDomainsContext(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "different host") {
		t.Fatalf("Expected the cross host redirect to be refused but got %v", err)
	}
}

func TestSameHostRedirectIsFollowed(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/domains" {
			http.Redirect(w, r, "/v2/domains", http.StatusFound)
			return
		}

		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("Expected the Authorization header to survive the same host redirect")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(domainsResponse))
	}))
	t.Cleanup(server.Close)

	api := newTLSTestAPI(t, server)

	domains, err := api.ListDomainsContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected the same host redirect to be followed but got %v", err)
	}

	if len(domains) != 1 {
		t.Errorf("Expected one domain but got %d", len(domains))
	}
}

func TestPlainHTTPIsRefusedByDefault(t *testing.T) {
	const token = "s3cr3t-api-token-value"

	api, err := NewWithToken(token)
	if err != nil {
		t.Fatalf("Unexpected error creating the API client: %v", err)
	}
	// No connection is attempted: the scheme check happens before the request is sent.
	api.BaseURL = "http://api.example.invalid/%s"

	_, err = api.ListDomainsContext(context.Background(), nil)
	if !errors.Is(err, ErrInsecureEndpoint) {
		t.Fatalf("Expected ErrInsecureEndpoint but got %v", err)
	}

	if strings.Contains(err.Error(), token) {
		t.Errorf("Expected the error not to leak the credentials but got %q", err)
	}
}

func TestPlainHTTPAllowedWithOptIn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(domainsResponse))
	}))
	t.Cleanup(server.Close)

	api, err := NewWithToken("token")
	if err != nil {
		t.Fatalf("Unexpected error creating the API client: %v", err)
	}
	// A non loopback host name whose connections are routed to the local test server, so the
	// opt in is exercised without touching the network.
	api.BaseURL = "http://api.example.invalid/%s"
	if err := api.SetHTTPClient(&http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, server.Listener.Addr().String())
		},
	}}); err != nil {
		t.Fatalf("Unexpected error setting the HTTP client: %v", err)
	}

	if _, err := api.ListDomainsContext(context.Background(), nil); !errors.Is(err, ErrInsecureEndpoint) {
		t.Fatalf("Expected ErrInsecureEndpoint without the opt in but got %v", err)
	}

	api.SetInsecureHTTP(true)

	domains, err := api.ListDomainsContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected the request to succeed after the opt in but got %v", err)
	}

	if len(domains) != 1 {
		t.Errorf("Expected one domain but got %d", len(domains))
	}

	api.SetInsecureHTTP(false)

	if _, err := api.ListDomainsContext(context.Background(), nil); !errors.Is(err, ErrInsecureEndpoint) {
		t.Fatalf("Expected ErrInsecureEndpoint after revoking the opt in but got %v", err)
	}
}

func TestPlainHTTPToLoopbackIsAllowed(t *testing.T) {
	// newTestAPI points the client at a plain HTTP server on 127.0.0.1.
	api, _ := newTestAPI(t, map[string]testResponse{
		"GET /domains": {Status: http.StatusOK, Body: domainsResponse},
	})

	if _, err := api.ListDomainsContext(context.Background(), nil); err != nil {
		t.Fatalf("Expected plain HTTP to the loopback address to be allowed but got %v", err)
	}
}

func TestCheckScheme(t *testing.T) {
	tests := []struct {
		url      string
		insecure bool
		wantErr  error
		wantAny  bool
	}{
		{"https://api.example.test/domains", false, nil, false},
		{"http://localhost/domains", false, nil, false},
		{"http://LOCALHOST:8080/domains", false, nil, false},
		{"http://127.0.0.1:8080/domains", false, nil, false},
		{"http://[::1]:8080/domains", false, nil, false},
		{"http://api.example.test/domains", false, ErrInsecureEndpoint, false},
		{"http://10.0.0.1/domains", false, ErrInsecureEndpoint, false},
		{"http://api.example.test/domains", true, nil, false},
		{"ftp://api.example.test/domains", true, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			api := &API{insecureHTTP: tt.insecure}

			err := api.checkScheme(mustParseURL(t, tt.url))
			switch {
			case tt.wantAny:
				if err == nil {
					t.Fatal("Expected an error for the unsupported scheme")
				}
			case tt.wantErr != nil:
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Expected %v but got %v", tt.wantErr, err)
				}
			default:
				if err != nil {
					t.Fatalf("Expected no error but got %v", err)
				}
			}
		})
	}
}

func TestRetryableStatus(t *testing.T) {
	tests := []struct {
		method string
		status int
		want   bool
	}{
		{http.MethodGet, http.StatusOK, false},
		{http.MethodGet, http.StatusBadRequest, false},
		{http.MethodGet, http.StatusForbidden, false},
		{http.MethodGet, http.StatusNotFound, false},
		{http.MethodGet, http.StatusTooManyRequests, true},
		{http.MethodPost, http.StatusTooManyRequests, true},
		{http.MethodPut, http.StatusTooManyRequests, true},
		{http.MethodDelete, http.StatusTooManyRequests, true},
		{http.MethodGet, http.StatusServiceUnavailable, true},
		{http.MethodPost, http.StatusServiceUnavailable, true},
		{http.MethodGet, http.StatusInternalServerError, true},
		{http.MethodPost, http.StatusInternalServerError, false},
		{http.MethodPut, http.StatusInternalServerError, false},
		{http.MethodDelete, http.StatusInternalServerError, false},
		{http.MethodGet, http.StatusBadGateway, true},
		{http.MethodPost, http.StatusBadGateway, false},
		{http.MethodPut, http.StatusBadGateway, false},
		{http.MethodGet, http.StatusGatewayTimeout, true},
		{http.MethodPost, http.StatusGatewayTimeout, false},
		{http.MethodDelete, http.StatusGatewayTimeout, false},
	}

	for _, tt := range tests {
		if got := retryableStatus(tt.method, tt.status); got != tt.want {
			t.Errorf("retryableStatus(%s, %d) = %t, want %t", tt.method, tt.status, got, tt.want)
		}
	}
}

func TestRetryDelay(t *testing.T) {
	api := &API{retrySleep: 1}

	response := func(retryAfter string) http.Header {
		header := http.Header{}
		if retryAfter != "" {
			header.Set("Retry-After", retryAfter)
		}

		return header
	}

	tests := []struct {
		name       string
		status     int
		retryAfter string
		attempt    int
		want       time.Duration
	}{
		{"first attempt", http.StatusBadGateway, "", 1, time.Second},
		{"second attempt", http.StatusBadGateway, "", 2, 2 * time.Second},
		{"third attempt", http.StatusBadGateway, "", 3, 4 * time.Second},
		{"growth is capped", http.StatusBadGateway, "", 20, maxRetryAfter},
		{"zero attempt", http.StatusBadGateway, "", 0, time.Second},
		{"retry after seconds", http.StatusTooManyRequests, "3", 1, 3 * time.Second},
		{"retry after zero is honoured", http.StatusTooManyRequests, "0", 2, 0},
		{"retry after negative", http.StatusTooManyRequests, "-5", 2, 0},
		{"retry after is capped", http.StatusTooManyRequests, "3600", 1, maxRetryAfter},
		{"retry after does not overflow", http.StatusTooManyRequests, "9223372036854775807", 1, maxRetryAfter},
		{"retry after garbage", http.StatusTooManyRequests, "soon", 2, 2 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := api.retryDelay(tt.status, response(tt.retryAfter), tt.attempt); got != tt.want {
				t.Errorf("retryDelay = %s, want %s", got, tt.want)
			}
		})
	}

	t.Run("retry after date", func(t *testing.T) {
		at := time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)

		got := api.retryDelay(http.StatusTooManyRequests, response(at), 1)
		if got < 25*time.Second || got > 30*time.Second {
			t.Errorf("retryDelay = %s, want about 30s", got)
		}
	})

	t.Run("retry after date in the past", func(t *testing.T) {
		at := time.Now().Add(-30 * time.Second).UTC().Format(http.TimeFormat)

		if got := api.retryDelay(http.StatusTooManyRequests, response(at), 1); got != 0 {
			t.Errorf("retryDelay = %s, want 0", got)
		}
	})

	t.Run("no sleep configured", func(t *testing.T) {
		api := &API{retrySleep: 0}

		if got := api.retryDelay(http.StatusBadGateway, response(""), 3); got != 0 {
			t.Errorf("retryDelay = %s, want 0", got)
		}

		for _, status := range []int{http.StatusTooManyRequests, http.StatusServiceUnavailable} {
			if got := api.retryDelay(status, response(""), 1); got != time.Second {
				t.Errorf("retryDelay for %d = %s, want 1s", status, got)
			}
		}
	})

	t.Run("long sleep is not doubled", func(t *testing.T) {
		api := &API{retrySleep: 90}

		if got := api.retryDelay(http.StatusBadGateway, response(""), 3); got != 90*time.Second {
			t.Errorf("retryDelay = %s, want 90s", got)
		}
	})
}

func TestRetryOnTooManyRequests(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(domainsResponse))
	})
	api.SetMaxRetries(3)
	api.SetRetrySleep(5)

	start := time.Now()
	domains, err := api.ListDomainsContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected the retries to succeed but got %v", err)
	}

	if len(domains) != 1 || hits != 3 {
		t.Errorf("Expected one domain after three requests but got %d domains after %d requests", len(domains), hits)
	}

	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Expected Retry-After: 0 to override the configured sleep, the call took %s", elapsed)
	}
}

func TestRetryOnServiceUnavailable(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(domainsResponse))
	})
	api.SetMaxRetries(2)
	api.SetRetrySleep(0)

	if _, err := api.ListDomainsContext(context.Background(), nil); err != nil {
		t.Fatalf("Expected the retry to succeed but got %v", err)
	}

	if hits != 2 {
		t.Errorf("Expected two requests but got %d", hits)
	}
}

func TestRetryOnBadGatewayDependsOnMethod(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusBadGateway)
	})
	api.SetMaxRetries(3)
	api.SetRetrySleep(0)

	_, err := api.ListDomainsContext(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected an error after the retries were exhausted")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("Expected an APIError with status 502 but got %v", err)
	}

	if hits != 3 {
		t.Errorf("Expected three attempts for the GET but got %d", hits)
	}

	hits = 0

	_, err = api.CreateDomainContext(context.Background(), &Domain{Name: "example.com"})
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("Expected an APIError with status 502 but got %v", err)
	}

	if hits != 1 {
		t.Errorf("Expected a single attempt for the POST but got %d", hits)
	}
}

func TestSingleAttemptByDefault(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	var apiErr *APIError
	if _, err := api.ListDomainsContext(context.Background(), nil); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("Expected an APIError with status 429 but got %v", err)
	}

	if hits != 1 {
		t.Errorf("Expected a single attempt with the default retry count but got %d", hits)
	}
}

func TestNoRetryOnServerErrorForWrites(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusInternalServerError)
	})
	api.SetMaxRetries(3)
	api.SetRetrySleep(0)

	var apiErr *APIError
	if _, err := api.CreateDomainContext(context.Background(), &Domain{Name: "example.com"}); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected an APIError with status 500 but got %v", err)
	}

	if hits != 1 {
		t.Errorf("Expected a single attempt for the POST but got %d", hits)
	}

	hits = 0

	if _, err := api.ListDomainsContext(context.Background(), nil); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected an APIError with status 500 but got %v", err)
	}

	if hits != 3 {
		t.Errorf("Expected three attempts for the GET but got %d", hits)
	}
}

func TestRetryWaitStopsOnContextCancel(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	api.SetMaxRetries(3)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := api.ListDomainsContext(ctx, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected the context error but got %v", err)
	}

	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("Expected the wait to stop with the context, the call took %s", elapsed)
	}

	if hits != 1 {
		t.Errorf("Expected a single attempt before the context expired but got %d", hits)
	}
}

func TestNoRetryOnClientError(t *testing.T) {
	var hits int
	api := newTestAPIWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusBadRequest)
	})
	api.SetMaxRetries(3)
	api.SetRetrySleep(0)

	if _, err := api.ListDomainsContext(context.Background(), nil); err == nil {
		t.Fatal("Expected an error")
	}

	if hits != 1 {
		t.Errorf("Expected a single attempt but got %d", hits)
	}
}

func TestSetHTTPClientAppliesRedirectPolicy(t *testing.T) {
	api, err := NewWithToken("token")
	if err != nil {
		t.Fatalf("Unexpected error creating the API client: %v", err)
	}

	if api.client.CheckRedirect == nil {
		t.Fatal("Expected the default client to carry the redirect policy")
	}

	plain := &http.Client{Timeout: time.Second}
	if err := api.SetHTTPClient(plain); err != nil {
		t.Fatalf("Unexpected error setting the HTTP client: %v", err)
	}

	if api.client.CheckRedirect == nil {
		t.Error("Expected the client without a policy to get the redirect policy")
	}

	if api.client == plain || plain.CheckRedirect != nil {
		t.Error("Expected the caller's client to stay untouched")
	}

	if api.client.Timeout != time.Second {
		t.Error("Expected the copy to keep the settings of the passed client")
	}

	called := false
	custom := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		called = true

		return nil
	}}
	if err := api.SetHTTPClient(custom); err != nil {
		t.Fatalf("Unexpected error setting the HTTP client: %v", err)
	}

	if api.client != custom {
		t.Error("Expected a client with its own policy to be used as is")
	}

	_ = api.client.CheckRedirect(nil, nil)
	if !called {
		t.Error("Expected the custom policy to stay in place")
	}
}

// newTLSTestAPI returns an API client that trusts the certificate of the passed TLS test
// server and is pointed at it.
func newTLSTestAPI(t *testing.T, server *httptest.Server) *API {
	t.Helper()

	api, err := NewWithToken("token")
	if err != nil {
		t.Fatalf("Unexpected error creating the API client: %v", err)
	}

	if err := api.SetHTTPClient(server.Client()); err != nil {
		t.Fatalf("Unexpected error setting the HTTP client: %v", err)
	}
	api.BaseURL = server.URL + "/%s"

	return api
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("Unexpected error parsing %q: %v", raw, err)
	}

	return u
}
