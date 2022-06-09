package myrasec

import (
	"strings"
	"testing"
)

func TestListMaintenances(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/www.example.com/maintenances",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "fqdn": "www.example.com", "start": "2022-06-10T00:00:00+0200", "end": "2022-06-10T06:00:00+0200", "content":"<!DOCTYPE html><html><head><title>Maintenance</title></head><body><h1>Maintenance</h1></body></html>"}
			]}`,
			methods["listMaintenances"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	maintenances, err := api.ListMaintenances(1, "www.example.com", nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(maintenances) != 1 {
		t.Errorf("Expected to get [%d] error pages but got [%d]", 1, len(maintenances))
	}

	for _, m := range maintenances {
		if m.ID != 1 {
			t.Errorf("Expected to get Maintenance with ID [%d] but got [%d]", 1, m.ID)
		}

		if m.Start.Format("2006-01-02 15:04:05") != "2022-06-10 00:00:00" {
			t.Errorf("Expected to have Start date [%s] but got [%s]", "2022-06-10 00:00:00", m.End.Format("2006-01-02 15:04:05"))
		}

		if m.End.Format("2006-01-02 15:04:05") != "2022-06-10 06:00:00" {
			t.Errorf("Expected to have End date [%s] but got [%s]", "2022-06-10 06:00:00", m.End.Format("2006-01-02 15:04:05"))
		}

		if !strings.Contains(m.Content, "<h1>Maintenance</h1>") {
			t.Errorf("Expected to have [\"%s\"] in the Maintenance content but did not find it.", "<h1>Maintenance</h1>")
		}
	}
}
