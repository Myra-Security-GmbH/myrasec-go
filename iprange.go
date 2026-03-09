package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getIPRangeMethods returns IP range related API calls
func getIPRangeMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listIPRanges": {
			Name:   "listIPRanges",
			Action: "ip-ranges",
			Method: http.MethodGet,
			Result: []IPRange{},
		},
	}
}

// IPRange represents a specific network block defined by a CIDR notation.
// It is used to define validity periods and access status for specific network ranges.
type IPRange struct {
	// ID is the unique identifier for the IP range.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the IP range. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the IP range was added.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Network defines the IP address range in CIDR notation.
	Network string `json:"network" jsonschema:"The network address in CIDR notation (e.g., '192.168.0.0/24' or '2001:db8::/32')."`

	// ValidFrom specifies the start date for the range's validity.
	// If nil, the range is considered valid immediately upon creation.
	ValidFrom *types.DateTime `json:"validFrom,omitempty" jsonschema:"The timestamp (ISO 8601) when this range becomes valid. If null, it is valid immediately."`

	// ValidTo specifies the expiration date for the range.
	// If nil, the range remains valid indefinitely.
	ValidTo *types.DateTime `json:"validTo,omitempty" jsonschema:"The timestamp (ISO 8601) when this range expires. If null, it remains valid indefinitely."`

	// Enabled controls whether the IP range is currently active.
	Enabled bool `json:"enabled" jsonschema:"Indicates if this IP range is currently active (enabled) or ignored."`

	// Comment provides a descriptive note for the IP range.
	Comment string `json:"comment,omitempty" jsonschema:"A descriptive comment or note for this IP range."`
}

// ListIPRanges returns a slice containing all ip ranges
func (api *API) ListIPRanges(params map[string]string) ([]IPRange, error) {
	if _, ok := methods["listIPRanges"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listIPRanges")
	}

	definition := methods["listIPRanges"]

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]IPRange)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}
