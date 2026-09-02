package myrasec

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetErrorPage(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/errorpages/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "errorCode": 500, "content": "<!DOCTYPE html><html><head><title>Error 500</title></head><body><h1>HTTP 500 error</h1></body></html>", "subDomainName": "www.example.com"}
			]}`,
			"getErrorPage",
		),
	)
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
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/errorpages",
			`{"error": false, "pageSize": 10, "page": 1, "count": 2, "data": [
				{"id": 1, "errorCode": 500, "content": "<!DOCTYPE html><html><head><title>Error 500</title></head><body><h1>HTTP 500 error</h1></body></html>", "subDomainName": "www.example.com"}, 
				{"id": 2, "errorCode": 404, "content": "<!DOCTYPE html><html><head><title>Error 404</title></head><body><h1>HTTP 404 error</h1></body></html>", "subDomainName": "test.example.com"}
			]}`,
			"listErrorPages",
		),
	)
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

// TestCreateErrorPage verifies that the create response body is read back so the
// returned page carries the server-assigned id, created and modified instead of the
// zero values.
func TestCreateErrorPage(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/errorpages",
			`{"error": false, "data": [
				{"id": 7, "errorCode": 404, "content": "<h1>HTTP 404 error</h1>", "subDomainName": "www.example.com", "created": "2025-01-09T16:31:13+0100", "modified": "2025-04-02T10:15:49+0200"}
			]}`,
			"createErrorPage",
		),
	)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	page, err := api.CreateErrorPage(&ErrorPage{
		ErrorCode:     404,
		Content:       "<h1>HTTP 404 error</h1>",
		SubDomainName: "www.example.com",
	}, 1)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if page.ID != 7 {
		t.Errorf("Expected to get ErrorPage with ID [%d] but got [%d]", 7, page.ID)
	}

	if page.Created == nil {
		t.Errorf("Expected the server-assigned Created timestamp to be copied onto the returned page but got nil")
	}

	if page.Modified == nil {
		t.Errorf("Expected the server-assigned Modified timestamp to be copied onto the returned page but got nil")
	}

	if page.ErrorCode != 404 {
		t.Errorf("Expected to get ErrorPage for ErrorCode [%d] but got [%d]", 404, page.ErrorCode)
	}

	if page.SubDomainName != "www.example.com" {
		t.Errorf("Expected to get ErrorPage for SubDomainName [%s] but got [%s]", "www.example.com", page.SubDomainName)
	}
}

// TestUpdateErrorPage mirrors TestCreateErrorPage for the update endpoint.
func TestUpdateErrorPage(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/errorpages",
			`{"error": false, "data": [
				{"id": 9, "errorCode": 500, "content": "<h1>HTTP 500 error</h1>", "subDomainName": "www.example.com", "created": "2025-01-09T16:31:13+0100", "modified": "2025-04-02T10:15:49+0200"}
			]}`,
			"updateErrorPage",
		),
	)
	if err != nil {
		t.Fatal("Unexpected error")
	}

	page, err := api.UpdateErrorPage(&ErrorPage{
		ErrorCode:     500,
		Content:       "<h1>HTTP 500 error</h1>",
		SubDomainName: "www.example.com",
	}, 1)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if page.ID != 9 {
		t.Errorf("Expected to get ErrorPage with ID [%d] but got [%d]", 9, page.ID)
	}

	if page.Modified == nil {
		t.Errorf("Expected the server-assigned Modified timestamp to be copied onto the returned page but got nil")
	}
}

// TestDecodeErrorPageResponse exercises the decoder directly, including the
// backward-compatibility fallback for a server without the companion change, which
// answers with an empty data array.
func TestDecodeErrorPageResponse(t *testing.T) {
	definition := APIMethod{Result: ErrorPage{}}

	// Populated data envelope: the first element is decoded into an *ErrorPage.
	populated := http.Response{
		Body: io.NopCloser(bytes.NewBufferString(`{"error": false, "data": [{"id": 7, "errorCode": 404, "subDomainName": "www.example.com"}]}`)),
	}
	res, err := decodeErrorPageResponse(&populated, definition)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}
	page, ok := res.(*ErrorPage)
	if !ok || page == nil {
		t.Fatalf("Expected a non-nil *ErrorPage but got [%T]", res)
	}
	if page.ID != 7 {
		t.Errorf("Expected decoded ErrorPage ID [%d] but got [%d]", 7, page.ID)
	}

	// Empty data array (server without the companion change): the decoder returns nil
	// so the caller falls back to the object it sent.
	empty := http.Response{
		Body: io.NopCloser(bytes.NewBufferString(`{"error": false, "data": []}`)),
	}
	res, err = decodeErrorPageResponse(&empty, definition)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}
	if res != nil {
		t.Errorf("Expected a nil result for an empty data envelope but got [%T]", res)
	}
}
