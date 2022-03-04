package myrasec

import (
	"fmt"
	"strconv"

	"github.com/Myra-Security-GmbH/myrasec-go/pkg/types"
)

//
// IPRange ...
//
type IPRange struct {
	ID        int             `json:"id,omitempty"`
	Created   *types.DateTime `json:"created,omitempty"`
	Modified  *types.DateTime `json:"modified,omitempty"`
	ValidFrom *types.DateTime `json:"validFrom,omitempty"`
	ValidTo   *types.DateTime `json:"validTo,omitempty"`
	Network   string          `json:"network"`
	Comment   string          `json:"comment,omitempty"`
	Enabled   bool            `json:"enabled,omitempty"`
}

//
// ListIPRanges returns a slice containing all ip ranges
//
func (api *API) ListIPRanges(params map[string]string) ([]IPRange, error) {
	if _, ok := methods["listIPRanges"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listIPRanges")
	}

	page := 1
	var err error
	if pageParam, ok := params[ParamPage]; ok {
		delete(params, ParamPage)
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			page = 1
		}
	}

	definition := methods["listIPRanges"]
	definition.Action = fmt.Sprintf(definition.Action, page)

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
