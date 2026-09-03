package myrasec

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

// testRequest records one request the test server received.
type testRequest struct {
	Method string
	Path   string
	Query  url.Values
	Body   []byte
}

// testResponse is a canned answer of the test server.
type testResponse struct {
	Status int
	Body   string
}

// testRequests collects the requests the test server received.
type testRequests struct {
	mu       sync.Mutex
	requests []testRequest
}

// all returns a copy of the recorded requests.
func (r *testRequests) all() []testRequest {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]testRequest(nil), r.requests...)
}

// last returns the most recent recorded request and fails the test when there is none.
func (r *testRequests) last(t *testing.T) testRequest {
	t.Helper()

	requests := r.all()
	if len(requests) == 0 {
		t.Fatal("Expected the test server to have received a request")
	}

	return requests[len(requests)-1]
}

// newTestAPI starts an HTTP test server answering the passed routes ("METHOD /path") and
// returns an API client pointed at it together with the recorded requests. Unknown routes
// answer 404 with an empty body.
func newTestAPI(t *testing.T, routes map[string]testResponse) (*API, *testRequests) {
	t.Helper()

	recorded := &testRequests{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		recorded.mu.Lock()
		recorded.requests = append(recorded.requests, testRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.Query(),
			Body:   body,
		})
		recorded.mu.Unlock()

		route, ok := routes[r.Method+" "+r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(route.Status)
		_, _ = w.Write([]byte(route.Body))
	}))
	t.Cleanup(server.Close)

	api, err := NewWithToken("token")
	if err != nil {
		t.Fatalf("Unexpected error creating the API client: %v", err)
	}
	api.BaseURL = server.URL + "/%s"

	return api, recorded
}
