package myrasec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// The list and update responses never carry the token string.
//
//nolint:gosec // G101: the name matches, the value is a JSON fixture, not a credential
const apiTokenFixtureJSON = `{
	"objectType": "TokenVO", "id": 7,
	"created": "2026-08-01T09:00:00+0200", "modified": "2026-08-15T16:45:00+0200",
	"name": "CI deploy", "enabled": true, "deleted": false,
	"expiryDate": "2027-01-01T00:00:00+0100", "lastUsedDate": "2026-09-08T11:30:00+0200"
}`

func assertAPIToken(t *testing.T, token *APIToken) {
	t.Helper()

	if token.ID != 7 {
		t.Errorf("Expected ID [%d] but got [%d]", 7, token.ID)
	}

	if token.Created == nil || token.Created.Format("2006-01-02") != "2026-08-01" {
		t.Errorf("Expected Created [%s] but got %v", "2026-08-01", token.Created)
	}

	if token.Modified.Format("2006-01-02") != "2026-08-15" {
		t.Errorf("Expected Modified [%s] but got [%s]", "2026-08-15", token.Modified.Format("2006-01-02"))
	}

	if token.Name != "CI deploy" {
		t.Errorf("Expected Name [%s] but got [%s]", "CI deploy", token.Name)
	}

	if !token.Enabled {
		t.Error("Expected the token to be enabled")
	}

	if token.ExpiryDate == nil || token.ExpiryDate.Format("2006-01-02") != "2027-01-01" {
		t.Errorf("Expected ExpiryDate [%s] but got %v", "2027-01-01", token.ExpiryDate)
	}

	if token.LastUsedDate == nil || token.LastUsedDate.Format("2006-01-02T15:04") != "2026-09-08T11:30" {
		t.Errorf("Expected LastUsedDate [%s] but got %v", "2026-09-08T11:30", token.LastUsedDate)
	}
}

func TestListApiTokens(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /user/tokens": {Status: http.StatusOK, Body: `{"error": false, "pageSize": 50, "page": 1, "count": 2, "data": [` + apiTokenFixtureJSON + `,
			{"objectType": "TokenVO", "id": 8, "name": "never used", "enabled": false, "deleted": false, "expiryDate": null, "lastUsedDate": null}
		]}`},
	})

	tokens, err := api.ListApiTokensContext(context.Background(), map[string]string{"search": "deploy", "pageSize": "10"})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(tokens) != 2 {
		t.Fatalf("Expected [%d] tokens but got [%d]", 2, len(tokens))
	}

	assertAPIToken(t, &tokens[0])

	if tokens[0].Token != "" {
		t.Error("Expected the token string not to be returned by the list")
	}

	if tokens[1].Enabled {
		t.Error("Expected the second token to be disabled")
	}

	if tokens[1].ExpiryDate != nil || tokens[1].LastUsedDate != nil {
		t.Errorf("Expected null dates to decode to nil, got expiry %v and last used %v", tokens[1].ExpiryDate, tokens[1].LastUsedDate)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodGet || sent.Path != "/user/tokens" {
		t.Errorf("Expected GET /user/tokens but got %s %s", sent.Method, sent.Path)
	}

	if sent.Query.Get("search") != "deploy" || sent.Query.Get("pageSize") != "10" {
		t.Errorf("Expected search and pageSize query parameters, got [%s]", sent.Query.Encode())
	}
}

func TestListApiTokensEmpty(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"GET /user/tokens": {Status: http.StatusOK, Body: `{"error": false, "pageSize": 50, "page": 1, "count": 0, "data": []}`},
	})

	tokens, err := api.ListApiTokensContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(tokens) != 0 {
		t.Errorf("Expected no tokens but got %v", tokens)
	}
}

func TestCreateApiToken(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /user/tokens": {Status: http.StatusCreated, Body: `{"error": false, "violationList": [], "warningList": [], "data": [
			{"objectType": "TokenVO", "id": 7, "created": "2026-08-01T09:00:00+0200", "modified": "2026-08-15T16:45:00+0200",
			 "name": "CI deploy", "enabled": true, "deleted": false, "expiryDate": "2027-01-01T00:00:00+0100", "lastUsedDate": "2026-09-08T11:30:00+0200",
			 "token": "aGVsbG8gd29ybGQgdGhpcyBpcyBhIHRlc3QgdG9rZW4="}
		]}`},
	})

	expiry := time.Date(2027, 1, 1, 0, 0, 0, 0, time.FixedZone("CET", 60*60))
	created, err := api.CreateApiTokenContext(context.Background(), &APIToken{
		Name:       "CI deploy",
		Enabled:    true,
		ExpiryDate: &types.DateTime{Time: expiry},
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	assertAPIToken(t, created)

	if created.Token != "aGVsbG8gd29ybGQgdGhpcyBpcyBhIHRlc3QgdG9rZW4=" {
		t.Errorf("Expected the created token to carry the token string, got [%s]", created.Token)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodPost || sent.Path != "/user/tokens" {
		t.Errorf("Expected POST /user/tokens but got %s %s", sent.Method, sent.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	if payload["name"] != "CI deploy" || payload["enabled"] != true || payload["expiryDate"] != "2027-01-01T00:00:00+0100" {
		t.Errorf("Expected name, enabled and expiryDate in the create payload, got %v", payload)
	}

	for _, key := range []string{"id", "created", "modified", "token", "lastUsedDate"} {
		if _, present := payload[key]; present {
			t.Errorf("Expected [%s] to be omitted from the create payload, got %v", key, payload[key])
		}
	}
}

// The enabled flag is always sent: a zero-value APIToken asks for a disabled token
// instead of leaving the choice to the server.
func TestCreateApiTokenSendsDisabledFlag(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /user/tokens": {Status: http.StatusCreated, Body: `{"error": false, "data": [{"objectType": "TokenVO", "id": 9, "name": "disabled", "enabled": false, "token": "x"}]}`},
	})

	created, err := api.CreateApiTokenContext(context.Background(), &APIToken{Name: "disabled"})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if created.Enabled {
		t.Error("Expected the created token to be disabled")
	}

	var payload map[string]any
	if err := json.Unmarshal(requests.last(t).Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", requests.last(t).Body)
	}

	if enabled, present := payload["enabled"]; !present || enabled != false {
		t.Errorf("Expected enabled to be sent as false, got %v", payload)
	}

	if _, present := payload["expiryDate"]; present {
		t.Errorf("Expected a nil ExpiryDate to be omitted, got %v", payload["expiryDate"])
	}
}

func TestCreateApiTokenValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		token   *APIToken
		body    string
		path    string
		message string
	}{
		{
			name:    "expiry date in the past",
			token:   &APIToken{Name: "old", Enabled: true, ExpiryDate: &types.DateTime{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}},
			body:    `{"error": true, "violationList": [{"propertyPath": "expiryDate", "message": "The expire date is in the past."}], "warningList": [], "data": []}`,
			path:    "expiryDate",
			message: "The expire date is in the past.",
		},
		{
			name:    "reserved name",
			token:   &APIToken{Name: "reserved-name", Enabled: true},
			body:    `{"error": true, "violationList": [{"propertyPath": "name", "message": "The name \"reserved-name\" is reserved."}], "warningList": [], "data": []}`,
			path:    "name",
			message: `The name "reserved-name" is reserved.`,
		},
		{
			name:    "missing name",
			token:   &APIToken{Enabled: true},
			body:    `{"error": true, "violationList": [{"propertyPath": "name", "message": "This value should not be blank."}], "warningList": [], "data": []}`,
			path:    "name",
			message: "This value should not be blank.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, _ := newTestAPI(t, map[string]testResponse{
				"POST /user/tokens": {Status: http.StatusBadRequest, Body: tt.body},
			})

			_, err := api.CreateApiTokenContext(context.Background(), tt.token)
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

			if len(apiErr.Violations) != 1 || apiErr.Violations[0].Path != tt.path || apiErr.Violations[0].Message != tt.message {
				t.Errorf("Expected a violation on [%s] with message [%s], got %v", tt.path, tt.message, apiErr.Violations)
			}
		})
	}
}

func TestUpdateApiToken(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /user/tokens/7": {Status: http.StatusOK, Body: `{"error": false, "violationList": [], "warningList": [], "data": [
			{"objectType": "TokenVO", "id": 7, "created": "2026-08-01T09:00:00+0200", "modified": "2026-09-09T10:00:00+0200",
			 "name": "CI deploy (paused)", "enabled": false, "deleted": false, "expiryDate": null, "lastUsedDate": "2026-09-08T11:30:00+0200"}
		]}`},
	})

	// A caller may reuse the object CreateApiTokenContext returned, which carries the token string.
	modified := time.Date(2026, 8, 15, 16, 45, 0, 0, time.FixedZone("CEST", 2*60*60))
	existing := &APIToken{
		ID:       7,
		Modified: &types.DateTime{Time: modified},
		Name:     "CI deploy (paused)",
		Enabled:  false,
		Token:    "string-from-the-create-response",
	}

	updated, err := api.UpdateApiTokenContext(context.Background(), existing)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if existing.Token == "" {
		t.Error("Expected the passed token not to be modified")
	}

	if updated.Name != "CI deploy (paused)" || updated.Enabled {
		t.Errorf("Expected the updated token to be renamed and disabled, got name [%s] enabled [%v]", updated.Name, updated.Enabled)
	}

	if updated.Modified.Format("2006-01-02") != "2026-09-09" {
		t.Errorf("Expected the new Modified [%s] but got [%s]", "2026-09-09", updated.Modified.Format("2006-01-02"))
	}

	if updated.ExpiryDate != nil {
		t.Errorf("Expected the cleared expiry date to decode to nil, got %v", updated.ExpiryDate)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodPut || sent.Path != "/user/tokens/7" {
		t.Errorf("Expected PUT /user/tokens/7 but got %s %s", sent.Method, sent.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	if payload["id"] != float64(7) || payload["modified"] != "2026-08-15T16:45:00+0200" {
		t.Errorf("Expected id and modified in the update payload, got %v", payload)
	}

	if payload["name"] != "CI deploy (paused)" || payload["enabled"] != false {
		t.Errorf("Expected name and enabled in the update payload, got %v", payload)
	}

	// The API replaces the stored expiry date with the sent one, so an omitted
	// ExpiryDate must not be sent as null or a zero date either.
	if _, present := payload["expiryDate"]; present {
		t.Errorf("Expected a nil ExpiryDate to be omitted from the update payload, got %v", payload["expiryDate"])
	}

	if _, present := payload["token"]; present {
		t.Error("Expected the token string to be stripped from the update payload")
	}
}

// The API answers a foreign token with an empty 403, a stale Modified with a 400 carrying
// one violation without a property path and the name it reserves for internal use with a
// 400 on name, whether the caller sets it or the stored token carries it. All must surface
// as an *APIError.
func TestUpdateApiTokenErrors(t *testing.T) {
	tests := []struct {
		name    string
		token   *APIToken
		status  int
		body    string
		path    string
		message string
	}{
		{
			name:   "token of another user",
			token:  &APIToken{ID: 7, Modified: types.DateTimeNow(), Name: "CI deploy", Enabled: true},
			status: http.StatusForbidden,
			body:   "",
		},
		{
			name:    "stale modified",
			token:   &APIToken{ID: 7, Modified: types.DateTimeNow(), Name: "CI deploy", Enabled: true},
			status:  http.StatusBadRequest,
			body:    `{"error": true, "violationList": [{"message": "The record has been edited in the meantime."}], "warningList": [], "data": []}`,
			message: "The record has been edited in the meantime.",
		},
		{
			name:    "rename to the reserved name",
			token:   &APIToken{ID: 7, Modified: types.DateTimeNow(), Name: "reserved-name", Enabled: true},
			status:  http.StatusBadRequest,
			body:    `{"error": true, "violationList": [{"propertyPath": "name", "message": "The name \"reserved-name\" is reserved."}], "warningList": [], "data": []}`,
			path:    "name",
			message: `The name "reserved-name" is reserved.`,
		},
		{
			// The stored token carries the reserved name, the caller only tries to disable it.
			name:    "stored token carries the reserved name",
			token:   &APIToken{ID: 7, Modified: types.DateTimeNow(), Name: "renamed", Enabled: false},
			status:  http.StatusBadRequest,
			body:    `{"error": true, "violationList": [{"propertyPath": "name", "message": "The name \"reserved-name\" is reserved."}], "warningList": [], "data": []}`,
			path:    "name",
			message: `The name "reserved-name" is reserved.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, _ := newTestAPI(t, map[string]testResponse{
				"PUT /user/tokens/7": {Status: tt.status, Body: tt.body},
			})

			_, err := api.UpdateApiTokenContext(context.Background(), tt.token)
			if err == nil {
				t.Fatalf("Expected an error for a %d response", tt.status)
			}

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("Expected an *APIError but got %T", err)
			}

			if apiErr.StatusCode != tt.status {
				t.Errorf("Expected status code [%d] but got [%d]", tt.status, apiErr.StatusCode)
			}

			if tt.message == "" {
				if len(apiErr.Violations) != 0 {
					t.Errorf("Expected no violations for an empty body, got %v", apiErr.Violations)
				}
				return
			}

			if len(apiErr.Violations) != 1 || apiErr.Violations[0].Path != tt.path || apiErr.Violations[0].Message != tt.message {
				t.Errorf("Expected one violation on [%s] with message [%s], got %v", tt.path, tt.message, apiErr.Violations)
			}
		})
	}
}

func TestDeleteApiToken(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"DELETE /user/tokens/7": {Status: http.StatusNoContent, Body: ""},
	})

	existing := &APIToken{ID: 7, Token: "string-from-the-create-response"}
	deleted, err := api.DeleteApiTokenContext(context.Background(), existing)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if deleted != existing || existing.Token == "" {
		t.Error("Expected the passed token to be returned unmodified")
	}

	sent := requests.last(t)
	if sent.Method != http.MethodDelete || sent.Path != "/user/tokens/7" {
		t.Errorf("Expected DELETE /user/tokens/7 but got %s %s", sent.Method, sent.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	if payload["id"] != float64(7) {
		t.Errorf("Expected the id in the delete payload, got %v", payload)
	}

	if _, present := payload["token"]; present {
		t.Error("Expected the token string to be stripped from the delete payload")
	}
}

func TestDeleteApiTokenForbidden(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"DELETE /user/tokens/7": {Status: http.StatusForbidden, Body: ""},
	})

	_, err := api.DeleteApiTokenContext(context.Background(), &APIToken{ID: 7})
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
}

func TestCreateApiTokenWithCancelledContext(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /user/tokens": {Status: http.StatusCreated, Body: `{"error": false, "data": [{"id": 7, "token": "x"}]}`},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := api.CreateApiTokenContext(ctx, &APIToken{Name: "CI deploy", Enabled: true})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled but got %v", err)
	}

	if len(requests.all()) != 0 {
		t.Error("Expected no token to be created for a cancelled context")
	}
}
