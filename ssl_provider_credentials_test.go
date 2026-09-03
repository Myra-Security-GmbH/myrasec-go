package myrasec

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// The customer-facing representation never carries eabHmac or privateKey.
const sslProviderFixtureJSON = `{
	"objectType": "SslProviderCredentialsVO", "id": 3,
	"created": "2026-07-01T09:00:00+0200", "modified": "2026-07-15T16:45:00+0200",
	"name": "Sectigo OV", "provider": "SECTIGO",
	"cert": "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----",
	"email": "pki@example.com", "endpoint": "https://acme.sectigo.com/v2/OV",
	"eabKid": "kid-123", "comment": "OV account"
}`

func assertSSLProviderCredentials(t *testing.T, credentials *SSLProviderCredentials) {
	t.Helper()

	if credentials.ID != 3 {
		t.Errorf("Expected ID [%d] but got [%d]", 3, credentials.ID)
	}

	if credentials.Modified.Format("2006-01-02") != "2026-07-15" {
		t.Errorf("Expected Modified [%s] but got [%s]", "2026-07-15", credentials.Modified.Format("2006-01-02"))
	}

	if credentials.Name != "Sectigo OV" {
		t.Errorf("Expected Name [%s] but got [%s]", "Sectigo OV", credentials.Name)
	}

	if credentials.Provider != SSLProviderSectigo {
		t.Errorf("Expected Provider [%s] but got [%s]", SSLProviderSectigo, credentials.Provider)
	}

	if credentials.Endpoint != "https://acme.sectigo.com/v2/OV" {
		t.Errorf("Expected Endpoint [%s] but got [%s]", "https://acme.sectigo.com/v2/OV", credentials.Endpoint)
	}

	if credentials.Email != "pki@example.com" {
		t.Errorf("Expected Email [%s] but got [%s]", "pki@example.com", credentials.Email)
	}

	if credentials.EABKid != "kid-123" {
		t.Errorf("Expected EABKid [%s] but got [%s]", "kid-123", credentials.EABKid)
	}

	if credentials.EABHmac != "" || credentials.PrivateKey != "" {
		t.Error("Expected the secrets not to be returned by the API")
	}

	if credentials.Comment != "OV account" {
		t.Errorf("Expected Comment [%s] but got [%s]", "OV account", credentials.Comment)
	}
}

func TestListSSLProviderCredentials(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /ssl/providers": {Status: http.StatusOK, Body: `{"error": false, "pageSize": 50, "page": 1, "count": 2, "data": [` + sslProviderFixtureJSON + `,
			{"objectType": "SslProviderCredentialsVO", "id": 4, "name": "D-Trust", "provider": "DTRUST", "cert": "", "email": "", "endpoint": "", "eabKid": "kid-456"}
		]}`},
	})

	list, err := api.ListSSLProviderCredentials(map[string]string{"provider": "SECTIGO", "search": "OV"})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(list) != 2 {
		t.Fatalf("Expected [%d] credentials but got [%d]", 2, len(list))
	}

	assertSSLProviderCredentials(t, &list[0])

	if list[1].Provider != SSLProviderDTrust {
		t.Errorf("Expected second credentials provider [%s] but got [%s]", SSLProviderDTrust, list[1].Provider)
	}

	sent := requests.last(t)
	if sent.Query.Get("provider") != "SECTIGO" || sent.Query.Get("search") != "OV" {
		t.Errorf("Expected provider and search query parameters, got [%s]", sent.Query.Encode())
	}
}

func TestGetSSLProviderCredentials(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /ssl/providers/3": {Status: http.StatusOK, Body: `{"error": false, "data": [` + sslProviderFixtureJSON + `]}`},
	})

	credentials, err := api.GetSSLProviderCredentials(3)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	assertSSLProviderCredentials(t, credentials)

	if sent := requests.last(t); sent.Method != http.MethodGet || sent.Path != "/ssl/providers/3" {
		t.Errorf("Expected GET /ssl/providers/3 but got %s %s", sent.Method, sent.Path)
	}
}

func TestCreateSSLProviderCredentials(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /ssl/providers": {Status: http.StatusCreated, Body: `{"error": false, "data": [` + sslProviderFixtureJSON + `]}`},
	})

	created, err := api.CreateSSLProviderCredentials(&SSLProviderCredentials{
		Name:     "Sectigo OV",
		Provider: SSLProviderSectigo,
		Endpoint: "https://acme.sectigo.com/v2/OV",
		Email:    "pki@example.com",
		EABKid:   "kid-123",
		EABHmac:  "hmac-secret",
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	assertSSLProviderCredentials(t, created)

	var payload map[string]any
	if err := json.Unmarshal(requests.last(t).Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", requests.last(t).Body)
	}

	if payload["provider"] != "SECTIGO" || payload["eabKid"] != "kid-123" || payload["eabHmac"] != "hmac-secret" {
		t.Errorf("Expected provider and EAB values in the payload, got %v", payload)
	}

	// Empty cert and privateKey let the server generate the account key pair.
	for _, key := range []string{"id", "cert", "privateKey", "created", "modified", "comment"} {
		if _, present := payload[key]; present {
			t.Errorf("Expected [%s] to be omitted from the create payload, got %v", key, payload[key])
		}
	}
}

func TestUpdateSSLProviderCredentials(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /ssl/providers/3": {Status: http.StatusOK, Body: `{"error": false, "data": [` + sslProviderFixtureJSON + `]}`},
	})

	modified := time.Date(2026, 7, 15, 16, 45, 0, 0, time.FixedZone("CEST", 2*60*60))
	updated, err := api.UpdateSSLProviderCredentials(&SSLProviderCredentials{
		ID:       3,
		Modified: &types.DateTime{Time: modified},
		Name:     "Sectigo OV",
		Provider: SSLProviderSectigo,
		Endpoint: "https://acme.sectigo.com/v2/OV",
		EABKid:   "kid-123",
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	assertSSLProviderCredentials(t, updated)

	sent := requests.last(t)
	if sent.Method != http.MethodPut || sent.Path != "/ssl/providers/3" {
		t.Errorf("Expected PUT /ssl/providers/3 but got %s %s", sent.Method, sent.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	if payload["id"] != float64(3) || payload["modified"] != "2026-07-15T16:45:00+0200" {
		t.Errorf("Expected id and modified in the update payload, got %v", payload)
	}

	// Omitted secrets keep their stored value on the server.
	for _, key := range []string{"eabHmac", "privateKey", "cert"} {
		if _, present := payload[key]; present {
			t.Errorf("Expected empty [%s] to be omitted from the update payload, got %v", key, payload[key])
		}
	}
}

func TestDeleteSSLProviderCredentials(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"DELETE /ssl/providers/3": {Status: http.StatusNoContent, Body: ""},
	})

	credentials := &SSLProviderCredentials{ID: 3, Modified: types.DateTimeNow()}
	deleted, err := api.DeleteSSLProviderCredentials(credentials)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if deleted != credentials {
		t.Error("Expected the passed credentials to be returned")
	}

	if sent := requests.last(t); sent.Method != http.MethodDelete || sent.Path != "/ssl/providers/3" {
		t.Errorf("Expected DELETE /ssl/providers/3 but got %s %s", sent.Method, sent.Path)
	}
}

func TestListSSLProviderCertificates(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /ssl/providers/3/certificates": {Status: http.StatusOK, Body: `{"error": false, "pageSize": 50, "page": 1, "count": 1, "data": [
			{"objectType": "SslCertSummaryVO", "id": 500, "created": "2026-08-03T08:00:00+0200", "modified": "2026-08-03T08:00:00+0200",
			 "subject": "CN=www.example.com", "subjectAlternatives": ["www.example.com", "*.example.com"], "algorithm": "sha256WithRSAEncryption",
			 "validFrom": "2026-08-03T00:00:00+0200", "validTo": "2027-08-03T00:00:00+0200", "fingerprint": "AA:BB", "serialNumber": "01",
			 "wildcard": true, "extendedValidation": false, "managed": true, "multidomain": false, "sslConfigurationName": "Myra-Global-TLS-Default",
			 "requestId": 42, "domainId": 12, "subdomains": ["www.example.com"]}
		]}`},
	})

	certificates, err := api.ListSSLProviderCertificates(3, map[string]string{"pageSize": "10"})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(certificates) != 1 {
		t.Fatalf("Expected [%d] certificate but got [%d]", 1, len(certificates))
	}

	cert := certificates[0]
	if cert.ID != 500 || cert.RequestID != 42 || cert.DomainID != 12 {
		t.Errorf("Expected certificate 500 for request 42 on domain 12, got %+v", cert)
	}

	if !cert.Managed || !cert.Wildcard || cert.Multidomain {
		t.Errorf("Expected a managed wildcard certificate, got %+v", cert)
	}

	if cert.ValidTo.Format("2006-01-02") != "2027-08-03" {
		t.Errorf("Expected ValidTo [%s] but got [%s]", "2027-08-03", cert.ValidTo.Format("2006-01-02"))
	}

	if len(cert.SubjectAlternatives) != 2 || len(cert.Subdomains) != 1 {
		t.Errorf("Expected two SANs and one subdomain, got %+v", cert)
	}

	sent := requests.last(t)
	if sent.Path != "/ssl/providers/3/certificates" || sent.Query.Get("pageSize") != "10" {
		t.Errorf("Expected GET /ssl/providers/3/certificates?pageSize=10 but got %s?%s", sent.Path, sent.Query.Encode())
	}
}

func TestSSLProviderCredentialsAccessDenied(t *testing.T) {
	// Sub-group members and organizations without the Myra-Certificate feature receive 403.
	api, _ := newTestAPI(t, map[string]testResponse{
		"POST /ssl/providers": {Status: http.StatusForbidden, Body: ""},
	})

	_, err := api.CreateSSLProviderCredentials(&SSLProviderCredentials{Name: "x", Provider: SSLProviderDTrust, EABKid: "k", EABHmac: "h"})

	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected an *APIError with status 403 but got %v", err)
	}
}
