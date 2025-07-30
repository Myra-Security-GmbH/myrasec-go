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
	ID       int             `json:"id,omitempty"`
	Created  *types.DateTime `json:"created,omitempty"`
	Modified *types.DateTime `json:"modified,omitempty"`
	Login    string          `json:"login,omitempty"`
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
