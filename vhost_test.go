package myrasec

import (
	"testing"
)

func TestListAllSubdomains(t *testing.T) {

	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/subdomains",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"id": 1, "label": "www.example.com", "value": "www.example.com.", "domainName": "example.com", "access": true, "paused": false},
				{"id": 2, "label": "example.com", "value": "example.com.", "domainName": "example.com", "access": true, "paused": false},
				{"id": 3, "label": "www.myrasecurity.com", "value": "www.myrasecurity.com.", "domainName": "myrasecurity.com", "access": true, "paused": false}
			]}`,
			methods["listAllSubdomains"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	subdomains, err := api.ListAllSubdomains(nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(subdomains) != 3 {
		t.Errorf("Expected to get [%d] subdomains but got [%d]", 3, len(subdomains))
	}

	for _, s := range subdomains {
		if !intInSlice(s.ID, []int{1, 2, 3}) {
			t.Errorf("Unexpected IP range ID [%d]", s.ID)
		}

		if s.Access != true {
			t.Errorf("Expected to get Subdomain with Access [%t] but got [%t]", true, s.Access)
		}

		if s.Paused != false {
			t.Errorf("Expected to get Subdomain with Paused [%t] but got [%t]", false, s.Paused)
		}

		if s.ID == 1 {
			if s.Label != "www.example.com" {
				t.Errorf("Expected to get Subdomain with Label [%s] but got [%s]", "www.example.com", s.Label)
			}

			if s.Value != "www.example.com." {
				t.Errorf("Expected to get Subdomain with Value [%s] but got [%s]", "www.example.com.", s.Value)
			}

			if s.DomainName != "example.com" {
				t.Errorf("Expected to get Subdomain with DomainName [%s] but got [%s]", "example.com", s.DomainName)
			}
		}

		if s.ID == 2 {
			if s.Label != "example.com" {
				t.Errorf("Expected to get Subdomain with Label [%s] but got [%s]", "example.com", s.Label)
			}

			if s.Value != "example.com." {
				t.Errorf("Expected to get Subdomain with Value [%s] but got [%s]", "example.com.", s.Value)
			}

			if s.DomainName != "example.com" {
				t.Errorf("Expected to get Subdomain with DomainName [%s] but got [%s]", "example.com", s.DomainName)
			}
		}

		if s.ID == 3 {
			if s.Label != "www.myrasecurity.com" {
				t.Errorf("Expected to get Subdomain with Label [%s] but got [%s]", "www.myrasecurity.com", s.Label)
			}

			if s.Value != "www.myrasecurity.com." {
				t.Errorf("Expected to get Subdomain with Value [%s] but got [%s]", "www.myrasecurity.com.", s.Value)
			}

			if s.DomainName != "myrasecurity.com" {
				t.Errorf("Expected to get Subdomain with DomainName [%s] but got [%s]", "myrasecurity.com", s.DomainName)
			}
		}
	}

}

func TestListAllSubdomainsForDomain(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/subdomains",
			`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [
				{"id": 1, "label": "www.example.com", "value": "www.example.com.", "domainName": "example.com", "access": true, "paused": false},
				{"id": 2, "label": "example.com", "value": "example.com.", "domainName": "example.com", "access": true, "paused": false}
			]}`,
			methods["listSubdomainsForDomain"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	subdomains, err := api.ListAllSubdomainsForDomain(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(subdomains) != 2 {
		t.Errorf("Expected to get [%d] subdomains but got [%d]", 2, len(subdomains))
	}

	for _, s := range subdomains {
		if s.DomainName != "example.com" {
			t.Errorf("Expected to get Subdomain with DomainName [%s] but got [%s]", "example.com", s.DomainName)
		}

		if !intInSlice(s.ID, []int{1, 2}) {
			t.Errorf("Unexpected IP range ID [%d]", s.ID)
		}

		if s.Access != true {
			t.Errorf("Expected to get Subdomain with Access [%t] but got [%t]", true, s.Access)
		}

		if s.Paused != false {
			t.Errorf("Expected to get Subdomain with Paused [%t] but got [%t]", false, s.Paused)
		}

		if s.ID == 1 {
			if s.Label != "www.example.com" {
				t.Errorf("Expected to get Subdomain with Label [%s] but got [%s]", "www.example.com", s.Label)
			}

			if s.Value != "www.example.com." {
				t.Errorf("Expected to get Subdomain with Value [%s] but got [%s]", "www.example.com.", s.Value)
			}
		}

		if s.ID == 2 {
			if s.Label != "example.com" {
				t.Errorf("Expected to get Subdomain with Label [%s] but got [%s]", "example.com", s.Label)
			}

			if s.Value != "example.com." {
				t.Errorf("Expected to get Subdomain with Value [%s] but got [%s]", "example.com.", s.Value)
			}
		}
	}

}
