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
	ID             int             `json:"id,omitempty"`
	Created        *types.DateTime `json:"created,omitempty"`
	Modified       *types.DateTime `json:"modified,omitempty"`
	Name           string          `json:"name"`
	VhostId        int             `json:"vhostId"`
	SubDomainName  string          `json:"subDomainName"`
	MaxConcurrent  int             `json:"maxConcurrent"`
	SessionTimeout int             `json:"sessionTimeout"`
	WaitRefresh    int             `json:"waitRefresh"`
	Paths          []string        `json:"paths"`
	Content        string          `json:"content"`
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
	var records []WaitingRoom
	records = append(records, *result.(*[]WaitingRoom)...)

	return records, nil
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
	var records []WaitingRoom
	records = append(records, *result.(*[]WaitingRoom)...)

	return records, nil
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
