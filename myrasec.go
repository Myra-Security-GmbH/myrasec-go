package myrasec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"golang.org/x/time/rate"

	"github.com/Myra-Security-GmbH/signature"
)

const (
	// APIBaseURL is the base URL template for the MYRA API.
	APIBaseURL = "https://apiv2.myracloud.com/%s"
	// DefaultAPILanguage is the default language for API responses.
	DefaultAPILanguage = "en"
	// DefaultAPIUserAgent is the default User-Agent header sent with API requests.
	DefaultAPIUserAgent = "myrasec-go"
	// DefaultCachingTTL is the default cache time-to-live in seconds.
	DefaultCachingTTL = 10
	// DefaultRetryCount is the default number of total attempts for a request.
	DefaultRetryCount = 1
	// DefaultRetrySleep is the default sleep duration in seconds between retries.
	DefaultRetrySleep = 0
	// ErrorMsgRateLimitReached is the error message returned when the rate limiter blocks a request.
	ErrorMsgRateLimitReached = "rate limit reached - too many requests"
)

// APILanguages defines the supported languages for API responses.
var APILanguages = map[string]bool{
	"en": true,
	"de": true,
}

// API holds the configuration for the current API client.
type API struct {
	BaseURL   string
	Language  string
	UserAgent string

	key        string
	secret     string
	token      string
	cache      map[string]*responseCache
	caching    bool
	cacheTTL   int
	headers    http.Header
	client     *http.Client
	limiter    *rate.Limiter
	maxRetries int
	retrySleep int

	methods map[string]APIMethod
}

// Response defines a response, returned by the MYRA API
type Response struct {
	Error         bool         `json:"error,omitempty"`
	ViolationList []*Violation `json:"violationList,omitempty"`
	ErrorMessage  string       `json:"errorMessage,omitempty"`
	WarningList   []*Warning   `json:"warningList,omitempty"`
	TargetObject  []any        `json:"targetObject,omitempty"`
	Data          []any        `json:"data,omitempty"`
	List          []any        `json:"list,omitempty"`
	Result        []any        `json:"result,omitempty"`
	Page          int          `json:"page,omitempty"`
	Count         int          `json:"count,omitempty"`
	PageSize      int          `json:"pageSize,omitempty"`
	Domain        []any        `json:"domain,omitempty"`
}

// Violation defines a violation VO, returned by the MYRA API
type Violation struct {
	Path    string `json:"propertypath,omitempty"`
	Message string `json:"message,omitempty"`
}

// Warning defines a warning VO, returned by the MYRA API
type Warning struct {
	Path    string `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
}

// New returns a new MYRA API Client
func New(key, secret string) (*API, error) {
	if key == "" || secret == "" {
		return nil, errors.New("missing API credentials")
	}

	return buildApi(key, secret, ""), nil
}

// NewWithToken returns a new MYRA API Client using API token authentication.
func NewWithToken(token string) (*API, error) {
	if token == "" {
		return nil, errors.New("missing API token")
	}

	return buildApi("", "", token), nil
}

func buildApi(key, secret, token string) *API {
	return &API{
		BaseURL:    getEnvOrDefault("MYRASEC_GO_BASE_URL", APIBaseURL),
		Language:   getEnvOrDefault("MYRASEC_GO_LANGUAGE", DefaultAPILanguage),
		UserAgent:  getEnvOrDefault("MYRASEC_GO_USER_AGENT", DefaultAPIUserAgent),
		cache:      make(map[string]*responseCache),
		caching:    false,
		cacheTTL:   0,
		key:        key,
		secret:     secret,
		token:      token,
		headers:    make(http.Header),
		client:     &http.Client{Timeout: 30 * time.Second},
		limiter:    rate.NewLimiter(rate.Limit(5), 1), // 5rps = 300req/min
		maxRetries: DefaultRetryCount,
		retrySleep: DefaultRetrySleep,
		methods:    initializeMethods(),
	}
}

// EnableCaching enables the caching of the response. Note: Only GET requests are cached.
// NOTE: The caching feature is still in development and may not work as expected.
func (api *API) EnableCaching() {
	api.caching = true
	api.cacheTTL = DefaultCachingTTL
}

// DisableCaching disables the caching of the response
// NOTE: The caching feature is still in development and may not work as expected.
func (api *API) DisableCaching() {
	api.caching = false
	api.cacheTTL = 0
}

// SetCachingTTL sets a ttl value for the caching. You have to first call the EnableCaching function to enable the caching.
// NOTE: The caching feature is still in development and may not work as expected.
func (api *API) SetCachingTTL(ttl int) {
	api.cacheTTL = ttl
}

// SetUserAgent sets the User-Agent for the API.
func (api *API) SetUserAgent(userAgent string) {
	api.UserAgent = userAgent
}

// SetLanguage changes the API language.
func (api *API) SetLanguage(language string) error {
	if _, ok := APILanguages[language]; !ok {
		return fmt.Errorf("passed language [\"%s\"] is not supported", language)
	}

	api.Language = language

	return nil
}

// SetMaxRetries sets the maxRetries value in the API struct. In case of a non-successful request, it will try (in total) n times.
func (api *API) SetMaxRetries(n int) {
	api.maxRetries = n
}

// SetRetrySleep sets a sleep value. It will wait for n-seconds to do the request again in case of retry operation.
func (api *API) SetRetrySleep(n int) {
	api.retrySleep = n
}

// SetProxy allows to set a custom proxyURL for the api client.
// Must be used after SetHTTPClient when SetHTTPClient is used.
func (api *API) SetProxy(proxyURL string) error {
	if proxyURL == "" {
		api.client.Transport = nil
		return nil
	}

	transport := &http.Transport{}
	if api.client.Transport != nil {
		var ok bool

		t, ok := api.client.Transport.(*http.Transport)
		if !ok {
			return fmt.Errorf("the client transport is not an http.Transport: %T", api.client.Transport)
		}

		transport = t.Clone()
	}

	purl, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("error setting proxy url [%q] is not a valid url: %v", proxyURL, err)
	}

	transport.Proxy = http.ProxyURL(purl)

	api.client.Transport = transport

	return nil
}

// SetHTTPClient sets the HTTP client used by the api client.
// Must be used before SetProxy if SetProxy is used.
func (api *API) SetHTTPClient(client *http.Client) error {
	if client == nil {
		return errors.New("nil http client")
	}

	api.client = client

	return nil
}

// call executes/sends the request to the MYRA API
func (api *API) call(definition APIMethod, payload ...any) (any, error) {
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

	resp, err := api.sendRequest(definition, payload...)
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

// sendRequest performs the concrete send-action
func (api *API) sendRequest(definition APIMethod, payload ...any) (*http.Response, error) {
	var retries int

	for {
		req, err := api.prepareRequest(definition, payload...)
		if err != nil {
			return nil, err
		}

		if err = api.limiter.Wait(context.Background()); err != nil {
			return nil, errors.New(ErrorMsgRateLimitReached)
		}

		var request *http.Request

		if api.key != "" && api.secret != "" {
			sig := signature.New(api.secret, api.key, req)

			request, err = sig.Append()
			if err != nil {
				return nil, err
			}
		}

		if api.token != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.token))
			request = req
		}

		if request == nil {
			return nil, errors.New("problem creating the request. Check API credentials")
		}

		resp, err := api.client.Do(request)
		if err != nil {
			return resp, err
		}

		retries++

		if resp.StatusCode != http.StatusInternalServerError || retries >= api.maxRetries {
			return resp, err
		}

		time.Sleep(time.Duration(api.retrySleep) * time.Second)
	}
}

// errorMessage returns the error message (error) from the response passed to the function.
func errorMessage(resp *http.Response) (*Response, error) {
	res, err := decodeBaseResponse(resp)
	if err != nil {
		return res, err
	}
	return res, nil
}

// decodeDefaultResponse handles the default decoding of a response.
func decodeDefaultResponse(resp *http.Response, definition APIMethod) (any, error) {
	if definition.Method == http.MethodDelete {
		return nil, nil
	}

	res, err := decodeBaseResponse(resp)
	if err != nil {
		return nil, err
	}

	return prepareResult(*res, definition)
}

// decodeSingleElementResponse decodes the response for a single element (like GetDomain or GetDNSRecord)
func decodeSingleElementResponse(resp *http.Response, definition APIMethod) (any, error) {
	res, err := decodeBaseResponse(resp)
	if err != nil {
		return nil, err
	}

	return prepareSingleElementResult(*res, definition)
}

// decodeBaseResponse decodes the passed http.Response to a Response struct for further processing
func decodeBaseResponse(resp *http.Response) (*Response, error) {
	var res Response
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if res.Error {
		return nil, formatAPIError(res.ErrorMessage, res.ViolationList)
	}

	return &res, nil
}

// formatAPIError builds an error from an API error response's errorMessage and violationList.
func formatAPIError(errorMessage string, violations []*Violation) error {
	var msg string
	for _, v := range violations {
		// Global violations (e.g. "Wrong format for provided dates!") carry no
		// property path; do not prefix them with a bare colon.
		if v.Path == "" {
			msg += fmt.Sprintf("%s\n", v.Message)
			continue
		}
		msg += fmt.Sprintf("%s: %s\n", v.Path, v.Message)
	}
	if errorMessage != "" {
		msg += fmt.Sprintf("%s\n", errorMessage)
	}
	if msg == "" {
		msg = "The API returned an error."
	}
	return errors.New(msg)
}

// prepareRequest builds an HTTP request from the given API method definition and payload.
func (api *API) prepareRequest(definition APIMethod, payload ...any) (*http.Request, error) {
	var err error
	var req *http.Request

	baseURL := api.BaseURL
	if definition.BaseURL != "" {
		baseURL = definition.BaseURL
	}
	apiURL := fmt.Sprintf(baseURL, definition.Action)
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

	if definition.AdditionalHeaders != nil {
		for k, v := range definition.AdditionalHeaders {
			req.Header.Set(k, v)
		}
	}

	return req, err
}

// prepareGETRequest handles/prepares GET requests
func (api *API) prepareGETRequest(apiURL string, payload ...any) (*http.Request, error) {
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

	queryMap, ok := payload[0].(map[string]string)
	if !ok {
		return nil, fmt.Errorf("GET request payload must be map[string]string, got %T", payload[0])
	}
	params := baseURL.Query()
	for k, v := range queryMap {
		params.Add(k, v)
	}

	baseURL.RawQuery = params.Encode()

	return http.NewRequest(http.MethodGet, baseURL.String(), nil)
}

// preparePOSTRequest handles/prepares POST requests
func (api *API) preparePOSTRequest(apiURL string, payload ...any) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(data))
}

// preparePUTRequest handles/prepares PUT requests
func (api *API) preparePUTRequest(apiURL string, payload ...any) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodPut, apiURL, bytes.NewBuffer(data))
}

// prepareDELETERequest handles/prepares DELETE requests
func (api *API) prepareDELETERequest(apiURL string, payload ...any) (*http.Request, error) {
	data, err := preparePayload(payload...)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(http.MethodDelete, apiURL, bytes.NewBuffer(data))
}

// prepareResult prepares the response for further processing
func prepareResult(response Response, definition APIMethod) (any, error) {
	var result any
	if response.TargetObject != nil {
		if len(response.TargetObject) == 0 {
			return nil, fmt.Errorf("empty TargetObject in API response")
		}
		result = response.TargetObject[0]
	} else if response.List != nil {
		result = response.List
	} else if response.Data != nil {
		if definition.Method == http.MethodGet {
			result = response.Data
		} else {
			if len(response.Data) == 0 {
				return nil, fmt.Errorf("empty Data in API response")
			}
			result = response.Data[0]
		}
	} else if response.Result != nil {
		result = response.Result
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
	err = decoder.Decode(res)

	return res, err
}

// prepareSingleElementResult extracts and decodes a single element from the API response.
func prepareSingleElementResult(response Response, definition APIMethod) (any, error) {
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("unable to detect element in API response")
	}

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
	err = decoder.Decode(res)

	return res, err
}

// preparePayload serializes the given payload to JSON bytes for use in HTTP request bodies.
func preparePayload(payload ...any) ([]byte, error) {
	var pl any
	pl = payload
	if len(payload) == 1 {
		pl = payload[0]
	}

	var data []byte
	var err error

	switch v := pl.(type) {
	case []byte:
		data = v
	default:
		data, err = json.Marshal(pl)
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

// returns the value of the passed environment variable (key) or - if not found - the fallback
func getEnvOrDefault(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
