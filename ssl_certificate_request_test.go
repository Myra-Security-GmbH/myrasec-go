package myrasec

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

const sslCertificateRequestJSON = `{
	"objectType": "SslCertRequestVO", "id": 42,
	"created": "2026-08-01T10:00:00+0200", "modified": "2026-08-02T11:30:00+0200",
	"algorithm": "ECDSA256", "provider": "SECTIGO", "status": "FAILED",
	"failureReason": "CNAME_TIMEOUT", "customerActionable": true,
	"subjectAlternativeNames": [
		{"objectType": "SslCertRequestSanVO", "id": 7, "created": "2026-08-01T10:00:00+0200", "modified": "2026-08-01T10:00:00+0200", "name": "www.example.com"},
		{"objectType": "SslCertRequestSanVO", "id": 8, "created": "2026-08-01T10:00:00+0200", "modified": "2026-08-01T10:00:00+0200", "name": "*.example.com"}
	],
	"assignments": [
		{"objectType": "SslCertRequestAssignmentVO", "id": 9, "created": "2026-08-01T10:00:00+0200", "modified": "2026-08-01T10:00:00+0200", "subDomainName": "www.example.com"}
	],
	"multiDomain": false, "sslProviderCredentialsId": 3, "renewalInterval": 30, "signatureAlgorithm": "SHA384"
}`

func assertSSLCertificateRequest(t *testing.T, request *SSLCertificateRequest) {
	t.Helper()

	if request.ID != 42 {
		t.Errorf("Expected ID [%d] but got [%d]", 42, request.ID)
	}

	if request.Created.Format("2006-01-02") != "2026-08-01" {
		t.Errorf("Expected Created [%s] but got [%s]", "2026-08-01", request.Created.Format("2006-01-02"))
	}

	if request.Algorithm != SSLCertificateRequestAlgorithmECDSA256 {
		t.Errorf("Expected Algorithm [%s] but got [%s]", SSLCertificateRequestAlgorithmECDSA256, request.Algorithm)
	}

	if request.Provider != SSLProviderSectigo {
		t.Errorf("Expected Provider [%s] but got [%s]", SSLProviderSectigo, request.Provider)
	}

	if request.Status != SSLCertificateRequestStatusFailed {
		t.Errorf("Expected Status [%s] but got [%s]", SSLCertificateRequestStatusFailed, request.Status)
	}

	if request.FailureReason != SSLCertificateRequestFailureReasonCNAMETimeout {
		t.Errorf("Expected FailureReason [%s] but got [%s]", SSLCertificateRequestFailureReasonCNAMETimeout, request.FailureReason)
	}

	if !request.CustomerActionable {
		t.Error("Expected CustomerActionable to be true")
	}

	if len(request.SubjectAlternativeNames) != 2 {
		t.Fatalf("Expected [%d] subject alternative names but got [%d]", 2, len(request.SubjectAlternativeNames))
	}

	if request.SubjectAlternativeNames[1].Name != "*.example.com" || request.SubjectAlternativeNames[1].ID != 8 {
		t.Errorf("Expected second SAN [%s] with ID [%d] but got [%s] with ID [%d]", "*.example.com", 8, request.SubjectAlternativeNames[1].Name, request.SubjectAlternativeNames[1].ID)
	}

	if len(request.Assignments) != 1 || request.Assignments[0].SubDomainName != "www.example.com" {
		t.Errorf("Expected one assignment for [%s] but got %+v", "www.example.com", request.Assignments)
	}

	if request.MultiDomain {
		t.Error("Expected MultiDomain to be false")
	}

	if request.SSLProviderCredentialsID != 3 {
		t.Errorf("Expected SSLProviderCredentialsID [%d] but got [%d]", 3, request.SSLProviderCredentialsID)
	}

	if request.RenewalInterval != 30 {
		t.Errorf("Expected RenewalInterval [%d] but got [%d]", 30, request.RenewalInterval)
	}

	if request.SignatureAlgorithm != SSLCertificateRequestSignatureAlgorithmSHA384 {
		t.Errorf("Expected SignatureAlgorithm [%s] but got [%s]", SSLCertificateRequestSignatureAlgorithmSHA384, request.SignatureAlgorithm)
	}
}

func TestListSSLCertificateRequests(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /ssl/requests": {Status: http.StatusOK, Body: `{"error": false, "pageSize": 50, "page": 1, "count": 2, "data": [` + sslCertificateRequestJSON + `,
			{"objectType": "SslCertRequestVO", "id": 43, "algorithm": "RSA2048", "provider": "LETS_ENCRYPT", "status": "CREATED", "subjectAlternativeNames": [], "assignments": [], "multiDomain": false}
		]}`},
	})

	list, err := api.ListSSLCertificateRequests(map[string]string{
		"status": "OPEN,CREATED",
		"domain": "example.com",
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(list) != 2 {
		t.Fatalf("Expected [%d] requests but got [%d]", 2, len(list))
	}

	assertSSLCertificateRequest(t, &list[0])

	if list[1].Status != SSLCertificateRequestStatusCreated || list[1].SSLProviderCredentialsID != 0 {
		t.Errorf("Expected second request in status CREATED without credentials, got %+v", list[1])
	}

	sent := requests.last(t)
	if sent.Query.Get("status") != "OPEN,CREATED" || sent.Query.Get("domain") != "example.com" {
		t.Errorf("Expected status and domain query parameters, got [%s]", sent.Query.Encode())
	}
}

func TestGetSSLCertificateRequest(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /ssl/requests/42": {Status: http.StatusOK, Body: `{"error": false, "data": [` + sslCertificateRequestJSON + `]}`},
	})

	request, err := api.GetSSLCertificateRequest(42)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	assertSSLCertificateRequest(t, request)

	if sent := requests.last(t); sent.Method != http.MethodGet || sent.Path != "/ssl/requests/42" {
		t.Errorf("Expected GET /ssl/requests/42 but got %s %s", sent.Method, sent.Path)
	}
}

func TestCreateSSLCertificateRequest(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /ssl/requests": {Status: http.StatusCreated, Body: `{"error": false, "data": [
			{"objectType": "SslCertRequestVO", "id": 44, "created": "2026-08-01T10:00:00+0200", "modified": "2026-08-01T10:00:00+0200", "algorithm": "RSA2048", "provider": "LETS_ENCRYPT", "status": "OPEN",
			 "subjectAlternativeNames": [{"id": 10, "name": "www.example.com"}], "assignments": [{"id": 11, "subDomainName": "www.example.com"}], "multiDomain": false}
		]}`},
	})

	created, err := api.CreateSSLCertificateRequest(&SSLCertificateRequest{
		Provider:                SSLProviderLetsEncrypt,
		Algorithm:               SSLCertificateRequestAlgorithmRSA2048,
		SubjectAlternativeNames: []SSLCertificateRequestSAN{{Name: "www.example.com"}},
		Assignments:             []SSLCertificateRequestAssignment{{SubDomainName: "www.example.com"}},
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if created.ID != 44 || created.Status != SSLCertificateRequestStatusOpen {
		t.Errorf("Expected request 44 in status OPEN, got %+v", created)
	}

	var payload map[string]any
	if err := json.Unmarshal(requests.last(t).Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", requests.last(t).Body)
	}

	if payload["provider"] != "LETS_ENCRYPT" || payload["algorithm"] != "RSA2048" {
		t.Errorf("Expected provider and algorithm in the payload, got %v", payload)
	}

	sans, ok := payload["subjectAlternativeNames"].([]any)
	if !ok || len(sans) != 1 || sans[0].(map[string]any)["name"] != "www.example.com" {
		t.Errorf("Expected one subject alternative name in the payload, got %v", payload["subjectAlternativeNames"])
	}

	for _, key := range []string{"id", "status", "created", "modified", "sslProviderCredentialsId", "renewalInterval", "signatureAlgorithm"} {
		if _, present := payload[key]; present {
			t.Errorf("Expected [%s] to be omitted from the create payload, got %v", key, payload[key])
		}
	}
}

func TestCreateSSLCertificateRequestSendsEmptyCollections(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /ssl/requests": {Status: http.StatusCreated, Body: `{"error": false, "data": [{"id": 45, "algorithm": "RSA2048", "provider": "LETS_ENCRYPT", "status": "OPEN", "subjectAlternativeNames": [], "assignments": []}]}`},
	})

	request := &SSLCertificateRequest{Provider: SSLProviderLetsEncrypt, Algorithm: SSLCertificateRequestAlgorithmRSA2048}
	if _, err := api.CreateSSLCertificateRequest(request); err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(requests.last(t).Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", requests.last(t).Body)
	}

	if string(payload["subjectAlternativeNames"]) != "[]" || string(payload["assignments"]) != "[]" {
		t.Errorf("Expected nil collections to be sent as empty arrays, got %s", requests.last(t).Body)
	}

	if request.SubjectAlternativeNames != nil || request.Assignments != nil {
		t.Error("Expected the passed request not to be modified")
	}
}

func TestUpdateSSLCertificateRequest(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /ssl/requests/42": {Status: http.StatusOK, Body: `{"error": false, "data": [` + sslCertificateRequestJSON + `]}`},
	})

	modified := time.Date(2026, 8, 2, 11, 30, 0, 0, time.FixedZone("CEST", 2*60*60))
	updated, err := api.UpdateSSLCertificateRequest(&SSLCertificateRequest{
		ID:        42,
		Modified:  &types.DateTime{Time: modified},
		Algorithm: SSLCertificateRequestAlgorithmECDSA256,
		Provider:  SSLProviderSectigo,
		SubjectAlternativeNames: []SSLCertificateRequestSAN{
			{ID: 7, Name: "www.example.com"},
			{ID: 8, Name: "*.example.com"},
		},
		Assignments:              []SSLCertificateRequestAssignment{{ID: 9, SubDomainName: "www.example.com"}},
		SSLProviderCredentialsID: 3,
		RenewalInterval:          30,
		SignatureAlgorithm:       SSLCertificateRequestSignatureAlgorithmSHA384,
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	assertSSLCertificateRequest(t, updated)

	sent := requests.last(t)
	if sent.Method != http.MethodPut || sent.Path != "/ssl/requests/42" {
		t.Errorf("Expected PUT /ssl/requests/42 but got %s %s", sent.Method, sent.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	if payload["id"] != float64(42) || payload["modified"] != "2026-08-02T11:30:00+0200" {
		t.Errorf("Expected id and modified in the update payload, got %v", payload)
	}

	if payload["algorithm"] != "ECDSA256" || payload["sslProviderCredentialsId"] != float64(3) || payload["renewalInterval"] != float64(30) {
		t.Errorf("Expected algorithm, credentials and renewal interval in the update payload, got %v", payload)
	}

	sans, ok := payload["subjectAlternativeNames"].([]any)
	if !ok || len(sans) != 2 || sans[0].(map[string]any)["id"] != float64(7) {
		t.Errorf("Expected the subject alternative names with their ids in the payload, got %v", payload["subjectAlternativeNames"])
	}
}

func TestDeleteSSLCertificateRequest(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"DELETE /ssl/requests/42": {Status: http.StatusNoContent, Body: ""},
	})

	request := &SSLCertificateRequest{ID: 42, Modified: types.DateTimeNow()}
	deleted, err := api.DeleteSSLCertificateRequest(request)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if deleted != request {
		t.Error("Expected the passed request to be returned")
	}

	if sent := requests.last(t); sent.Method != http.MethodDelete || sent.Path != "/ssl/requests/42" {
		t.Errorf("Expected DELETE /ssl/requests/42 but got %s %s", sent.Method, sent.Path)
	}
}

func TestUpdateSSLCertificateRequestConfiguration(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /ssl/requests/42/ssl-configuration": {Status: http.StatusOK, Body: `{"error": false, "data": [` + sslCertificateRequestJSON + `]}`},
	})

	request, err := api.UpdateSSLCertificateRequestConfiguration(42, "2023-mozilla-modern")
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if request.ID != 42 {
		t.Errorf("Expected request [%d] but got [%d]", 42, request.ID)
	}

	sent := requests.last(t)
	if string(sent.Body) != `{"sslConfigurationName":"2023-mozilla-modern"}` {
		t.Errorf("Expected the configuration name as payload but got [%s]", sent.Body)
	}
}

func TestCheckSSLCertificateRequestDomains(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /ssl/requests/check-domains": {Status: http.StatusOK, Body: `{"error": false, "data": [{"domains": {
			"www.example.com": {"exists": true, "isMyraNS": true, "challengeName": null, "expectedCName": null},
			"*.example.org": {"exists": true, "isMyraNS": false, "challengeName": "_acme-challenge.example.org.", "expectedCName": "_acme--challenge-example-org.ax4z.com."},
			"missing.example.net": {"exists": false, "isMyraNS": false, "challengeName": null, "expectedCName": null}
		}}]}`},
	})

	checks, err := api.CheckSSLCertificateRequestDomains([]string{"www.example.com", "*.example.org", "missing.example.net"})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(checks) != 3 {
		t.Fatalf("Expected [%d] results but got [%d]", 3, len(checks))
	}

	if !checks["www.example.com"].IsMyraNS || !checks["www.example.com"].Exists {
		t.Errorf("Expected www.example.com to be served by Myra name servers, got %+v", checks["www.example.com"])
	}

	if checks["*.example.org"].IsMyraNS || checks["*.example.org"].ExpectedCName != "_acme--challenge-example-org.ax4z.com." {
		t.Errorf("Expected a CNAME challenge for *.example.org, got %+v", checks["*.example.org"])
	}

	if checks["missing.example.net"].Exists {
		t.Errorf("Expected missing.example.net not to exist, got %+v", checks["missing.example.net"])
	}

	if string(requests.last(t).Body) != `{"domains":["www.example.com","*.example.org","missing.example.net"]}` {
		t.Errorf("Expected the domains as payload but got [%s]", requests.last(t).Body)
	}
}

func TestCheckSSLCertificateRequestDomainsEmptyResult(t *testing.T) {
	// PHP serializes an empty associative array as [] rather than {}.
	api, _ := newTestAPI(t, map[string]testResponse{
		"POST /ssl/requests/check-domains": {Status: http.StatusOK, Body: `{"error": false, "data": [{"domains": []}]}`},
	})

	checks, err := api.CheckSSLCertificateRequestDomains([]string{"unresolvable.example.test"})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(checks) != 0 {
		t.Errorf("Expected an empty result but got %+v", checks)
	}
}

func TestCheckSSLCertificateRequestDomainsWithoutDomains(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{})

	for _, domains := range [][]string{nil, {}} {
		if _, err := api.CheckSSLCertificateRequestDomains(domains); err == nil {
			t.Errorf("Expected an error for domains %v", domains)
		}
	}

	if len(requests.all()) != 0 {
		t.Error("Expected no request to be sent without domains")
	}
}
