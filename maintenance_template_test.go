package myrasec

import (
	"strings"
	"testing"
)

func TestListMaintenanceTemplates(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/maintenance-templates",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "Tesing Maintenance Templates", "content":"<!DOCTYPE html><html><head><title>Maintenance</title></head><body><h1>Maintenance</h1></body></html>"}
			]}`,
			methods["listMaintenanceTemplates"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	templates, err := api.ListMaintenanceTemplates(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(templates) != 1 {
		t.Errorf("Expected to get [%d] error pages but got [%d]", 1, len(templates))
	}

	for _, mt := range templates {
		if mt.ID != 1 {
			t.Errorf("Expected to get Maintenance with ID [%d] but got [%d]", 1, mt.ID)
		}

		if !strings.Contains(mt.Content, "<h1>Maintenance</h1>") {
			t.Errorf("Expected to have [\"%s\"] in the Maintenance content but did not find it.", "<h1>Maintenance</h1>")
		}
	}
}
