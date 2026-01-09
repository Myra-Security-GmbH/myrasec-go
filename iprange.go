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

// IPRange ...
type IPRange struct {
	ID        int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object."`
	Created   *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format."`
	Modified  *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. This value is always a date type with an ISO 8601 format."`
	Network   string          `json:"network" jsonschema:"The network (CIDR notation)."`
	ValidFrom *types.DateTime `json:"validFrom,omitempty"`
	ValidTo   *types.DateTime `json:"validTo,omitempty"`
	Enabled   bool            `json:"enabled"`
	Comment   string          `json:"comment,omitempty"`
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

	return *result.(*[]IPRange), nil
}
