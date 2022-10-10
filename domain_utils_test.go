package myrasec

import (
	"testing"
)

func TestFetchDomainForSubdomainName(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domains/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "example.com"}
			]}`,
			methods["getDomain"],
		),
		preCacheRequest(
			"https://apiv2.myracloud.com/domains?search=example.com",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "example.com"}
			]}`,
			methods["listDomains"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	domain, err := api.FetchDomainForSubdomainName("ALL-1")
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if domain.ID != 1 {
		t.Errorf("Expected to get domain with ID [%d] but got ID [%d]", 1, domain.ID)
	}

	if domain.Name != "example.com" {
		t.Errorf("Expected to get domain with name [%s] but got name [%s]", "example.com", domain.Name)
	}

	domain, err = api.FetchDomainForSubdomainName("ALL:example.com")
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if domain.ID != 1 {
		t.Errorf("Expected to get domain with ID [%d] but got ID [%d]", 1, domain.ID)
	}
}

func TestFetchDomain(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domains?search=example.com",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "example.com"}
			]}`,
			methods["listDomains"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	domain, err := api.FetchDomain("example.com")
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if domain.ID != 1 {
		t.Errorf("Expected to get domain with ID [%d] but got ID [%d]", 1, domain.ID)
	}

	if domain.Name != "example.com" {
		t.Errorf("Expected to get domain with name [%s] but got name [%s]", "example.com", domain.Name)
	}
}

func TestEnsureTrailingDot(t *testing.T) {
	var name string

	name = EnsureTrailingDot("www.example.com")
	if name != "www.example.com." {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "www.example.com.", "www.example.com", name)
	}

	name = EnsureTrailingDot("www.example.com.")
	if name != "www.example.com." {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "www.example.com.", "www.example.com.", name)
	}

	name = EnsureTrailingDot("www.example.com......")
	if name != "www.example.com." {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "www.example.com.", "www.example.com......", name)
	}

	name = EnsureTrailingDot(".www.example.com..")
	if name != ".www.example.com." {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", ".www.example.com.", ".www.example.com..", name)
	}

	name = EnsureTrailingDot("")
	if name != "." {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", ".", ".", name)
	}
	name = EnsureTrailingDot(".")
	if name != "." {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", ".", ".", name)
	}
	name = EnsureTrailingDot("..")
	if name != "." {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", ".", "..", name)
	}
}

func TestRemoveTrailingDot(t *testing.T) {
	var name string
	name = RemoveTrailingDot("www.example.com")
	if name != "www.example.com" {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "www.example.com", "www.example.com", name)
	}

	name = RemoveTrailingDot("www.example.com.")
	if name != "www.example.com" {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "www.example.com", "www.example.com.", name)
	}

	name = RemoveTrailingDot("www.example.com.....")
	if name != "www.example.com" {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "www.example.com", "www.example.com.....", name)
	}

	name = RemoveTrailingDot(".www.example.com..")
	if name != ".www.example.com" {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", ".www.example.com", ".www.example.com..", name)
	}

	name = RemoveTrailingDot("")
	if name != "" {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "", "", name)
	}

	name = RemoveTrailingDot(".")
	if name != "" {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "", ".", name)
	}

	name = RemoveTrailingDot("..")
	if name != "" {
		t.Errorf("Expected to get [%s] for [%s] but got [%s]", "", "..", name)
	}

}

func TestIsGeneralDomainName(t *testing.T) {
	var res bool

	res = IsGeneralDomainName("example.com")
	if res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", false, "example.com", res)
	}

	res = IsGeneralDomainName("ALL:example.com")
	if !res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", true, "ALL:example.com", res)
	}

	res = IsGeneralDomainName("ALL-example.com")
	if !res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", true, "ALL-example.com", res)
	}

	res = IsGeneralDomainName("ALL:1234")
	if !res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", true, "ALL:1234", res)
	}

	res = IsGeneralDomainName("ALL-1234")
	if !res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", true, "ALL-1234", res)
	}

	res = IsGeneralDomainName("ALL|example.com")
	if res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", false, "ALL|example.com", res)
	}

	res = IsGeneralDomainName("")
	if res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", false, "", res)
	}

	res = IsGeneralDomainName("    ")
	if res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", false, "    ", res)
	}

	res = IsGeneralDomainName("-:")
	if res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", false, "-:", res)
	}

	res = IsGeneralDomainName("ALL")
	if res {
		t.Errorf("Expected to get [%t] for [%s] but got [%t]", false, "ALL", res)
	}
}
