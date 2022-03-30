package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// RateLimit ...
//
type RateLimit struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	Network       string          `json:"network"`
	SubDomainName string          `json:"subDomainName"`
	Type          string          `json:"type"`
	Burst         int             `json:"burst"`
	Timeframe     int             `json:"timeframe"`
	Value         int             `json:"value"`
}

//
// ListRateLimits returns a slice containing all visible rate limit settings
//
func (api *API) ListRateLimits(domainId int, subDomainName string, params map[string]string) ([]RateLimit, error) {
	if _, ok := methods["listRateLimits"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listRateLimits")
	}

	definition := methods["listRateLimits"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}
	var records []RateLimit
	records = append(records, *result.(*[]RateLimit)...)

	return records, nil
}

//
// CreateRateLimit creates a new rate limit setting for the passed subdomain (name) using the MYRA API
//
func (api *API) CreateRateLimit(ratelimit *RateLimit, domainId int, subDomainName string) (*RateLimit, error) {
	if _, ok := methods["createRateLimit"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createRateLimit")
	}

	definition := methods["createRateLimit"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, ratelimit)
	if err != nil {
		return nil, err
	}
	return result.(*RateLimit), nil
}

//
// UpdateRateLimit updates the passed rate limit setting using the MYRA API
//
func (api *API) UpdateRateLimit(ratelimit *RateLimit, domainId int, subDomainName string) (*RateLimit, error) {
	if _, ok := methods["updateRateLimit"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateRateLimit")
	}

	definition := methods["updateRateLimit"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, ratelimit.ID)

	result, err := api.call(definition, ratelimit)
	if err != nil {
		return nil, err
	}
	return result.(*RateLimit), nil
}

//
// DeleteRateLimit deletes the passed rate limit setting using the MYRA API
//
func (api *API) DeleteRateLimit(ratelimit *RateLimit, domainId int, subDomainName string) (*RateLimit, error) {
	if _, ok := methods["deleteRateLimit"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteRateLimit")
	}

	definition := methods["deleteRateLimit"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, ratelimit.ID)

	_, err := api.call(definition, ratelimit)
	if err != nil {
		return nil, err
	}
	return ratelimit, nil
}
