package myrasec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	if api.key != key {
		t.Errorf("Expected key to be [%s] but got [%s]\n", key, api.key)
	}

	if api.secret != secret {
		t.Errorf("Expected key to be [%s] but got [%s]\n", key, api.key)
	}
}

func TestNewWithEmptyKey(t *testing.T) {
	key := ""
	secret := "123abc"

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty key should fail")
	}

	if err.Error() != "Missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "Missing API credentials", err.Error())
	}
}

func TestNewWithEmptySecret(t *testing.T) {
	key := "abc123"
	secret := ""

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty secret should fail")
	}

	if err.Error() != "Missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "Missing API credentials", err.Error())
	}
}

func TestNewWithEmptyParams(t *testing.T) {
	key := ""
	secret := ""

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty key/secret should fail")
	}

	if err.Error() != "Missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "Missing API credentials", err.Error())
	}
}

func TestSetUserAgent(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	if api.UserAgent != DefaultAPIUserAgent {
		t.Errorf("Expected default UserAgent to be [%s] but got [%s]", DefaultAPIUserAgent, api.UserAgent)
	}

	api.SetUserAgent("Testing")

	if api.UserAgent != "Testing" {
		t.Errorf("Expected UserAgent to be [%s] but got [%s]\n", "Testing", api.UserAgent)
	}
}

func TestSetLanguage(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	if api.Language != DefaultAPILanguage {
		t.Errorf("Expected default language to be [%s] but got [%s]\n", DefaultAPILanguage, api.Language)
	}

	err = api.SetLanguage("de")
	if err != nil {
		t.Errorf("Expected [%s] to be a valid language [%s]\n", "de", err.Error())
	}

	if api.Language != "de" {
		t.Errorf("Expected language to be [%s] but got [%s]\n", "de", api.Language)
	}

	err = api.SetLanguage("en")
	if err != nil {
		t.Errorf("Expected [%s] to be a valid language [%s]\n", "en", err.Error())
	}

	if api.Language != "en" {
		t.Errorf("Expected language to be [%s] but got [%s]\n", "en", api.Language)
	}

	err = api.SetLanguage("fr")
	if err == nil {
		t.Errorf("Expected [%s] to be invalid as a language setting\n", "fr")
	}
}

func TestPrepareRequest(t *testing.T) {
	definition := APIMethod{
		Name:   "Test",
		Action: "test",
		Method: http.MethodGet,
		Result: Domain{},
	}

	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.prepareRequest(definition, map[string]string{})
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	expectedURL := fmt.Sprintf(APIBaseURL, "test")
	if req.URL.String() != expectedURL {
		t.Errorf("Expected request URL path to be [%s] but got [%s]\n", expectedURL, req.URL.String())
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected request header \"Content-Type\" to be [%s] but got [%s]\n", "application/json", req.Header.Get("Content-Type"))
	}

	if req.Header.Get("User-Agent") != DefaultAPIUserAgent {
		t.Errorf("Expected request header \"User-Agent\" to be [%s] but got [%s]\n", DefaultAPIUserAgent, req.Header.Get("User-Agent"))
	}
}

func TestPrepareRequestInvalidMethod(t *testing.T) {
	definition := APIMethod{
		Name:   "Test",
		Action: "test",
		Method: "INVALID",
		Result: Domain{},
	}

	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.prepareRequest(definition, map[string]string{})
	if err == nil {
		t.Error("Expected an error as the passed HTTP method is not supported")
	}

	if req != nil {
		t.Error("Expected req to be null")
	}
}

func TestPrepareGETRequest(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.prepareGETRequest("/test", map[string]string{
		"test": "me",
		"foo":  "bar",
	})
	if err != nil {
		t.Errorf("Unexpected error passing a valid payload")
	}

	if req.Body != nil {
		t.Errorf("Expected request body to be nil")
	}

	if req.Method != http.MethodGet {
		t.Errorf("Expected request method to be [%s] but got [%s]\n", http.MethodGet, req.Method)
	}

	if req.URL.Path != "/test" {
		t.Errorf("Expected request URL path to be [%s] but got [%s]\n", "/test", req.URL.Path)
	}

	if req.URL.Query().Get("test") != "me" {
		t.Errorf("Expected request URL query value for [%s] to be [%s] but got [%s]\n", "test", "me", req.URL.Query().Get("test"))
	}
	if req.URL.Query().Get("foo") != "bar" {
		t.Errorf("Expected request URL query value for [%s] to be [%s] but got [%s]\n", "foo", "bar", req.URL.Query().Get("foo"))
	}
}

func TestPrepareGETRequestWithMultiplePayloads(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	_, err = api.prepareGETRequest("/test", map[string]string{
		"test": "me",
	}, map[string]string{
		"foo": "bar",
	})
	if err == nil {
		t.Errorf("Passing multiple payloads should fail")
	}
}

func TestPrepareGETRequestWithEmptyPayload(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.prepareGETRequest("/test", map[string]string{})
	if err != nil {
		t.Errorf("Unexpected error passing empty payload")
	}

	if req.Body != nil {
		t.Errorf("Expected request body to be nil")
	}

	if req.Method != http.MethodGet {
		t.Errorf("Expected request method to be [%s] but got [%s]\n", http.MethodGet, req.Method)
	}

	if req.URL.String() != "/test" {
		t.Errorf("Expected request URL to be [%s] but got [%s]\n", "/test", req.URL.String())
	}
}

func TestPrepareGETRequestPassingURLWithQuery(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.prepareGETRequest("/test?search=test", map[string]string{
		"test":   "me",
		"foo":    "bar",
		"search": "test",
	})
	if err != nil {
		t.Errorf("Unexpected error passing a valid payload")
	}

	if req.Body != nil {
		t.Errorf("Expected request body to be nil")
	}

	if req.Method != http.MethodGet {
		t.Errorf("Expected request method to be [%s] but got [%s]\n", http.MethodGet, req.Method)
	}

	if req.URL.Path != "/test" {
		t.Errorf("Expected request URL path to be [%s] but got [%s]\n", "/test", req.URL.Path)
	}

	if len(req.URL.Query()) != 3 {
		t.Errorf("Expected request to have [%d] query params but got [%d] (%s)\n", 3, len(req.URL.Query()), req.URL.Query().Encode())
	}

	if req.URL.Query().Get("search") != "test" {
		t.Errorf("Expected request URL query value for [%s] to be [%s] but got [%s]\n", "search", "test", req.URL.Query().Get("search"))
	}

	if req.URL.Query().Get("test") != "me" {
		t.Errorf("Expected request URL query value for [%s] to be [%s] but got [%s]\n", "test", "me", req.URL.Query().Get("test"))
	}

	if req.URL.Query().Get("foo") != "bar" {
		t.Errorf("Expected request URL query value for [%s] to be [%s] but got [%s]\n", "foo", "bar", req.URL.Query().Get("foo"))
	}
}

func TestPreparePOSTRequest(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.preparePOSTRequest("/test", nil)
	if err != nil {
		t.Errorf("Unexpected error preparing a POST request")
	}
	if req.Body == nil {
		t.Errorf("Expected request body not to be nil")
	}

	if req.Method != http.MethodPost {
		t.Errorf("Expected request method to be [%s] but got [%s]\n", http.MethodPost, req.Method)
	}

	if req.URL.Path != "/test" {
		t.Errorf("Expected request URL path to be [%s] but got [%s]\n", "/test", req.URL.Path)
	}
}

func TestPreparePUTRequest(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.preparePUTRequest("/test", nil)
	if err != nil {
		t.Errorf("Unexpected error preparing a PUT request")
	}
	if req.Body == nil {
		t.Errorf("Expected request body not to be nil")
	}

	if req.Method != http.MethodPut {
		t.Errorf("Expected request method to be [%s] but got [%s]\n", http.MethodPut, req.Method)
	}

	if req.URL.Path != "/test" {
		t.Errorf("Expected request URL path to be [%s] but got [%s]\n", "/test", req.URL.Path)
	}
}

func TestPrepareDELETERequest(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	req, err := api.prepareDELETERequest("/test", nil)
	if err != nil {
		t.Errorf("Unexpected error preparing a DELETE request")
	}
	if req.Body == nil {
		t.Errorf("Expected request body not to be nil")
	}

	if req.Method != http.MethodDelete {
		t.Errorf("Expected request method to be [%s] but got [%s]\n", http.MethodDelete, req.Method)
	}

	if req.URL.Path != "/test" {
		t.Errorf("Expected request URL path to be [%s] but got [%s]\n", "/test", req.URL.Path)
	}
}

func TestPreparePayloadWithMapPayload(t *testing.T) {
	payload := map[string]string{
		"test": "me",
		"foo":  "bar",
	}
	current, err := preparePayload(payload)
	if err != nil {
		t.Errorf("Unexpected error preparing a valid payload")
	}

	if current == nil {
		t.Errorf("Expected result to be not nil")
	}

	expected, err := json.Marshal(payload)
	if err != nil {
		t.Error("Unexpected error parsing the payload")
	}

	if !bytes.Equal(expected, current) {
		t.Errorf("Expected that [%s] is [%s]", string(current), string(expected))
	}

}

func TestPreparePayloadWithStructPayload(t *testing.T) {
	payload := &Domain{Name: "example.com"}
	current, err := preparePayload(payload)
	if err != nil {
		t.Errorf("Unexpected error preparing a valid payload")
	}

	if current == nil {
		t.Errorf("Expected result to be not nil")
	}

	expected, err := json.Marshal(payload)
	if err != nil {
		t.Error("Unexpected error parsing the payload")
	}

	if !bytes.Equal(expected, current) {
		t.Errorf("Expected that [%s] is [%s]", string(current), string(expected))
	}
}

func TestPreparePayloadWithMultiplePayloads(t *testing.T) {

	firstPayload := &Domain{Name: "example.com"}
	secondPayload := map[string]string{
		"test": "me",
		"foo":  "bar",
	}
	current, err := preparePayload(firstPayload, secondPayload)
	if err != nil {
		t.Errorf("Unexpected error preparing a valid payload")
	}

	if current == nil {
		t.Errorf("Expected result to be not nil")
	}

	expected, err := json.Marshal([]interface{}{
		firstPayload,
		secondPayload,
	})
	if err != nil {
		t.Error("Unexpected error parsing the payload")
	}

	if !bytes.Equal(expected, current) {
		t.Errorf("Expected that [%s] is [%s]", string(current), string(expected))
	}
}

func TestPreparePayloadWithNilPayload(t *testing.T) {

	current, err := preparePayload(nil)
	if err != nil {
		t.Errorf("Unexpected error preparing a valid payload")
	}

	if current == nil {
		t.Errorf("Expected result to be not nil")
	}

	expected, err := json.Marshal(nil)
	if err != nil {
		t.Error("Unexpected error parsing the payload")
	}

	if !bytes.Equal(expected, current) {
		t.Errorf("Expected that [%s] is [%s]", string(current), string(expected))
	}
}
