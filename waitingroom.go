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

// WaitingRoom represents a virtual queue system for high-traffic scenarios.
// It limits the number of concurrent users allowed on the origin server to prevent overloads.
type WaitingRoom struct {
	// ID is the unique identifier for the waiting room configuration.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the waiting room. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the waiting room was configured.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Name is a descriptive label for the waiting room.
	Name string `json:"name" jsonschema:"A descriptive name for this waiting room configuration."`

	// VhostId is the identifier of the Virtual Host this waiting room protects.
	VhostId int `json:"vhostId" jsonschema:"The unique identifier of the target VHost (Virtual Host)."`

	// SubDomainName is the FQDN associated with the VHost.
	// Optional: If omitted, it is determined automatically from the VhostId.
	SubDomainName string `json:"subDomainName" jsonschema:"The FQDN (Fully Qualified Domain Name) the waiting room belongs to. Optional; if not provided, it is inferred from the 'vhostId'."`

	// MaxConcurrent sets the limit of simultaneous users allowed on the origin.
	// Exceeding this limit triggers the waiting room for new visitors.
	MaxConcurrent int `json:"maxConcurrent" jsonschema:"The maximum number of concurrent active users allowed on the origin. Once exceeded, new visitors are redirected to the waiting room."`

	// SessionTimeout defines the idle timeout for active users.
	// If a user is inactive for this period, they lose their spot.
	SessionTimeout int `json:"sessionTimeout" jsonschema:"Idle timeout in seconds. If an active session does not access the server within this time, access is revoked."`

	// WaitRefresh defines the auto-reload interval for the waiting page.
	// Logic: If the session is not accessed after the 3rd reload, it is removed from the queue.
	WaitRefresh int `json:"waitRefresh" jsonschema:"The auto-reload interval in seconds for the waiting page. Note: If the client does not poll/reload, the session is removed from the queue after 3 intervals."`

	// Paths defines the URL patterns covered by this waiting room.
	// Expects Regex. Default is "." (match all).
	// Warning: "." overrides specific path settings on the same domain.
	Paths []string `json:"paths" jsonschema:"List of URL paths (Regular Expressions) to protect. Default: '.' (matches everything). Warning: Using '.' overwrites/takes precedence over other specific path rules for this domain."`

	// Content contains the HTML code displayed to users in the queue.
	Content string `json:"content" jsonschema:"The raw HTML content displayed to visitors while they are in the waiting queue."`
}

// ListWaitingRoomsForDomain returns a slice containing all visible waiting rooms for domain
func (api *API) ListWaitingRoomsForDomain(domainId int, params map[string]string) ([]WaitingRoom, error) {
	if _, ok := api.methods["listWaitingRoomsForDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWaitingRoomsForDomain")
	}

	definition := api.methods["listWaitingRoomsForDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]WaitingRoom)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListWaitingRoomsForSubDomain returns a slice containing all visible waiting rooms for subdomain
func (api *API) ListWaitingRoomsForSubDomain(subDomainName string, params map[string]string) ([]WaitingRoom, error) {
	if _, ok := api.methods["listWaitingRoomsForSubDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWaitingRoomsForSubDomain")
	}

	definition := api.methods["listWaitingRoomsForSubDomain"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]WaitingRoom)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// GetWaitingRoom returns the waiting room
func (api *API) GetWaitingRoom(id int) (*WaitingRoom, error) {
	if _, ok := api.methods["getWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getWaitingRoom")
	}

	definition := api.methods["getWaitingRoom"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}
	res, ok := result.(*WaitingRoom)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateWaitingRoom creates a new waiting room
func (api *API) CreateWaitingRoom(waitingroom *WaitingRoom) (*WaitingRoom, error) {
	if _, ok := api.methods["createWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createWaitingRoom")
	}

	definition := api.methods["createWaitingRoom"]

	result, err := api.call(definition, waitingroom)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*WaitingRoom)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateWaitingRoom updates the waiting room
func (api *API) UpdateWaitingRoom(waitingroom *WaitingRoom) (*WaitingRoom, error) {
	if _, ok := api.methods["updateWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateWaitingRoom")
	}

	definition := api.methods["updateWaitingRoom"]
	definition.Action = fmt.Sprintf(definition.Action, waitingroom.ID)

	result, err := api.call(definition, waitingroom)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*WaitingRoom)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteWaitingRoom deletes the passed waiting room using the MYRA API
func (api *API) DeleteWaitingRoom(waitingroom *WaitingRoom) (*WaitingRoom, error) {
	if _, ok := api.methods["deleteWaitingRoom"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteWaitingRoom")
	}

	definition := api.methods["deleteWaitingRoom"]
	definition.Action = fmt.Sprintf(definition.Action, waitingroom.ID)

	_, err := api.call(definition, waitingroom)
	if err != nil {
		return nil, err
	}
	return waitingroom, nil
}
