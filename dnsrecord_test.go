package myrasec

import (
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
