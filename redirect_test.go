package myrasec

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestGetRedirect(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/redirects/www.example.com/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "source": "/from", "destination": "/to", "type": "redirect", "subDomainName": "www.example.com", "matchingType": "exact", "sort": 0, "expertMode": false, "enabled": true}
			]}`,
			"getRedirect",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	redirect, err := api.GetRedirect(1, "www.example.com", 1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if redirect.ID != 1 {
		t.Errorf("Expected to get Redirect with ID [%d] but got [%d]", 1, redirect.ID)
	}

	if redirect.Source != "/from" {
		t.Errorf("Expected to get Redirect with Source [%s] but got [%s]", "/from", redirect.Source)
	}

	if redirect.Destination != "/to" {
		t.Errorf("Expected to get Redirect with Destination [%s] but got [%s]", "/to", redirect.Destination)
	}

	if redirect.Type != "redirect" {
		t.Errorf("Expected to get Redirect with Type [%s] but got [%s]", "redirect", redirect.Type)
	}

	if redirect.MatchingType != "exact" {
		t.Errorf("Expected to get Redirect with MatchingType [%s] but got [%s]", "exact", redirect.MatchingType)
	}

	if redirect.SubDomainName != "www.example.com" {
		t.Errorf("Expected to get Redirect with SubDomainName [%s] but got [%s]", "www.example.com", redirect.SubDomainName)
	}

	if redirect.Enabled != true {
		t.Errorf("Expected to get Redirect with Enabled [%t] but got [%t]", true, redirect.Enabled)
	}

	if redirect.ExpertMode != false {
		t.Errorf("Expected to get Redirect with ExperMode [%t] but got [%t]", false, redirect.ExpertMode)
	}

	if redirect.Sort != 0 {
		t.Errorf("Expected to get Redirect with sort [%d] but got [%d]", 0, redirect.Sort)
	}
}

func TestListRedirects(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/redirects/www.example.com",
			`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [
				{"id": 1, "source": "/from", "destination": "/to", "type": "redirect", "subDomainName": "www.example.com", "matchingType": "exact", "sort": 0, "expertMode": false, "enabled": true},
		 		{"id": 2, "source": "/index.html", "destination": "/index.php", "type": "permanent", "subDomainName": "www.example.com", "matchingType": "prefix", "sort": 1, "expertMode": false, "enabled": false}
			]}`,
			"listRedirects",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	redirects, err := api.ListRedirects(1, "www.example.com", nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(redirects) != 2 {
		t.Errorf("Expected to get [%d] error pages but got [%d]", 2, len(redirects))
	}

	for _, r := range redirects {
		if !intInSlice(r.ID, []int{1, 2}) {
			t.Errorf("Unexpected redirect ID [%d]", r.ID)
		}

		if r.ID == 1 {
			if r.Source != "/from" {
				t.Errorf("Expected to get redirect with Source [%s] but got [%s]", "/from", r.Source)
			}

			if r.Destination != "/to" {
				t.Errorf("Expected to get redirect with Destination [%s] but got [%s]", "/to", r.Destination)
			}

			if r.MatchingType != "exact" {
				t.Errorf("Expected to get redirect with MatchingType [%s] but got [%s]", "exact", r.MatchingType)
			}

			if r.Type != "redirect" {
				t.Errorf("Expected to get redirect with Type [%s] but got [%s]", "redirect", r.Type)
			}

			if !r.Enabled {
				t.Errorf("Expected to get an enabled redirect")
			}

			if r.Sort != 0 {
				t.Errorf("Expected to get redirect with Sort [%d] but got [%d]", 0, r.Sort)
			}
		}

		if r.ID == 2 {
			if r.Source != "/index.html" {
				t.Errorf("Expected to get redirect with Source [%s] but got [%s]", "/index.html", r.Source)
			}

			if r.Destination != "/index.php" {
				t.Errorf("Expected to get redirect with Destination [%s] but got [%s]", "/index.php", r.Destination)
			}

			if r.MatchingType != "prefix" {
				t.Errorf("Expected to get redirect with MatchingType [%s] but got [%s]", "prefix", r.MatchingType)
			}

			if r.Type != "permanent" {
				t.Errorf("Expected to get redirect with Type [%s] but got [%s]", "permanent", r.Type)
			}

			if r.Enabled {
				t.Errorf("Expected to get a disabled redirect")
			}

			if r.Sort != 1 {
				t.Errorf("Expected to get redirect with Sort [%d] but got [%d]", 0, r.Sort)
			}
		}

	}
}

func TestRedirectSendsDisabledExpertMode(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /domain/1/redirects/www.example.com": {Status: http.StatusOK, Body: `{"error": false, "targetObject": [{"id": 1, "type": "permanent", "subDomainName": "www.example.com", "source": "/old", "destination": "/new", "matchingType": "exact", "enabled": true, "expertMode": false}]}`},
	})

	_, err := api.CreateRedirectContext(context.Background(), &Redirect{Type: "permanent", Source: "/old", Destination: "/new", MatchingType: "exact", Enabled: true}, 1, "www.example.com")
	if err != nil {
		t.Fatalf("Expected not to get an error but got %v", err)
	}

	body := string(requests.last(t).Body)
	for _, want := range []string{`"enabled":true`, `"expertMode":false`} {
		if !strings.Contains(body, want) {
			t.Errorf("Expected the request body to contain %s but got %s", want, body)
		}
	}
}
