package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getAPIKeyMethods returns API key related API calls.
func getAPIKeyMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listApiKeys": {
			Name:   "listApiKeys",
			Action: "user/%d/api-keys",
			Method: http.MethodGet,
			Result: []APIKey{},
		},
		"createApiKey": {
			Name:   "createApiKey",
			Action: "user/%d/api-keys",
			Method: http.MethodPost,
			Result: APIKey{},
		},
		"deleteApiKey": {
			Name:   "deleteApiKey",
			Action: "user/%d/api-keys/%d",
			Method: http.MethodDelete,
			Result: APIKey{},
		},
	}
}

// APIKey represents an authentication credential used to access the API.
// It encapsulates the public identifier, the private secret and metadata
// regarding the lifecycle of the key.
type APIKey struct {
	// ID is the unique identifier for the API key.
	// This value is server-generated and cannot be set during creation.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the API key. This value is server-generated and cannot be set during creation."`

	// Created indicates the timestamp when the API key was generated.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp indicating when the API key was generated (ISO 8601 format). This is a server-managed, read-only value."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format), serving as a version identifier. This field is required for updates and deletes to ensure data consistency."`

	// Name is an arbitrary, user-defined label for the API key.
	Name string `json:"name,omitempty" jsonschema:"An arbitrary, user-defined label for the API key."`

	// Key is the public token part of the credential.
	Key string `json:"key,omitempty" jsonschema:"The public token string of the API key."`

	// Secret is the private portion of the credential.
	// Note: This value is returned only once upon creation and cannot be retrieved later.
	Secret string `json:"secret,omitempty" jsonschema:"The private secret of the API key. Visible only once upon creation; it cannot be retrieved later."`
}

// ListApiKeys returns a slice containing all available API keys
func (api *API) ListApiKeys(params map[string]string) ([]APIKey, error) {
	if _, ok := api.methods["listApiKeys"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listApiKeys")
	}

	user, err := api.Me()
	if err != nil {
		return nil, err
	}

	definition := api.methods["listApiKeys"]
	definition.Action = fmt.Sprintf(definition.Action, user.ID)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]APIKey)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateApiKey creates a new API key using the MYRA API
func (api *API) CreateApiKey(apikey *APIKey) (*APIKey, error) {
	if _, ok := api.methods["createApiKey"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createApiKey")
	}

	user, err := api.Me()
	if err != nil {
		return nil, err
	}

	definition := api.methods["createApiKey"]
	definition.Action = fmt.Sprintf(definition.Action, user.ID)

	result, err := api.call(definition, apikey)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*APIKey)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteApiKey deletes the passed API key using the MYRA API
func (api *API) DeleteApiKey(apikey *APIKey) (*APIKey, error) {
	if _, ok := api.methods["deleteApiKey"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteApiKey")
	}

	user, err := api.Me()
	if err != nil {
		return nil, err
	}

	definition := api.methods["deleteApiKey"]
	definition.Action = fmt.Sprintf(definition.Action, user.ID, apikey.ID)

	_, err = api.call(definition, apikey)
	if err != nil {
		return nil, err
	}

	return apikey, nil
}
