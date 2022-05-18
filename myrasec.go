package myrasec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"golang.org/x/time/rate"

	"github.com/Myra-Security-GmbH/signature"
)

const (
	// APIBaseURL ...
	APIBaseURL = "https://apiv2.myracloud.com/%s"
	// DefaultAPILanguage ...
	DefaultAPILanguage = "en"
	// DefaultAPIUserAgent ...
	DefaultAPIUserAgent = "myrasec-go"
	// DefaultCachingTTL ...
	DefaultCachingTTL = 10
	// ErrorMsgRateLimitReached ...
	ErrorMsgRateLimitReached = "rate limit reached - too many requests"
)

// APILanguages ...
var APILanguages = map[string]bool{
	"en": true,
	"de": true,
}

//
// API holds the configuration for the current API client.
//
type API struct {
	BaseURL   string
	Language  string
	UserAgent string
	key       string
	secret    string
	cache     map[string]*responseCache
	caching   bool
	cacheTTL  int
	headers   http.Header
	client    *http.Client
	limiter   *rate.Limiter
}

//
// Response defines a response, returned by the MYRA API
//
type Response struct {
	Error         bool          `json:"error,omitempty"`
	ViolationList []*Violation  `json:"violationList,omitempty"`
	WarningList   []*Warning    `json:"warningList,omitempty"`
	TargetObject  []interface{} `json:"targetObject,omitempty"`
	Data          []interface{} `json:"data,omitempty"`
	List          []interface{} `json:"list,omitempty"`
	Page          int           `json:"page,omitempty"`
	Count         int           `json:"count,omitempty"`
	PageSize      int           `json:"pageSize,omitempty"`
}

//
// Violation defines a violation VO, returned by the MYRA API
//
type Violation struct {
	Path    string `json:"propertypath,omitempty"`
	Message string `json:"message,omitempty"`
}

//
// Warning defines a warning VO, returned by the MYRA API
//
type Warning struct {
	Path    string `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
}

//
// New returns a new MYRA API Client
//
func New(key, secret string) (*API, error) {
	if key == "" || secret == "" {
		return nil, errors.New("missing API credentials")
	}

	api := &API{
		BaseURL:   APIBaseURL,
		Language:  DefaultAPILanguage,
		UserAgent: DefaultAPIUserAgent,
		cache:     make(map[string]*responseCache),
		caching:   false,
		cacheTTL:  0,
		key:       key,
		secret:    secret,
		headers:   make(http.Header),
		client:    http.DefaultClient,
		limiter:   rate.NewLimiter(rate.Limit(5), 1), //5rps = 300req/min
	}
	return api, nil
}

//
// EnableCaching enables the caching of the response. Note: Only GET requests are cached.
// NOTE: The caching feature is still in development and may not work as expected.
//
func (api *API) EnableCaching() {
	api.caching = true
	api.cacheTTL = DefaultCachingTTL
}

//
// DisableCaching disables the caching of the response
// NOTE: The caching feature is still in development and may not work as expected.
//
func (api *API) DisableCaching() {
	api.caching = false
	api.cacheTTL = 0
}

//
// SetCachingTTL sets a ttl value for the caching. You have to first call the EnableCaching function to enable the caching.
// NOTE: The caching feature is still in development and may not work as expected.
//
func (api *API) SetCachingTTL(ttl int) {
	api.cacheTTL = ttl
}

//
// SetUserAgent sets the User-Agent for the API.
//
func (api *API) SetUserAgent(userAgent string) {
	api.UserAgent = userAgent
}

//
// SetLanguage changes the API language.
//
func (api *API) SetLanguage(language string) error {
	if _, ok := APILanguages[language]; !ok {
		return fmt.Errorf("passed language [\"%s\"] is not supported", language)
	}

	api.Language = language

	return nil
}

//
// call executes/sends the request to the MYRA API
//
func (api *API) call(definition APIMethod, payload ...interface{}) (interface{}, error) {

	req, err := api.prepareRequest(definition, payload...)
	if err != nil {
		return nil, err
	}

	if api.caching && api.inCache(req) {
		res := api.fromCache(req)
		if res != nil {
			return res, nil
		}
	}

	if err = api.limiter.Wait(context.Background()); err != nil {
		return nil, fmt.Errorf(ErrorMsgRateLimitReached)
	}

	sig := signature.New(api.secret, api.key, req)

	request, err := sig.Append()
	if err != nil {
		return nil, err
	}

	resp, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !intInSlice(resp.StatusCode, []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent,
	}) {
		_, err = errorMessage(resp)
		if err != nil {
			return nil, fmt.Errorf("%s (%d):\n%s", http.StatusText(resp.StatusCode), resp.StatusCode, err.Error())
		}
		return nil, fmt.Errorf("%s (%d)", http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	if definition.ResponseDecodeFunc != nil {
		res, err := definition.ResponseDecodeFunc(resp, definition)
		if err == nil && api.caching && isCachable(req) && !api.inCache(req) {
			api.cacheResponse(req, res)
		}

		return res, err
	}

	res, err := decodeDefaultResponse(resp, definition)
	if err == nil && api.caching && isCachable(req) && !api.inCache(req) {
		api.cacheResponse(req, res)
	}

	return res, err
}

//
// errorMessage returns the error message (error) from the response passed to the function.
//
func errorMessage(resp *http.Response) (*Response, error) {
	res, err := decodeBaseResponse(resp)
	if err != nil {
		return res, err
	}
	return res, nil
}

//
// decodeDefaultResponse handles the default decoding of a response.
//
func decodeDefaultResponse(resp *http.Response, definition APIMethod) (interface{}, error) {

	if definition.Method == http.MethodDelete {
		return nil, nil
	}

	res, err := decodeBaseResponse(resp)
	if err != nil {
		return nil, err
	}

	return prepareResult(*res, definition)
}

//
// decodeSingleElementResponse decodes the response for a single element (like GetDomain or GetDNSRecord)
//
func decodeSingleElementResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	res, err := decodeBaseResponse(resp)
	if err != nil {
		return nil, err
	}

	return prepareSingleElementResult(*res, definition)
}

//
// decodeBaseResponse decodes the passed http.Response to a Response struct for further processing
//
func decodeBaseResponse(resp *http.Response) (*Response, error) {
	var res Response
	err := json.NewDecoder(resp.Body).Decode(&res)
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

	return &res, nil
}

//
// prepareRequest ...
//
func (api *API) prepareRequest(definition APIMethod, payload ...interface{}) (*http.Request, error) {
	var err error
	var req *http.Request

	apiURL := fmt.Sprintf(api.BaseURL, definition.Action)
	switch definition.Method {
	case http.MethodGet:
		req, err = api.prepareGETRequest(apiURL, payload...)
	case http.MethodPost:
		req, err = api.preparePOSTRequest(apiURL, payload...)
	case http.MethodPut:
		req, err = api.preparePUTRequest(apiURL, payload...)
	case http.MethodDelete:
		req, err = api.prepareDELETERequest(apiURL, payload...)
	default:
		req, err = nil, fmt.Errorf("passed APIMethod definition has a not supported HTTP method - [%s] is not supported", definition.Method)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if api.UserAgent != "" {
		req.Header.Set("User-Agent", api.UserAgent)
	}

	return req, err
}

//
// prepareGETRequest handles/prepares GET requests
//
func (api *API) prepareGETRequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	if len(payload) <= 0 {
		return http.NewRequest(http.MethodGet, apiURL, nil)
	}

	if len(payload) > 1 {
		return nil, fmt.Errorf("unable to handle more than one payload in a GET call - payload should be a map[string]string")
	}

	baseURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	queryMap := payload[0].(map[string]string)
	params := baseURL.Query()
	for k, v := range queryMap {
		params.Add(k, v)
	}

	baseURL.RawQuery = params.Encode()

	return http.NewRequest(http.MethodGet, baseURL.String(), nil)
}

//
// preparePOSTRequest handles/prepares POST requests
//
func (api *API) preparePOSTRequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(data))
}

//
// preparePUTRequest handles/prepares PUT requests
//
func (api *API) preparePUTRequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPut, apiURL, bytes.NewBuffer(data))
}

//
// prepareDELETERequest handles/prepares DELETE requests
//
func (api *API) prepareDELETERequest(apiURL string, payload ...interface{}) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodDelete, apiURL, bytes.NewBuffer(data))
}

//
// prepareResult prepares the response for further processing
//
func prepareResult(response Response, definition APIMethod) (interface{}, error) {
	var result interface{}
	if response.TargetObject != nil {
		result = response.TargetObject[0]
	} else if response.List != nil {
		result = response.List
	} else if response.Data != nil {
		if definition.Method == http.MethodGet {
			result = response.Data
		} else {
			result = response.Data[0]
		}
	}

	tmp, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if definition.Result == nil {
		return tmp, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(tmp))
	retValue := reflect.New(reflect.TypeOf(definition.Result))
	res := retValue.Interface()
	decoder.Decode(res)

	return res, err
}

//
// prepareSingleElementResult ...
//
func prepareSingleElementResult(response Response, definition APIMethod) (interface{}, error) {
	result := response.Data[0]

	tmp, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if definition.Result == nil {
		return tmp, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(tmp))
	retValue := reflect.New(reflect.TypeOf(definition.Result))
	res := retValue.Interface()
	decoder.Decode(res)

	return res, err
}

//
// preparePayload ...
//
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
