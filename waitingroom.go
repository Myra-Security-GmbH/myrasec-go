package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getWaitingRoomMethods returns WaitingRoom related API calls
func getWaitingRoomMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listWaitingRoomsForDomain": {
			Name:   "listWaitingRoomsForDomain",
			Action: "waiting-rooms?domainId=%d",
			Method: http.MethodGet,
			Result: []WaitingRoom{},
		},
		"listWaitingRoomsForSubDomain": {
			Name:   "listWaitingRoomsForSubDomain",
			Action: "waiting-rooms?subDomainName=%s",
			Method: http.MethodGet,
			Result: []WaitingRoom{},
		},
		"createWaitingRoom": {
			Name:   "createWaitingRoom",
			Action: "waiting-room",
			Method: http.MethodPost,
			Result: WaitingRoom{},
		},
		"updateWaitingRoom": {
			Name:   "updateWaitingRoom",
			Action: "waiting-room/%d",
			Method: http.MethodPut,
			Result: WaitingRoom{},
		},
		"deleteWaitingRoom": {
			Name:   "deleteWaitingRoom",
			Action: "waiting-room/%d",
			Method: http.MethodDelete,
			Result: WaitingRoom{},
		},
		"getWaitingRoom": {
			Name:               "getWaitingRoom",
			Action:             "waiting-room/%d",
			Method:             http.MethodGet,
			Result:             WaitingRoom{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
	}
}

// WaitingRoom ...
type WaitingRoom struct {
	ID             int             `json:"id,omitempty" jsonschema:"ID is a unique iDentifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a Waiting Room it is necessary to add this attribute to your object."`
	Created        *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be set by the server after creating a new cache setting object. This value is only informational, so it is not necessary to add this as an attribute to any API call."`
	Modified       *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always need to add the modified timestamp for updates and deletes. This value is always a date type with an ISO 8601 format."`
	Name           string          `json:"name" jsonschema:"Name of the Waiting Room."`
	VhostId        int             `json:"vhostId" jsonschema:"ID of the VHost for the Waiting Room."`
	SubDomainName  string          `json:"subDomainName" jsonschema:"Identifies the subdomain via a FQDN (Full Qualified Domain Name) that the Waiting Room belongs to. This value is optional and is determined from the VHost based on its ID."`
	MaxConcurrent  int             `json:"maxConcurrent" jsonschema:"The maximum number of visitors allowed on the Origin server at the same time. As soon as this value is exceeded, each additional visitor is directed to the Waiting Room."`
	SessionTimeout int             `json:"sessionTimeout" jsonschema:"Defines the duration in seconds during which an inactive session may access the Origin server. If the same session does not access the server again during this time, access for that session will be disabled."`
	WaitRefresh    int             `json:"waitRefresh" jsonschema:"Defines the duration in seconds after which the waiting page is reloaded. If the session is not accessed again after the third reload, the session will be removed from the queue."`
	Paths          []string        `json:"paths" jsonschema:"Defines a specific path within the apex domain or subdomain for which the waiting room is to be valid. The path needs to be defined as a regular expression. The default value in the PATH field is . . If the default value . is used as the path, the waiting pages and settings of all waiting rooms with a specific path of the corresponding apex domain or subdomain are overwritten."`
	Content        string          `json:"content" jsonschema:"The HTML content of the Waiting Room."`
}

// ListWaitingRoomsForDomain returns a slice containing all visible waiting rooms for domain
func (api *API) ListWaitingRoomsForDomain(domainId int, params map[string]string) ([]WaitingRoom, error) {
	if _, ok := methods["listWaitingRoomsForDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWaitingRoomsForDomain")
	}

	definition := methods["listWaitingRoomsForDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]WaitingRoom), nil
}

// ListWaitingRoomsForSubDomain returns a slice containing all visible waiting rooms for subdomain
func (api *API) ListWaitingRoomsForSubDomain(subDomainName string, params map[string]string) ([]WaitingRoom, error) {
	if _, ok := methods["listWaitingRoomsForSubDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWaitingRoomsForSubDomain")
	}

	definition := methods["listWaitingRoomsForSubDomain"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]WaitingRoom), nil
}

// GetWaitingRoom returns the waiting room
func (api *API) GetWaitingRoom(id int) (*WaitingRoom, error) {
	if _, ok := methods["getWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getWaitingRoom")
	}

	definition := methods["getWaitingRoom"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}
	return result.(*WaitingRoom), nil
}

// CreateWaitingRoom creates a new waiting room
func (api *API) CreateWaitingRoom(waitingroom *WaitingRoom) (*WaitingRoom, error) {
	if _, ok := methods["createWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createWaitingRoom")
	}

	definition := methods["createWaitingRoom"]

	result, err := api.call(definition, waitingroom)
	if err != nil {
		return nil, err
	}
	return result.(*WaitingRoom), nil
}

// UpdateWaitingRoom updates the waiting room
func (api *API) UpdateWaitingRoom(waitingroom *WaitingRoom) (*WaitingRoom, error) {
	if _, ok := methods["updateWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateWaitingRoom")
	}

	definition := methods["updateWaitingRoom"]
	definition.Action = fmt.Sprintf(definition.Action, waitingroom.ID)

	result, err := api.call(definition, waitingroom)
	if err != nil {
		return nil, err
	}
	return result.(*WaitingRoom), nil
}

// DeleteWaitingRoom updates the waiting room
func (api *API) DeleteWaitingRoom(waitingroom *WaitingRoom) (*WaitingRoom, error) {
	if _, ok := methods["deleteWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteWaitingRoom")
	}

	definition := methods["deleteWaitingRoom"]
	definition.Action = fmt.Sprintf(definition.Action, waitingroom.ID)

	_, err := api.call(definition, waitingroom)
	if err != nil {
		return nil, err
	}
	return waitingroom, nil
}
