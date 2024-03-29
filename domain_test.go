package myrasec

import (
	"testing"
)

func TestGetDomain(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domains/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "example.com"}
			]}`,
			methods["getDomain"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	domain, err := api.GetDomain(1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if domain.ID != 1 {
		t.Errorf("Expected to get Domain with ID [%d] but got [%d]", 1, domain.ID)
	}

	if domain.Name != "example.com" {
		t.Errorf("Expected to get Domain with Name [%s] but got [%s]", "example.com", domain.Name)
	}
}

func TestListDomains(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domains",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"id": 1, "name": "example.com"}, 
				{"id": 2, "name": "example2.com"}, 
				{"id": 3, "name": "example3.com"}
			]}`,
			methods["listDomains"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	domains, err := api.ListDomains(nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(domains) != 3 {
		t.Errorf("Expected to get [%d] domains but got [%d]", 3, len(domains))
	}
}
