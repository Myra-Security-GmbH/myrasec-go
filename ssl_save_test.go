package myrasec

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

const sslCertificateSavedBody = `{"error": false, "data": [
	{"id": 5, "subjectAlternatives": ["example.com"], "subject": "CN=example.com", "subdomains": ["example.com"]}
], "warningList": [
	{"propertyPath": "intermediates", "message": "The chain has no intermediate certificate.", "code": "CHAIN_INTERMEDIATE_MISSING",
	 "translations": {"en": "The chain has no intermediate certificate.", "de": "Der Kette fehlt das Zwischenzertifikat."}, "parameters": {}},
	{"propertyPath": "cert", "message": "The certificate expires in 10 days.", "code": "LEAF_EXPIRES_SOON",
	 "translations": {"en": "The certificate expires in 10 days.", "de": "Das Zertifikat läuft in 10 Tagen ab."}, "parameters": {"days": 10, "date": "2026-11-28"}}
]}`

const sslCertificateRejectedBody = `{"error": true, "data": [], "violationList": [
	{"propertyPath": "cert", "message": "The certificate has expired.", "code": "LEAF_EXPIRED",
	 "translations": {"en": "The certificate has expired.", "de": "Das Zertifikat ist abgelaufen."}, "parameters": {"date": "2026-01-01"}},
	{"propertyPath": "intermediates[0].cert", "message": "Intermediate certificate 1 was not issued with the expected key.", "code": "INTERMEDIATE_WRONG_KEY",
	 "translations": {"en": "Intermediate certificate 1 was not issued with the expected key.", "de": "Zwischenzertifikat 1 wurde nicht mit dem erwarteten Schlüssel ausgestellt."}, "parameters": {"index": 0}}
]}`

func newSSLCertificateFixture() *SSLCertificate {
	return &SSLCertificate{Certificate: &Certificate{ID: 5, Cert: "-----BEGIN CERTIFICATE-----"}, Key: "key"}
}

func assertSavedWarnings(t *testing.T, warnings []*Warning) {
	t.Helper()

	if len(warnings) != 2 {
		t.Fatalf("Expected [%d] warnings but got [%d]", 2, len(warnings))
	}

	if warnings[0].Code != "CHAIN_INTERMEDIATE_MISSING" || warnings[0].Path != "intermediates" {
		t.Errorf("Expected the missing intermediate on [intermediates] but got [%s] on [%s]", warnings[0].Code, warnings[0].Path)
	}

	if warnings[1].Translations["de"] != "Das Zertifikat läuft in 10 Tagen ab." {
		t.Errorf("Expected the German translation but got [%s]", warnings[1].Translations["de"])
	}

	if warnings[1].Parameters["days"] != float64(10) || warnings[1].Parameters["date"] != "2026-11-28" {
		t.Errorf("Expected the expiry parameters but got [%v]", warnings[1].Parameters)
	}

	if warnings[1].String() != "cert: The certificate expires in 10 days. [LEAF_EXPIRES_SOON]" {
		t.Errorf("Unexpected warning text [%s]", warnings[1].String())
	}
}

func TestCreateSSLCertificateWithWarningsKeepsTheWarnings(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"POST /domain/1/certificates": {Status: http.StatusOK, Body: sslCertificateSavedBody},
	})

	saved, warnings, err := api.CreateSSLCertificateWithWarningsContext(context.Background(), newSSLCertificateFixture(), 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if saved.ID != 5 || saved.Subject != "CN=example.com" {
		t.Errorf("Expected the saved certificate but got [%+v]", saved.Certificate)
	}

	assertSavedWarnings(t, warnings)
}

func TestUpdateSSLCertificateWithWarningsKeepsTheWarnings(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /domain/1/certificates/5": {Status: http.StatusOK, Body: sslCertificateSavedBody},
	})

	saved, warnings, err := api.UpdateSSLCertificateWithWarningsContext(context.Background(), newSSLCertificateFixture(), 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if saved.ID != 5 {
		t.Errorf("Expected the saved certificate [5] but got [%d]", saved.ID)
	}

	if requests.last(t).Method != http.MethodPut {
		t.Errorf("Expected a PUT but got [%s]", requests.last(t).Method)
	}

	assertSavedWarnings(t, warnings)
}

func TestSaveSSLCertificateWithoutWarningsStillReturnsTheCertificate(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"POST /domain/1/certificates":  {Status: http.StatusOK, Body: sslCertificateSavedBody},
		"PUT /domain/1/certificates/5": {Status: http.StatusOK, Body: sslCertificateSavedBody},
	})

	created, err := api.CreateSSLCertificateContext(context.Background(), newSSLCertificateFixture(), 1)
	if err != nil || created.ID != 5 {
		t.Fatalf("Expected the created certificate [5] but got [%v] with error [%v]", created, err)
	}

	updated, err := api.UpdateSSLCertificate(newSSLCertificateFixture(), 1)
	if err != nil || updated.ID != 5 {
		t.Fatalf("Expected the updated certificate [5] but got [%v] with error [%v]", updated, err)
	}
}

func TestCreateSSLCertificateRejectedKeepsCodeFieldAndParameters(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"POST /domain/1/certificates": {Status: http.StatusBadRequest, Body: sslCertificateRejectedBody},
	})

	saved, warnings, err := api.CreateSSLCertificateWithWarningsContext(context.Background(), newSSLCertificateFixture(), 1)
	if saved != nil || warnings != nil {
		t.Errorf("Expected no certificate and no warnings but got [%v] and [%v]", saved, warnings)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected an *APIError but got %T", err)
	}

	if len(apiErr.Violations) != 2 {
		t.Fatalf("Expected [%d] violations but got [%d]", 2, len(apiErr.Violations))
	}

	wrongKey := apiErr.Violations[1]
	if wrongKey.Code != "INTERMEDIATE_WRONG_KEY" || wrongKey.Path != "intermediates[0].cert" || wrongKey.Parameters["index"] != float64(0) {
		t.Errorf("Expected the wrong-key intermediate on [intermediates[0].cert] but got [%+v]", wrongKey)
	}

	if apiErr.Violations[0].Translations["de"] != "Das Zertifikat ist abgelaufen." {
		t.Errorf("Expected the German translation but got [%s]", apiErr.Violations[0].Translations["de"])
	}

	expected := "Bad Request (400):\ncert: The certificate has expired. [LEAF_EXPIRED]\n" +
		"intermediates[0].cert: Intermediate certificate 1 was not issued with the expected key. [INTERMEDIATE_WRONG_KEY]\n"
	if err.Error() != expected {
		t.Errorf("Expected error message [%s] but got [%s]", expected, err.Error())
	}
}
