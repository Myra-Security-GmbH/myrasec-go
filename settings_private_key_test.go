package myrasec

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestSettingsPrivateKeyRequests(t *testing.T) {
	variants := []struct {
		name, path, method, readBody string
		read                         func(*API) (*Settings, error)
		update                       func(*API, *Settings) error
		partial                      func(*API, map[string]any) error
	}{
		{"flat", "/domain/1/www.example.test/settings", http.MethodPost,
			`{"myra_ssl_certificate_key":["<redacted>"],"only_https":true}`,
			func(api *API) (*Settings, error) {
				return api.ListSettingsContext(context.Background(), 1, "www.example.test", nil)
			},
			func(api *API, s *Settings) error {
				_, err := api.UpdateSettingsContext(context.Background(), s, 1, "www.example.test")
				return err
			},
			func(api *API, s map[string]any) error {
				_, err := api.UpdateSettingsPartialContext(context.Background(), s, 1, "www.example.test")
				return err
			},
		},
		{"tag", "/tag/1/settings", http.MethodPut,
			`{"settings":{"myra_ssl_certificate_key":["<redacted>"],"only_https":true}}`,
			func(api *API) (*Settings, error) { return api.ListTagSettingsContext(context.Background(), 1) },
			func(api *API, s *Settings) error {
				_, err := api.UpdateTagSettingsContext(context.Background(), s, 1)
				return err
			},
			func(api *API, s map[string]any) error {
				_, err := api.UpdateTagSettingsPartialContext(context.Background(), s, 1)
				return err
			},
		},
	}
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			api, requests := newTestAPI(t, map[string]testResponse{
				"GET " + v.path:         {Status: 200, Body: v.readBody},
				v.method + " " + v.path: {Status: 200, Body: `{"error":false,"data":[{}]}`},
			})
			t.Run("read_modify_write", func(t *testing.T) {
				settings, err := v.read(api)
				if err != nil {
					t.Fatal(err)
				}
				if RedactedSettingsPrivateKey != "<redacted>" {
					t.Fatalf("RedactedSettingsPrivateKey = %q, want <redacted>", RedactedSettingsPrivateKey)
				}
				if !reflect.DeepEqual(settings.MyraSSLCertificateKey, []string{RedactedSettingsPrivateKey}) {
					t.Fatalf("decoded key = %#v, want a single server marker", settings.MyraSSLCertificateKey)
				}
				enabled := false
				settings.OnlyHTTPS = &enabled
				before := len(requests.all())
				if err := v.update(api, settings); err != nil {
					t.Fatal(err)
				}
				payload := privateKeyRequestPayload(t, requests, v.method, v.path, before)
				if !reflect.DeepEqual(payload["myra_ssl_certificate_key"], []any{"<redacted>"}) || payload["only_https"] != false {
					t.Fatalf("HTTP update must preserve the returned marker while changing another setting, got %#v", payload)
				}
			})
			for _, tc := range []struct {
				name string
				keys []string
			}{{"typed_nil", nil}, {"typed_empty", []string{}}} {
				t.Run(tc.name, func(t *testing.T) {
					before := len(requests.all())
					if err := v.update(api, &Settings{MyraSSLCertificateKey: tc.keys}); err != nil {
						t.Fatal(err)
					}
					payload := privateKeyRequestPayload(t, requests, v.method, v.path, before)
					if _, ok := payload["myra_ssl_certificate_key"]; ok {
						t.Fatalf("typed nil and empty slices must be omitted from the HTTP request, got %#v", payload)
					}
				})
			}
			for _, tc := range []struct {
				name  string
				value any
				want  any
			}{{"partial_nil", nil, nil}, {"partial_empty", []string{}, []any{}}, {"partial_new_key", []string{"NEW_KEY"}, []any{"NEW_KEY"}}} {
				t.Run(tc.name, func(t *testing.T) {
					before := len(requests.all())
					if err := v.partial(api, map[string]any{"myra_ssl_certificate_key": tc.value}); err != nil {
						t.Fatal(err)
					}
					payload := privateKeyRequestPayload(t, requests, v.method, v.path, before)
					got, ok := payload["myra_ssl_certificate_key"]
					if !ok || !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("HTTP key value = %#v, present %t; want %#v", got, ok, tc.want)
					}
				})
			}
			t.Run("partial_omitted", func(t *testing.T) {
				before := len(requests.all())
				if err := v.partial(api, map[string]any{"only_https": false}); err != nil {
					t.Fatal(err)
				}
				payload := privateKeyRequestPayload(t, requests, v.method, v.path, before)
				if _, ok := payload["myra_ssl_certificate_key"]; ok {
					t.Fatalf("partial update must leave omitted keys out of the HTTP request, got %#v", payload)
				}
			})
		})
	}
}

func privateKeyRequestPayload(t *testing.T, requests *testRequests, method, path string, before int) map[string]any {
	t.Helper()
	if got := len(requests.all()) - before; got != 1 {
		t.Fatalf("update sent %d requests, want exactly one", got)
	}
	request := requests.last(t)
	if request.Method != method || request.Path != path {
		t.Fatalf("unexpected request: %s %s", request.Method, request.Path)
	}
	var payload map[string]any
	if err := json.Unmarshal(request.Body, &payload); err != nil {
		t.Fatal(err)
	}
	return payload
}
