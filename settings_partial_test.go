package myrasec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

// partialSettingsUpdate runs one of the partial settings updates against the test API.
type partialSettingsUpdate struct {
	name  string
	route string
	call  func(api *API, settings map[string]any) (any, error)
}

var partialSettingsUpdates = []partialSettingsUpdate{
	{
		name:  "tag settings",
		route: "PUT /tag/7/settings",
		call: func(api *API, settings map[string]any) (any, error) {
			return api.UpdateTagSettingsPartialContext(context.Background(), settings, 7)
		},
	},
	{
		name:  "subdomain settings",
		route: "POST /domain/1/www.example.com/settings",
		call: func(api *API, settings map[string]any) (any, error) {
			return api.UpdateSettingsPartialContext(context.Background(), settings, 1, "www.example.com")
		},
	},
}

func TestUpdateSettingsPartialResponses(t *testing.T) {
	// The API echoes the submitted attributes that are not null in data[0].
	tests := []struct {
		name     string
		settings map[string]any
		body     string
		want     map[string]any
	}{
		{
			// Every submitted attribute is null: PHP encodes the empty map as a JSON array
			// (terraform-provider-myrasec issue #280).
			name:     "all submitted attributes null, empty array",
			settings: map[string]any{"only_https": nil, "cache_enabled": nil},
			body:     `{"error":false,"violationList":[],"warningList":[],"data":[[]]}`,
			want:     map[string]any{},
		},
		{
			name:     "all submitted attributes null, empty object",
			settings: map[string]any{"only_https": nil, "cache_enabled": nil},
			body:     `{"error":false,"violationList":[],"warningList":[],"data":[{}]}`,
			want:     map[string]any{},
		},
		{
			name:     "submitted values echoed",
			settings: map[string]any{"only_https": true, "source_protocol": "https", "cache_enabled": nil},
			body:     `{"error":false,"violationList":[],"warningList":[],"data":[{"only_https":true,"source_protocol":"https"}]}`,
			want:     map[string]any{"only_https": true, "source_protocol": "https"},
		},
	}

	for _, update := range partialSettingsUpdates {
		for _, tt := range tests {
			t.Run(update.name+"/"+tt.name, func(t *testing.T) {
				api, requests := newTestAPI(t, map[string]testResponse{
					update.route: {Status: http.StatusOK, Body: tt.body},
				})

				result, err := update.call(api, tt.settings)
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				got, ok := result.(*map[string]any)
				if !ok {
					t.Fatalf("Expected a *map[string]any, got %T", result)
				}
				if !reflect.DeepEqual(*got, tt.want) {
					t.Errorf("Expected %v, got %v", tt.want, *got)
				}

				var sent map[string]any
				if err := json.Unmarshal(requests.last(t).Body, &sent); err != nil {
					t.Fatalf("Unexpected request body: %v", err)
				}
				if !reflect.DeepEqual(sent, tt.settings) {
					t.Errorf("Expected request body %v, got %v", tt.settings, sent)
				}
			})
		}
	}
}

// The deprecated typed updates call the same routes. Every Settings field is omitempty, so an
// empty Settings is sent as {} and the API answers "data":[[]].
func TestUpdateSettingsTypedEmptyArrayResponse(t *testing.T) {
	updates := []struct {
		name  string
		route string
		call  func(api *API, settings *Settings) (*Settings, error)
	}{
		{
			name:  "tag settings",
			route: "PUT /tag/7/settings",
			call: func(api *API, settings *Settings) (*Settings, error) {
				return api.UpdateTagSettingsContext(context.Background(), settings, 7)
			},
		},
		{
			name:  "subdomain settings",
			route: "POST /domain/1/www.example.com/settings",
			call: func(api *API, settings *Settings) (*Settings, error) {
				return api.UpdateSettingsContext(context.Background(), settings, 1, "www.example.com")
			},
		},
	}

	for _, update := range updates {
		t.Run(update.name, func(t *testing.T) {
			api, requests := newTestAPI(t, map[string]testResponse{
				update.route: {Status: http.StatusOK, Body: `{"error":false,"violationList":[],"warningList":[],"data":[[]]}`},
			})

			got, err := update.call(api, &Settings{})
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if !reflect.DeepEqual(*got, Settings{}) {
				t.Errorf("Expected empty Settings, got %+v", *got)
			}
			if body := string(requests.last(t).Body); body != "{}" {
				t.Errorf("Expected request body {}, got %s", body)
			}
		})
	}
}

func TestUpdateSettingsPartialErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "API error",
			body: `{"error":true,"violationList":[{"propertyPath":"only_https","message":"invalid"}],"data":[[]]}`,
		},
		{
			name: "empty data",
			body: `{"error":false,"violationList":[],"warningList":[],"data":[]}`,
		},
		{
			name: "non-empty array in data",
			body: `{"error":false,"violationList":[],"warningList":[],"data":[["only_https"]]}`,
		},
	}

	for _, update := range partialSettingsUpdates {
		for _, tt := range tests {
			t.Run(update.name+"/"+tt.name, func(t *testing.T) {
				api, _ := newTestAPI(t, map[string]testResponse{
					update.route: {Status: http.StatusOK, Body: tt.body},
				})

				result, err := update.call(api, map[string]any{"only_https": nil})
				if err == nil {
					t.Fatalf("Expected an error, got result %v", result)
				}
			})
		}
	}
}

func TestUpdateSettingsPartialAPIErrorType(t *testing.T) {
	for _, update := range partialSettingsUpdates {
		t.Run(update.name, func(t *testing.T) {
			api, _ := newTestAPI(t, map[string]testResponse{
				update.route: {
					Status: http.StatusOK,
					Body:   `{"error":true,"violationList":[{"propertyPath":"only_https","message":"invalid"}],"data":[[]]}`,
				},
			})

			_, err := update.call(api, map[string]any{"only_https": nil})

			if _, ok := errors.AsType[*APIError](err); !ok {
				t.Fatalf("Expected an *APIError, got %T: %v", err, err)
			}
		})
	}
}

func TestUpdateSettingsWAFAfterWhitelistMatch(t *testing.T) {
	updates := []struct {
		name  string
		route string
		call  func(api *API, settings *Settings) (*Settings, error)
	}{
		{
			name:  "subdomain settings",
			route: "POST /domain/1/www.example.com/settings",
			call: func(api *API, settings *Settings) (*Settings, error) {
				return api.UpdateSettingsContext(context.Background(), settings, 1, "www.example.com")
			},
		},
		{
			name:  "domain settings",
			route: "POST /domain/1/ALL-1./settings",
			call: func(api *API, settings *Settings) (*Settings, error) {
				return api.UpdateSettingsContext(context.Background(), settings, 1, "ALL-1.")
			},
		},
		{
			name:  "tag settings",
			route: "PUT /tag/7/settings",
			call: func(api *API, settings *Settings) (*Settings, error) {
				return api.UpdateTagSettingsContext(context.Background(), settings, 7)
			},
		},
	}

	tests := []struct {
		name     string
		settings *Settings
		sent     string
		echoed   string
		want     *bool
	}{
		{
			name:     "enabled",
			settings: &Settings{WAFAfterWhitelistMatch: BoolPtr(true)},
			sent:     `{"waf_after_whitelist_match":true}`,
			echoed:   `[{"waf_after_whitelist_match":true}]`,
			want:     BoolPtr(true),
		},
		{
			name:     "disabled",
			settings: &Settings{WAFAfterWhitelistMatch: BoolPtr(false)},
			sent:     `{"waf_after_whitelist_match":false}`,
			echoed:   `[{"waf_after_whitelist_match":false}]`,
			want:     BoolPtr(false),
		},
		{
			name:     "not set",
			settings: &Settings{WAFEnable: BoolPtr(true)},
			sent:     `{"waf_enable":true}`,
			echoed:   `[{"waf_enable":true}]`,
			want:     nil,
		},
	}

	for _, update := range updates {
		for _, tt := range tests {
			t.Run(update.name+"/"+tt.name, func(t *testing.T) {
				api, requests := newTestAPI(t, map[string]testResponse{
					update.route: {Status: http.StatusOK, Body: `{"error":false,"violationList":[],"warningList":[],"data":` + tt.echoed + `}`},
				})

				got, err := update.call(api, tt.settings)
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if !reflect.DeepEqual(got.WAFAfterWhitelistMatch, tt.want) {
					t.Errorf("Expected WAFAfterWhitelistMatch [%s], got [%s]", boolPtrStr(tt.want), boolPtrStr(got.WAFAfterWhitelistMatch))
				}
				if body := string(requests.last(t).Body); body != tt.sent {
					t.Errorf("Expected request body %s, got %s", tt.sent, body)
				}
			})
		}
	}
}

func TestUpdateSettingsPartialWAFAfterWhitelistMatch(t *testing.T) {
	tests := []struct {
		name     string
		settings map[string]any
		body     string
		want     map[string]any
	}{
		{
			name:     "enabled",
			settings: map[string]any{"waf_after_whitelist_match": true},
			body:     `{"error":false,"violationList":[],"warningList":[],"data":[{"waf_after_whitelist_match":true}]}`,
			want:     map[string]any{"waf_after_whitelist_match": true},
		},
		{
			name:     "removed",
			settings: map[string]any{"waf_after_whitelist_match": nil},
			body:     `{"error":false,"violationList":[],"warningList":[],"data":[[]]}`,
			want:     map[string]any{},
		},
	}

	for _, update := range partialSettingsUpdates {
		for _, tt := range tests {
			t.Run(update.name+"/"+tt.name, func(t *testing.T) {
				api, requests := newTestAPI(t, map[string]testResponse{
					update.route: {Status: http.StatusOK, Body: tt.body},
				})

				result, err := update.call(api, tt.settings)
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				got, ok := result.(*map[string]any)
				if !ok {
					t.Fatalf("Expected a *map[string]any, got %T", result)
				}
				if !reflect.DeepEqual(*got, tt.want) {
					t.Errorf("Expected %v, got %v", tt.want, *got)
				}

				var sent map[string]any
				if err := json.Unmarshal(requests.last(t).Body, &sent); err != nil {
					t.Fatalf("Unexpected request body: %v", err)
				}
				if !reflect.DeepEqual(sent, tt.settings) {
					t.Errorf("Expected request body %v, got %v", tt.settings, sent)
				}
			})
		}
	}
}
