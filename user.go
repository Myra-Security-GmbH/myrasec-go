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

type User struct {
	ID       int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type."`
	Created  *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new User object."`
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format."`
	Login    string          `json:"login,omitempty" jsonschema:"The login name that is the same as the email address."`
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

	return result.(*User), nil
}
