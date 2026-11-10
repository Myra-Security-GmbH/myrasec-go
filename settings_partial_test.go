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
	tests := []struct {
		name string
		body string
		want map[string]any
	}{
		{
			// The API answers "data":[[]] when every setting was reset: PHP encodes the
			// empty map as a JSON array (terraform-provider-myrasec issue #280).
			name: "all settings reset, empty array",
			body: `{"error":false,"violationList":[],"warningList":[],"data":[[]]}`,
			want: map[string]any{},
		},
		{
			name: "empty object",
			body: `{"error":false,"violationList":[],"warningList":[],"data":[{}]}`,
			want: map[string]any{},
		},
		{
			name: "remaining settings",
			body: `{"error":false,"violationList":[],"warningList":[],"data":[{"only_https":true,"source_protocol":"https"}]}`,
			want: map[string]any{"only_https": true, "source_protocol": "https"},
		},
	}

	for _, update := range partialSettingsUpdates {
		for _, tt := range tests {
			t.Run(update.name+"/"+tt.name, func(t *testing.T) {
				api, requests := newTestAPI(t, map[string]testResponse{
					update.route: {Status: http.StatusOK, Body: tt.body},
				})

				result, err := update.call(api, map[string]any{"only_https": nil, "cache_enabled": nil})
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
				if v, ok := sent["only_https"]; !ok || v != nil {
					t.Errorf("Expected only_https to be sent as null, got %v (present: %v)", v, ok)
				}
			})
		}
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
