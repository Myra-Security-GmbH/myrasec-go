package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

func getUserMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"me": {
			Name:               "me",
			Action:             "user/me",
			Method:             http.MethodGet,
			Result:             User{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
	}
}

// User represents a registered account holder in the system.
// It contains the authentication identity (Login) and metadata.
type User struct {
	// ID is the unique identifier for the user.
	// This value is server-generated and read-only.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the user. Server-generated and read-only. Ignored during creation."`

	// Created indicates when the user was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Login is the unique username for the account.
	// This must be a valid email address.
	Login string `json:"login,omitempty" jsonschema:"The user's login name. Must be a valid email address (format: user@example.com)."`
}

// Me returns the active user information
func (api *API) Me() (*User, error) {
	if _, ok := methods["me"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "me")
	}

	definition := methods["me"]
	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*User)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}
