package myrasec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

	if err.Error() != "missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "missing API credentials", err.Error())
	}
}

func TestNewWithEmptySecret(t *testing.T) {
	key := "abc123"
	secret := ""

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty secret should fail")
	}

	if err.Error() != "missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "missing API credentials", err.Error())
	}
}

func TestNewWithEmptyParams(t *testing.T) {
	key := ""
	secret := ""

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty key/secret should fail")
	}

	if err.Error() != "missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "missing API credentials", err.Error())
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

func TestPrepareRequestGET(t *testing.T) {
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

	if req.Method != http.MethodGet {
		t.Errorf("Expected request method to be [%s] but got [%s]", http.MethodGet, req.Method)
	}
}

func TestPrepareRequestPOST(t *testing.T) {
	definition := APIMethod{
		Name:   "Test",
		Action: "test",
		Method: http.MethodPost,
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

	if req.Method != http.MethodPost {
		t.Errorf("Expected request method to be [%s] but got [%s]", http.MethodPost, req.Method)
	}
}

func TestPrepareRequestPUT(t *testing.T) {
	definition := APIMethod{
		Name:   "Test",
		Action: "test",
		Method: http.MethodPut,
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

	if req.Method != http.MethodPut {
		t.Errorf("Expected request method to be [%s] but got [%s]", http.MethodPut, req.Method)
	}
}

func TestPrepareRequestDELETE(t *testing.T) {
	definition := APIMethod{
		Name:   "Test",
		Action: "test",
		Method: http.MethodDelete,
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

	if req.Method != http.MethodDelete {
		t.Errorf("Expected request method to be [%s] but got [%s]", http.MethodDelete, req.Method)
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

func TestPreparePayloadWithInvalidPayload(t *testing.T) {
	_, err := preparePayload(func() {})
	if err == nil {
		t.Errorf("Expected to get an error on passing a non-valid payload")
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

func TestEnableDisableCaching(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	api.EnableCaching()

	if api.caching != true {
		t.Errorf("Expected to have caching enabled")
	}

	if api.cacheTTL != DefaultCachingTTL {
		t.Errorf("Expected cacheTTL to be %d", DefaultCachingTTL)
	}

	api.SetCachingTTL(10)

	if api.caching != true {
		t.Errorf("Expected to have caching enabled")
	}

	if api.cacheTTL != 10 {
		t.Errorf("Expected cacheTTL to be %d", 10)
	}

	api.DisableCaching()

	if api.caching != false {
		t.Errorf("Expected to have caching disabled")
	}

	if api.cacheTTL != 0 {
		t.Errorf("Expected cacheTTL to be %d", 0)
	}
}

func TestRetryFunctions(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	if api.maxRetries != DefaultRetryCount {
		t.Errorf("Expected maxRetries to be %d", DefaultRetryCount)
	}

	if api.retrySleep != DefaultRetrySleep {
		t.Errorf("Expected retrySleep to be %d", DefaultRetrySleep)
	}

	api.SetMaxRetries(3)

	if api.maxRetries != 3 {
		t.Errorf("Expected maxRetries to be %d", 3)
	}

	if api.retrySleep != DefaultRetrySleep {
		t.Errorf("Expected retrySleep to be %d", DefaultRetrySleep)
	}

	api.SetMaxRetries(DefaultRetryCount)
	api.SetRetrySleep(3)

	if api.maxRetries != DefaultRetryCount {
		t.Errorf("Expected maxRetries to be %d", DefaultRetryCount)
	}

	if api.retrySleep != 3 {
		t.Errorf("Expected retrySleep to be %d", 3)
	}
}

func TestErrorMessageWithErrorInResponse(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusBadRequest),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"error": true, "violationList": [{"propertypath": "test", "message": "this is a test message."}]}`)),
	}

	_, err := errorMessage(&resp)
	if err == nil {
		t.Errorf("Expected to have an error.")
	}

	if err.Error() != "test: this is a test message.\n" {
		t.Errorf("Expected to get error message [%s] but got [%s]", "test: this is a test message.", err.Error())
	}

}

func TestErrorMessageWithoutError(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusBadRequest),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"error": false}`)),
	}

	_, err := errorMessage(&resp)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s].", err.Error())
	}
}

func TestDecodeDefaultResponse(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"error": false, "pageSize": 10, "page": 1, "count": 1, "targetObject": [{"id": 1, "name": "example.com"}]}`)),
	}

	result, err := decodeDefaultResponse(&resp, APIMethod{
		Name:               "TEST",
		Action:             "TEST",
		Method:             http.MethodGet,
		Result:             Domain{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	})
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if result == nil {
		t.Errorf("Expected to get a domain but did not get anything.")
	}

	d := result.(*Domain)
	if d.ID != 1 {
		t.Errorf("Expected to get ID [%d] but got [%d]", 1, d.ID)
	}

	if d.Name != "example.com" {
		t.Errorf("Expected to get Name [%s] but got [%s]", "example.com", d.Name)
	}
}

func TestDecodeDefaultResponseWithInvalidBody(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"this will not work": func(){}}`)),
	}

	_, err := decodeDefaultResponse(&resp, APIMethod{
		Name:               "TEST",
		Action:             "TEST",
		Method:             http.MethodGet,
		Result:             Domain{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	})

	if err == nil {
		t.Errorf("Expected to get an error.")
	}
}

func TestDecodeDefaultResponseForDELETEMethod(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusNoContent),
	}

	result, err := decodeDefaultResponse(&resp, APIMethod{
		Name:               "TEST",
		Action:             "TEST",
		Method:             http.MethodDelete,
		Result:             nil,
		ResponseDecodeFunc: nil,
	})
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if result != nil {
		t.Errorf("Expected to get no result on MethodDelete.")
	}
}

func TestDecodeSingleElementResponse(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [{"id": 1, "name": "example.com"}]}`)),
	}

	result, err := decodeSingleElementResponse(&resp, APIMethod{
		Name:               "TEST",
		Action:             "TEST",
		Method:             http.MethodGet,
		Result:             Domain{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	})

	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if result == nil {
		t.Errorf("Expected to get a domain but did not get anything.")
	}

	d := result.(*Domain)
	if d.ID != 1 {
		t.Errorf("Expected to get ID [%d] but got [%d]", 1, d.ID)
	}

	if d.Name != "example.com" {
		t.Errorf("Expected to get Name [%s] but got [%s]", "example.com", d.Name)
	}
}

func TestDecodeSingleElementResponseWithCorruptBody(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"fooooo": func() {}}`)),
	}

	_, err := decodeSingleElementResponse(&resp, APIMethod{
		Name:               "TEST",
		Action:             "TEST",
		Method:             http.MethodGet,
		Result:             Domain{},
		ResponseDecodeFunc: decodeSingleElementResponse,
	})

	if err == nil {
		t.Errorf("Expected to get an error.")
	}
}

func TestDecodeBaseResponse(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"error": false, "pageSize": 10, "page": 1, "count": 0}`)),
	}

	r, err := decodeBaseResponse(&resp)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if r.Error != false {
		t.Errorf("Expected to have Error set to false")
	}

	if r.Page != 1 {
		t.Errorf("Expected to have Page set to [%d] but got [%d]", 1, r.Page)
	}

	if r.PageSize != 10 {
		t.Errorf("Expected to have PageSize set to [%d] but got [%d]", 10, r.PageSize)
	}

	if r.Count != 0 {
		t.Errorf("Expected to have Page set to [%d] but got [%d]", 0, r.Count)
	}
}

func TestDecodeBaseResponseWithCorruptBody(t *testing.T) {
	resp := http.Response{
		Status: strconv.Itoa(http.StatusBadRequest),
		Body:   ioutil.NopCloser(bytes.NewBufferString(`{"wrong": format}`)),
	}

	_, err := decodeBaseResponse(&resp)
	if err == nil {
		t.Errorf("Expected to get an error")
	}

}
