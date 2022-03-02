package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// IPRange ...
//
type IPRange struct {
	ID        int             `json:"id,omitempty"`
	Created   *types.DateTime `json:"created,omitempty"`
	Modified  *types.DateTime `json:"modified,omitempty"`
	Network   string          `json:"network"`
	ValidFrom *types.DateTime `json:"validFrom,omitempty"`
	ValidTo   *types.DateTime `json:"validTo,omitempty"`
	Enabled   bool            `json:"enabled,omitempty"`
	Comment   string          `json:"comment,omitempty"`
}

//
// ListIPRanges returns a slice containing all ip ranges
//
func (api *API) ListIPRanges(params map[string]string) ([]IPRange, error) {
	if _, ok := methods["listIPRanges"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listIPRanges")
	}

	definition := methods["listIPRanges"]

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []IPRange
	for _, v := range *result.(*[]IPRange) {
		records = append(records, v)
	}

	return records, nil
}
