package myrasec

import (
	"errors"
	"net/http"
	"testing"
)

func TestAPIErrorEmptyBody(t *testing.T) {
	// Access denied (missing feature or permission) answers 403 with an empty body.
	api, _ := newTestAPI(t, map[string]testResponse{
		"GET /ssl/requests/1": {Status: http.StatusForbidden, Body: ""},
	})

	_, err := api.GetSSLCertificateRequest(1)
	if err == nil {
		t.Fatal("Expected an error for a 403 response")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected an *APIError but got %T", err)
	}

	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code [%d] but got [%d]", http.StatusForbidden, apiErr.StatusCode)
	}

	if err.Error() != "Forbidden (403)" {
		t.Errorf("Expected error message [%s] but got [%s]", "Forbidden (403)", err.Error())
	}
}

func TestAPIErrorViolations(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"POST /ssl/requests": {
			Status: http.StatusBadRequest,
			Body: `{"error": true, "violationList": [
				{"propertyPath": "subjectAlternativeNames", "message": "Unknown domain: ghost.example.test"},
				{"propertyPath": "", "message": "Global violation"}
			], "errorMessage": "Validation failed"}`,
		},
	})

	_, err := api.CreateSSLCertificateRequest(&SSLCertificateRequest{Provider: SSLProviderLetsEncrypt})
	if err == nil {
		t.Fatal("Expected an error for a 400 response")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected an *APIError but got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code [%d] but got [%d]", http.StatusBadRequest, apiErr.StatusCode)
	}

	if len(apiErr.Violations) != 2 {
		t.Fatalf("Expected [%d] violations but got [%d]", 2, len(apiErr.Violations))
	}

	if apiErr.Violations[0].Path != "subjectAlternativeNames" {
		t.Errorf("Expected violation path [%s] but got [%s]", "subjectAlternativeNames", apiErr.Violations[0].Path)
	}

	if apiErr.ErrorMessage != "Validation failed" {
		t.Errorf("Expected error message [%s] but got [%s]", "Validation failed", apiErr.ErrorMessage)
	}

	expected := "Bad Request (400):\nsubjectAlternativeNames: Unknown domain: ghost.example.test\nGlobal violation\nValidation failed\n"
	if err.Error() != expected {
		t.Errorf("Expected error message [%q] but got [%q]", expected, err.Error())
	}
}

func TestAPIErrorNonJSONBody(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"GET /ssl/providers": {Status: http.StatusBadGateway, Body: "<html><body>502 Bad Gateway</body></html>"},
	})

	_, err := api.ListSSLProviderCredentials(nil)
	if err == nil {
		t.Fatal("Expected an error for a 502 response")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected an *APIError but got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status code [%d] but got [%d]", http.StatusBadGateway, apiErr.StatusCode)
	}

	if err.Error() != "Bad Gateway (502)" {
		t.Errorf("Expected error message [%s] but got [%s]", "Bad Gateway (502)", err.Error())
	}

	if apiErr.Body != "<html><body>502 Bad Gateway</body></html>" {
		t.Errorf("Expected the raw body to be kept but got [%s]", apiErr.Body)
	}
}

func TestAPIErrorMessageWithoutViolations(t *testing.T) {
	apiErr := &APIError{StatusCode: http.StatusBadRequest, ErrorMessage: "The object has been modified"}

	expected := "Bad Request (400):\nThe object has been modified\n"
	if apiErr.Error() != expected {
		t.Errorf("Expected error message [%q] but got [%q]", expected, apiErr.Error())
	}
}

func TestAPIErrorEnvelopeInSuccessfulResponse(t *testing.T) {
	// Some endpoints answer validation failures with HTTP 200 and error: true.
	api, _ := newTestAPI(t, map[string]testResponse{
		"GET /ssl/providers": {
			Status: http.StatusOK,
			Body:   `{"error": true, "violationList": [{"propertyPath": "provider", "message": "Invalid provider"}], "data": []}`,
		},
	})

	_, err := api.ListSSLProviderCredentials(nil)
	if err == nil {
		t.Fatal("Expected an error for an error envelope")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected an *APIError but got %T", err)
	}

	if apiErr.StatusCode != http.StatusOK || len(apiErr.Violations) != 1 {
		t.Errorf("Expected status 200 with one violation but got %+v", apiErr)
	}

	// The message of an envelope error carries no status prefix, as before the APIError type existed.
	if err.Error() != "provider: Invalid provider\n" {
		t.Errorf("Expected error message [%q] but got [%q]", "provider: Invalid provider\n", err.Error())
	}

	if (&APIError{StatusCode: http.StatusOK}).Error() != "The API returned an error." {
		t.Errorf("Expected the fallback message for an empty envelope")
	}
}
