package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

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

type APIKey struct {
	ID       int             `json:"id,omitempty"`
	Created  *types.DateTime `json:"created,omitempty"`
	Modified *types.DateTime `json:"modified,omitempty"`
	Name     string          `json:"name,omitempty"`
	Key      string          `json:"key,omitempty"`
	Secret   string          `json:"secret,omitempty"`
}

// ListApiKeys returns a slice containing all available API keys
func (api *API) ListApiKeys(params map[string]string) ([]APIKey, error) {
	if _, ok := methods["listApiKeys"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listApiKeys")
	}

	user, err := api.Me()
	if err != nil {
		return nil, err
	}

	definition := methods["listApiKeys"]
	definition.Action = fmt.Sprintf(definition.Action, user.ID)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]APIKey), nil
}

// CreateApiKey creates a new API key using the MYRA API
func (api *API) CreateApiKey(apikey *APIKey) (*APIKey, error) {
	if _, ok := methods["createApiKey"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createApiKey")
	}

	user, err := api.Me()
	if err != nil {
		return nil, err
	}

	definition := methods["createApiKey"]
	definition.Action = fmt.Sprintf(definition.Action, user.ID)

	result, err := api.call(definition, apikey)
	if err != nil {
		return nil, err
	}

	return result.(*APIKey), nil
}

// DeleteApiKey deletes the passed API key using the MYRA API
func (api *API) DeleteApiKey(apikey *APIKey) (*APIKey, error) {
	if _, ok := methods["deleteApiKey"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteApiKey")
	}

	user, err := api.Me()
	if err != nil {
		return nil, err
	}

	definition := methods["deleteApiKey"]
	definition.Action = fmt.Sprintf(definition.Action, user.ID, apikey.ID)

	_, err = api.call(definition, apikey)
	if err != nil {
		return nil, err
	}

	return apikey, nil
}
