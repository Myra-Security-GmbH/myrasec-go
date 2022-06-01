package myrasec

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func preCacheGetDomain(api *API) {
	responseJSON := `{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [{"id": 1, "name": "example.com"}]}`

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains/1", nil)
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(responseJSON)),
	}

	res, _ := methods["getDomain"].ResponseDecodeFunc(&resp, methods["getDomain"])
	api.cacheResponse(req, res)
}

func preCacheListDomains(api *API) {

	responseJSON := `{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [{"id": 1, "name": "example.com"}, {"id": 2, "name": "example2.com"}, {"id": 3, "name": "example3.com"}]}`

	req, _ := http.NewRequest(http.MethodGet, "https://apiv2.myracloud.com/domains", nil)
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(responseJSON)),
	}

	res, _ := decodeDefaultResponse(&resp, methods["listDomains"])
	api.cacheResponse(req, res)
}

//
// preCacheDomainAPI will mock the data, returned by the api.call function. Like this we can test without sending real API requests.
//
func preCacheDomainAPI() (*API, error) {
	api, _ := New("abc123", "123abc")
	api.EnableCaching()
	api.SetCachingTTL(120)

	preCacheGetDomain(api)
	preCacheListDomains(api)

	return api, nil
}

func TestGetDomain(t *testing.T) {

	api, err := preCacheDomainAPI()
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
	api, err := preCacheDomainAPI()
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
