package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getAPITokenMethods returns API token related API calls.
func getAPITokenMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listApiTokens": {
			Name:   "listApiTokens",
			Action: "user/tokens",
			Method: http.MethodGet,
			Result: []APIToken{},
		},
		"createApiToken": {
			Name:   "createApiToken",
			Action: "user/tokens",
			Method: http.MethodPost,
			Result: APIToken{},
		},
		"updateApiToken": {
			Name:   "updateApiToken",
			Action: "user/tokens/%d",
			Method: http.MethodPut,
			Result: APIToken{},
		},
		"deleteApiToken": {
			Name:   "deleteApiToken",
			Action: "user/tokens/%d",
			Method: http.MethodDelete,
			Result: APIToken{},
		},
	}
}

// APIToken represents a bearer token used to access the API, the credential
// NewWithToken takes. A token carries the permissions of the user that created
// it. Unlike an API key it can be disabled and it can carry an expiry date.
type APIToken struct {
	// ID is the unique identifier for the API token.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the API token. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates the timestamp when the API token was generated.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp indicating when the API token was generated (ISO 8601 format). This is a server-managed, read-only value."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format), serving as a version identifier. This field is required for updates to ensure data consistency."`

	// Name is an arbitrary, user-defined label for the API token. Required.
	// The API reserves one name for internal use: creating a token with it
	// and updating a token that carries it are rejected with 400 Bad Request,
	// the violation names it.
	Name string `json:"name,omitempty" jsonschema:"An arbitrary, user-defined label for the API token. Required. The API reserves one name for internal use, creating a token with it and updating a token that carries it are rejected."`

	// Enabled indicates whether the token can be used for authentication.
	// A disabled token is rejected with 403 Forbidden. The flag is always
	// sent, so set it to true when creating or updating a token.
	Enabled bool `json:"enabled" jsonschema:"Indicates whether the token can be used for authentication. A disabled token is rejected with 403 Forbidden. Always sent, so set it to true when creating or updating a token."`

	// ExpiryDate is the optional point in time after which the token is
	// rejected. The API validates it against the start of the current day
	// only: a value before today's midnight is rejected with 400 Bad Request,
	// a value earlier today is accepted but the token is rejected right away,
	// so pass a point in the future. Updates replace the stored value, so an
	// update without ExpiryDate clears the expiry date.
	ExpiryDate *types.DateTime `json:"expiryDate,omitempty" jsonschema:"The optional point in time after which the token is rejected (ISO 8601 format). Validated against the start of the current day only, so pass a point in the future. An update without this attribute clears the expiry date."`

	// LastUsedDate records when the token was last used for authentication.
	// This is a server-managed, read-only value; it is nil for a token that
	// has never been used.
	LastUsedDate *types.DateTime `json:"lastUsedDate,omitempty" jsonschema:"The timestamp of the last authentication with this token (ISO 8601 format). Server-managed, read-only; absent for a token that has never been used."`

	// Token is the bearer token string.
	// Note: This value is returned only once upon creation and cannot be
	// retrieved later. It is never sent back to the API: update and delete
	// strip it from the request body.
	Token string `json:"token,omitempty" jsonschema:"The bearer token string. Visible only once upon creation; it cannot be retrieved later. Never sent back to the API."`
}

// ListApiTokensContext returns the API tokens of the authenticated user.
// Supported query parameters are "search", "page", "pageSize" and "sort".
// The returned tokens never carry the Token attribute, it is visible only
// in the response of CreateApiTokenContext.
func (api *API) ListApiTokensContext(ctx context.Context, params map[string]string) ([]APIToken, error) {
	if _, ok := api.methods["listApiTokens"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listApiTokens")
	}

	definition := api.methods["listApiTokens"]

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]APIToken)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateApiTokenContext creates a new API token for the authenticated user
// using the MYRA API. Name is required and Enabled must be true for a token
// that can be used. The returned token carries the Token attribute, the
// only time the API reveals it; store it right away.
func (api *API) CreateApiTokenContext(ctx context.Context, token *APIToken) (*APIToken, error) {
	if _, ok := api.methods["createApiToken"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createApiToken")
	}

	definition := api.methods["createApiToken"]

	result, err := api.call(ctx, definition, token)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*APIToken)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateApiTokenContext updates the passed API token using the MYRA API.
// ID and Modified identify the token and its version. The update replaces
// Name, Enabled and ExpiryDate, so pass every attribute that should be kept:
// an omitted ExpiryDate clears the expiry date and a false Enabled disables
// the token. Updating a token of another user is rejected with 403 Forbidden.
// A token carrying the name the API reserves for internal use cannot be
// updated at all, neither renamed nor disabled, the API answers 400 Bad
// Request; delete it instead. The Token attribute is stripped from the
// request body.
func (api *API) UpdateApiTokenContext(ctx context.Context, token *APIToken) (*APIToken, error) {
	if _, ok := api.methods["updateApiToken"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateApiToken")
	}

	definition := api.methods["updateApiToken"]
	definition.Action = fmt.Sprintf(definition.Action, token.ID)

	// The API never reads the token string, keep the credential out of the request.
	payload := *token
	payload.Token = ""

	result, err := api.call(ctx, definition, &payload)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*APIToken)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteApiTokenContext deletes the passed API token using the MYRA API.
// Only the ID is needed. Deleting a token of another user is rejected with
// 403 Forbidden. The Token attribute is stripped from the request body.
func (api *API) DeleteApiTokenContext(ctx context.Context, token *APIToken) (*APIToken, error) {
	if _, ok := api.methods["deleteApiToken"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteApiToken")
	}

	definition := api.methods["deleteApiToken"]
	definition.Action = fmt.Sprintf(definition.Action, token.ID)

	// The API never reads the token string, keep the credential out of the request.
	payload := *token
	payload.Token = ""

	_, err := api.call(ctx, definition, &payload)
	if err != nil {
		return nil, err
	}

	return token, nil
}
