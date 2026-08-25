package myrasec

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestGetDNSRecord(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/dns-records/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "www.example.com.", "value": "127.0.0.1", "ttl": 300, "recordType": "A", "upstreamOptions": {"id": 1, "backup": false, "down": false, "failTimeout": "1", "maxFails": 100, "weight": 1}}
			]}`,
			methods["getDNSRecord"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	rec, err := api.GetDNSRecord(1, 1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if rec.ID != 1 {
		t.Errorf("Expected to get DNS record with ID [%d] but got [%d]", 1, rec.ID)
	}

	if rec.Name != "www.example.com." {
		t.Errorf("Expected to get DNS record with Name [%s] but got [%s]", "www.example.com.", rec.Name)
	}

	if rec.Value != "127.0.0.1" {
		t.Errorf("Expected to get DNS record with Value [%s] but got [%s]", "127.0.0.1", rec.Value)
	}

	if rec.RecordType != "A" {
		t.Errorf("Expected to get DNS record with RecordType [%s] but got [%s]", "A", rec.RecordType)
	}

	if rec.TTL != 300 {
		t.Errorf("Expected to get DNS record with TTL [%d] but got [%d]", 300, rec.TTL)
	}

	if rec.UpstreamOptions.ID != 1 {
		t.Errorf("Expected to get DNS record with Upstream-Options ID [%d] but got [%d]", 1, rec.UpstreamOptions.ID)
	}

	if rec.UpstreamOptions.Backup != false {
		t.Errorf("Expected to get DNS record with Upstream-Options Backup [%t] but got [%t]", false, rec.UpstreamOptions.Backup)
	}

	if rec.UpstreamOptions.Down != false {
		t.Errorf("Expected to get DNS record with Upstream-Options Down [%t] but got [%t]", false, rec.UpstreamOptions.Down)
	}

	if rec.UpstreamOptions.FailTimeout != "1" {
		t.Errorf("Expected to get DNS record with Upstream-Options FailTimeout [%s] but got [%s]", "1", rec.UpstreamOptions.FailTimeout)
	}

	if rec.UpstreamOptions.MaxFails != 100 {
		t.Errorf("Expected to get DNS record with Upstream-Options MaxFails [%d] but got [%d]", 100, rec.UpstreamOptions.MaxFails)
	}

	if rec.UpstreamOptions.Weight != 1 {
		t.Errorf("Expected to get DNS record with Upstream-Options Weight [%d] but got [%d]", 1, rec.UpstreamOptions.Weight)
	}
}

func TestListDNSRecords(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/dns-records",
			`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [
				{"id": 1, "name": "www.example.com.", "value": "127.0.0.1", "ttl": 300, "recordType": "A", "upstreamOptions": {"id": 1, "backup": false, "down": false, "failTimeout": "1", "maxFails": 100, "weight": 1}}, 
				{"id": 2, "name": "example.com.", "value": "127.0.0.1", "ttl": 300, "recordType": "A", "upstreamOptions": {"id": 2, "backup": false, "down": false, "failTimeout": "1", "maxFails": 100, "weight": 1}}
			]}`,
			methods["listDNSRecords"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	records, err := api.ListDNSRecords(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(records) != 2 {
		t.Errorf("Expected to get [%d] DNS records but got [%d]", 2, len(records))
	}
}

// TestCreateDNSRecordWithEmptyEndpoints covers SB-2970: the API answers a create request
// with "endpoints": [] (an empty JSON array) instead of an empty JSON object, which made
// every CreateDNSRecord call fail while decoding the response.
func TestCreateDNSRecordWithEmptyEndpoints(t *testing.T) {
	cache, err := preCacheRequestWithError(
		"https://apiv2.myracloud.com/domain/1/dns-records",
		`{"error": false, "violationList": [], "warningList": [], "data": [
			{"id": 1, "name": "www.example.com.", "value": "127.0.0.1", "ttl": 300, "recordType": "A", "endpoints": []}
		]}`,
		methods["createDNSRecord"],
	)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	api, err := setupPreCachedAPI([]*TestCache{cache})
	if err != nil {
		t.Error("Unexpected error.")
	}

	rec, err := api.CreateDNSRecord(&DNSRecord{Name: "www.example.com.", Value: "127.0.0.1", TTL: 300, RecordType: "A"}, 1)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if rec.ID != 1 {
		t.Errorf("Expected to get DNS record with ID [%d] but got [%d]", 1, rec.ID)
	}

	if len(rec.Endpoints) != 0 {
		t.Errorf("Expected to get DNS record with [%d] endpoints but got [%d]", 0, len(rec.Endpoints))
	}
}

// TestListDNSRecordsWithEndpoints covers SB-2970: a list response mixes records that carry
// endpoints as a JSON object with records that carry them as an empty JSON array.
func TestListDNSRecordsWithEndpoints(t *testing.T) {
	cache, err := preCacheRequestWithError(
		"https://apiv2.myracloud.com/domain/1/dns-records",
		`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [
			{"id": 1, "name": "www.example.com.", "value": "127.0.0.1", "ttl": 300, "recordType": "A", "endpoints": {"ipv4": ["127.0.0.1", "127.0.0.2"], "ipv6": ["::1"]}},
			{"id": 2, "name": "example.com.", "value": "some text", "ttl": 300, "recordType": "TXT", "endpoints": []}
		]}`,
		methods["listDNSRecords"],
	)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	api, err := setupPreCachedAPI([]*TestCache{cache})
	if err != nil {
		t.Error("Unexpected error.")
	}

	records, err := api.ListDNSRecords(1, nil)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(records) != 2 {
		t.Fatalf("Expected to get [%d] DNS records but got [%d]", 2, len(records))
	}

	if len(records[0].Endpoints["ipv4"]) != 2 {
		t.Errorf("Expected to get [%d] ipv4 endpoints but got [%d]", 2, len(records[0].Endpoints["ipv4"]))
	}

	if records[0].Endpoints["ipv4"][0] != "127.0.0.1" {
		t.Errorf("Expected first ipv4 endpoint to be [%s] but got [%s]", "127.0.0.1", records[0].Endpoints["ipv4"][0])
	}

	if len(records[0].Endpoints["ipv6"]) != 1 || records[0].Endpoints["ipv6"][0] != "::1" {
		t.Errorf("Expected to get ipv6 endpoint [%s] but got %v", "::1", records[0].Endpoints["ipv6"])
	}

	if records[1].Endpoints != nil {
		t.Errorf("Expected TXT record endpoints to be nil but got %v", records[1].Endpoints)
	}
}

func TestDNSRecordEndpointsUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		expectErr bool
		expected  DNSRecordEndpoints
	}{
		{
			name:     "object",
			json:     `{"endpoints": {"ipv4": ["127.0.0.1"], "ipv6": ["::1"]}}`,
			expected: DNSRecordEndpoints{"ipv4": {"127.0.0.1"}, "ipv6": {"::1"}},
		},
		{
			name:     "empty object",
			json:     `{"endpoints": {}}`,
			expected: nil,
		},
		{
			name:     "empty array",
			json:     `{"endpoints": []}`,
			expected: nil,
		},
		{
			name:     "empty array with inner whitespace",
			json:     `{"endpoints": [ ]}`,
			expected: nil,
		},
		{
			name:     "null",
			json:     `{"endpoints": null}`,
			expected: nil,
		},
		{
			name:     "attribute missing",
			json:     `{"id": 1}`,
			expected: nil,
		},
		{
			name:      "non-empty array",
			json:      `{"endpoints": ["127.0.0.1"]}`,
			expectErr: true,
		},
		{
			name:      "wrong value type",
			json:      `{"endpoints": {"ipv4": "127.0.0.1"}}`,
			expectErr: true,
		},
		{
			name:      "scalar",
			json:      `{"endpoints": 1}`,
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var rec DNSRecord
			err := json.Unmarshal([]byte(test.json), &rec)

			if test.expectErr {
				if err == nil {
					t.Fatalf("Expected to get an error for [%s] but got none", test.json)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected not to get an error but got [%s]", err.Error())
			}

			if !reflect.DeepEqual(rec.Endpoints, test.expected) {
				t.Errorf("Expected to get endpoints %#v but got %#v", test.expected, rec.Endpoints)
			}
		})
	}
}

// TestDNSRecordEndpointsMarshalJSON makes sure the read-only endpoints attribute is not
// sent back to the API when it is empty.
func TestDNSRecordEndpointsMarshalJSON(t *testing.T) {
	payload, err := json.Marshal(&DNSRecord{Name: "www.example.com.", Value: "127.0.0.1", RecordType: "A"})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if strings.Contains(string(payload), "endpoints") {
		t.Errorf("Expected payload not to contain the endpoints attribute but got [%s]", string(payload))
	}

	payload, err = json.Marshal(&DNSRecord{Name: "www.example.com.", Value: "127.0.0.1", RecordType: "A", Endpoints: DNSRecordEndpoints{"ipv4": {"127.0.0.1"}}})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if !strings.Contains(string(payload), `"endpoints":{"ipv4":["127.0.0.1"]}`) {
		t.Errorf("Expected payload to contain the endpoints attribute but got [%s]", string(payload))
	}
}
