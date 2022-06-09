package myrasec

import (
	"strings"
	"testing"
)

func TestGetErrorPage(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/errorpages/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "errorCode": 500, "content": "<!DOCTYPE html><html><head><title>Error 500</title></head><body><h1>HTTP 500 error</h1></body></html>", "subDomainName": "www.example.com"}
			]}`,
			methods["getErrorPage"],
		),
	})
	if err != nil {
		t.Error("Unexpected error")
	}

	page, err := api.GetErrorPage(1, 1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if page.ID != 1 {
		t.Errorf("Expected to get ErrorPage with ID [%d] but got [%d]", 1, page.ID)
	}

	if page.ErrorCode != 500 {
		t.Errorf("Expected to get ErrorPage for ErrorCode [%d] but got [%d]", 500, page.ErrorCode)
	}

	if !strings.Contains(page.Content, "<h1>HTTP 500 error</h1>") {
		t.Errorf("Expected to have [\"%s\"] in the ErrorPage content but did not find it.", "<h1>HTTP 500 error</h1>")
	}
}

func TestListErrorPages(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/errorpages",
			`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [
				{"id": 1, "errorCode": 500, "content": "<!DOCTYPE html><html><head><title>Error 500</title></head><body><h1>HTTP 500 error</h1></body></html>", "subDomainName": "www.example.com"}, 
				{"id": 2, "errorCode": 404, "content": "<!DOCTYPE html><html><head><title>Error 404</title></head><body><h1>HTTP 404 error</h1></body></html>", "subDomainName": "test.example.com"}
			]}`,
			methods["listErrorPages"],
		),
	})
	if err != nil {
		t.Error("Unexpected error")
	}

	pages, err := api.ListErrorPages(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(pages) != 2 {
		t.Errorf("Expected to get [%d] error pages but got [%d]", 2, len(pages))
	}
}
