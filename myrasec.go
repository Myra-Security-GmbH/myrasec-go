package myrasec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/Myra-Security-GmbH/signature"
)

//
// APIBaseURL ...
//
const APIBaseURL = "https://api.myracloud.com/%s/rapi/%s"

//
// API holds the configuration for the current API client.
//
type API struct {
	BaseURL  string
	Language string
	key      string
	secret   string
	headers  http.Header
	client   *http.Client
}

//
// Response ...
//
type Response struct {
	Error         bool          `json:"error,omitempty"`
	ViolationList []*Violation  `json:"violationList,omitempty"`
	TargetObject  interface{}   `json:"targetObject,omitempty"`
	List          []interface{} `json:"list,omitempty"`
	Page          int           `json:"page,omitempty"`
	Count         int           `json:"count,omitempty"`
	PageSize      int           `json:"pageSize,omitempty"`
}

//
// Violation ...
//
type Violation struct {
	Path    string `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
}

//
// New returns a new MYRA API Client
//
func New(key, secret string) (*API, error) {
	if key == "" || secret == "" {
		return nil, errors.New("Missing credentials")
	}

	api := &API{
		BaseURL:  APIBaseURL,
		Language: "en",
		key:      key,
		secret:   secret,
		headers:  make(http.Header),
		client:   http.DefaultClient,
	}
	return api, nil
}

//
// call ...
//
func (api *API) call(definition APIMethod, payload ...interface{}) (interface{}, error) {

	req, err := api.prepareRequest(definition, payload...)
	if err != nil {
		return nil, err
	}

	fmt.Println(req.URL)

	sig := signature.New(api.secret, api.key, req)
	request, err := sig.Append()

	resp, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	if res.Error {
		var errorMsg string
		for _, e := range res.ViolationList {
			errorMsg += fmt.Sprintf("%s: %s\n", e.Path, e.Message)
		}
		return nil, errors.New(errorMsg)
	}
	return prepareResult(res, definition.Result)
}

//
// prepareRequest ...
//
func (api *API) prepareRequest(definition APIMethod, payload ...interface{}) (*http.Request, error) {
	var err error
	var req *http.Request

	apiURL := fmt.Sprintf(api.BaseURL, api.Language, definition.Action)
	switch definition.Method {
	case http.MethodGet:
		req, err = api.prepareGETRequest(apiURL, payload...)
		break
	case http.MethodPost:
		req, err = api.preparePOSTRequest(apiURL, payload...)
		break
	case http.MethodPut:
		req, err = api.preparePUTRequest(apiURL, payload...)
		break
	case http.MethodDelete:
		req, err = api.prepareDELETERequest(apiURL, payload...)
		break
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, err
}

//
// prepareGETRequest ...
//
func (api *API) prepareGETRequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	if len(payload) <= 0 {
		return http.NewRequest(http.MethodGet, apiURL, nil)
	}

	if len(payload) > 1 {
		return nil, fmt.Errorf("Unable to handle more than one payload in a GET call. payload should be a map[string]string")
	}

	baseURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	queryMap := payload[0].(map[string]string)
	params := url.Values{}
	for k, v := range queryMap {
		params.Add(k, v)
	}

	baseURL.RawQuery = params.Encode()
	return http.NewRequest(http.MethodGet, baseURL.String(), nil)
}

//
// preparePOSTRequest - ...
//
func (api *API) preparePOSTRequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(data))
}

//
// preparePUTRequest ...
//
func (api *API) preparePUTRequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPut, apiURL, bytes.NewBuffer(data))
}

//
// prepareDELETERequest ...
//
func (api *API) prepareDELETERequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodDelete, apiURL, bytes.NewBuffer(data))
}

//
// prepareResult ...
//
func prepareResult(response Response, definition interface{}) (interface{}, error) {

	var result interface{}
	if response.TargetObject != nil {
		result = response.TargetObject
	} else if response.List != nil {
		result = response.List
	}

	tmp, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if definition == nil {
		return tmp, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(tmp))
	retValue := reflect.New(reflect.TypeOf(definition))
	res := retValue.Interface()
	decoder.Decode(res)

	return res, err
}

func preparePayload(payload ...interface{}) ([]byte, error) {
	var pl interface{}
	pl = payload
	if len(payload) == 1 {
		pl = payload[0]
	}

	data, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}
	return data, nil
}
