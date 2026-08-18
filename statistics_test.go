package myrasec

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

func TestDecodeStatisticsResponse(t *testing.T) {
	body := `{"query": {"startDate": "2026-08-01T00:00:00+0200", "endDate": "2026-08-08T00:00:00+0200", "type": "fqdn", "fqdn": ["www.example.com"], "dataSources": {"traffic": {"source": "bytes", "type": "stats"}}, "aggregationInterval": "day"},
		"result": {"traffic": {"min": 0, "max": 10, "avg": 5, "sum": 4096}}}`
	resp := &http.Response{Body: io.NopCloser(strings.NewReader(body))}

	res, err := decodeStatisticsResponse(resp, methods["queryStatistics"])
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}
	stats, ok := res.(*Statistics)
	if !ok {
		t.Fatalf("Expected *Statistics but got %T", res)
	}
	if stats.Query == nil || stats.Query.Type != "fqdn" || len(stats.Query.FQDN) != 1 {
		t.Errorf("Expected the query to be echoed, got %+v", stats.Query)
	}
	traffic, ok := stats.Result["traffic"].(map[string]any)
	if !ok || traffic["sum"] != float64(4096) {
		t.Errorf("Expected result traffic.sum 4096, got %+v", stats.Result)
	}
}

func TestDecodeStatisticsResponseShapes(t *testing.T) {
	// The API serializes an empty PHP array as [] and may omit query or result;
	// none of these is an error.
	cases := map[string]string{
		"result empty array":   `{"query": {"type": "fqdn"}, "result": []}`,
		"result omitted":       `{"query": {"type": "fqdn"}}`,
		"result null":          `{"query": {"type": "fqdn"}, "result": null}`,
		"query omitted":        `{"result": {}}`,
		"error false explicit": `{"error": false, "query": {"type": "fqdn"}, "result": {}}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			resp := &http.Response{Body: io.NopCloser(strings.NewReader(body))}
			res, err := decodeStatisticsResponse(resp, methods["queryStatistics"])
			if err != nil {
				t.Fatalf("Expected not to get an error but got [%s]", err.Error())
			}
			stats, ok := res.(*Statistics)
			if !ok {
				t.Fatalf("Expected *Statistics but got %T", res)
			}
			if len(stats.Result) != 0 {
				t.Errorf("Expected an empty result, got %+v", stats.Result)
			}
		})
	}
}

func TestDecodeStatisticsResponseErrorEnvelope(t *testing.T) {
	// The API answers validation failures with HTTP 200 and error: true. The
	// violation list carries the message; a global violation has no path.
	body := `{"error": true, "violationList": [{"propertyPath": "query.fqdn", "message": "statisticQueryTypeAndFqdn.object.denied"}, {"propertyPath": "", "message": "Wrong format for provided dates!"}], "data": [{}]}`
	resp := &http.Response{Body: io.NopCloser(strings.NewReader(body))}

	res, err := decodeStatisticsResponse(resp, methods["queryStatistics"])
	if err == nil {
		t.Fatalf("Expected an error for an error envelope, got %+v", res)
	}
	if !strings.Contains(err.Error(), "query.fqdn: statisticQueryTypeAndFqdn.object.denied") {
		t.Errorf("Expected the path-scoped violation in the error, got [%s]", err.Error())
	}
	if !strings.Contains(err.Error(), "\nWrong format for provided dates!") || strings.Contains(err.Error(), ": Wrong format") {
		t.Errorf("Expected the global violation without a colon prefix, got [%s]", err.Error())
	}
}

// statisticsTestServer serves POST /statistic/query with the given body and
// records the request payload the client sent.
func statisticsTestServer(t *testing.T, status int, body string) (*API, *json.RawMessage) {
	t.Helper()
	var sent json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/statistic/query" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		payload, _ := io.ReadAll(r.Body)
		sent = append(json.RawMessage(nil), payload...)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	api, err := NewWithToken("token")
	if err != nil {
		t.Fatalf("Unexpected error creating the API client: %v", err)
	}
	api.BaseURL = server.URL + "/%s"
	return api, &sent
}

func testStatisticQuery() *StatisticQuery {
	cest := time.FixedZone("CEST", 2*60*60)
	return &StatisticQuery{
		StartDate:   &types.DateTime{Time: time.Date(2026, 8, 1, 0, 0, 0, 0, cest)},
		EndDate:     &types.DateTime{Time: time.Date(2026, 8, 8, 0, 0, 0, 0, cest)},
		FQDN:        []string{"www.example.com"},
		DataSources: map[string]map[string]string{"requests": {"source": "requests", "type": "stats"}},
	}
}

func TestQueryStatisticsRoundTrip(t *testing.T) {
	api, sent := statisticsTestServer(t, http.StatusOK,
		`{"query": {"type": "fqdn"}, "result": {"requests": {"min": 1, "max": 9, "avg": 5, "sum": 42}}}`)

	stats, err := api.QueryStatistics(testStatisticQuery())
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}
	requests, ok := stats.Result["requests"].(map[string]any)
	if !ok || requests["sum"] != float64(42) {
		t.Errorf("Expected result requests.sum 42, got %+v", stats.Result)
	}

	// The payload is wrapped in {"query": ...}; unset scope and interval are
	// omitted so the API defaults ('fqdn', 'day') apply instead of "" being
	// rejected by the API's choice validation.
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(*sent, &payload); err != nil {
		t.Fatalf("Request payload is not JSON: %v", err)
	}
	var query map[string]any
	if err := json.Unmarshal(payload["query"], &query); err != nil {
		t.Fatalf("Request payload has no query object: %v", err)
	}
	if _, present := query["type"]; present {
		t.Errorf("Expected an unset type to be omitted from the payload, got %s", payload["query"])
	}
	if _, present := query["aggregationInterval"]; present {
		t.Errorf("Expected an unset aggregationInterval to be omitted from the payload, got %s", payload["query"])
	}
	fqdn, _ := query["fqdn"].([]any)
	if len(fqdn) != 1 || fqdn[0] != "www.example.com" || query["startDate"] != "2026-08-01T00:00:00+0200" {
		t.Errorf("Unexpected request payload: %s", payload["query"])
	}
}

func TestQueryStatisticsSurfacesErrorEnvelope(t *testing.T) {
	api, _ := statisticsTestServer(t, http.StatusOK,
		`{"error": true, "violationList": [{"propertyPath": "", "message": "Wrong format for provided dates!"}], "data": [{"query": {}}]}`)

	_, err := api.QueryStatistics(testStatisticQuery())
	if err == nil {
		t.Fatal("Expected an error when the API reports error: true")
	}
	if !strings.Contains(err.Error(), "Wrong format for provided dates!") {
		t.Errorf("Expected the API violation message, got [%s]", err.Error())
	}
}
