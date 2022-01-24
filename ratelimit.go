package myrasec

import (
	"fmt"
	"strconv"

	"github.com/Myra-Security-GmbH/myrasec-go/pkg/types"
)

//
// RateLimit ...
//
type RateLimit struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	Burst         int             `json:"burst"`
	Network       string          `json:"network"`
	SubDomainName string          `json:"subDomainName"`
	Timeframe     int             `json:"timeframe"`
	Type          string          `json:"type"`
	Value         int             `json:"value"`
}

//
// ListRateLimits returns a slice containing all visible rate limit settings
// Valid rateLimitType values are "dns" or "tag"
//
// Rate limit settings can be filtered using the params map
//
// Avalilable filters/query parameters:
//		search (string) - filter by the specified search query
// Additional valid filters/query parameters for ruleType = "dns":
//		subDomainName (string) - filter rate limit settings for this subdomain (name)
//		reference (int) - filter rate limit settings for this domain (ID)
// Additional valid filters/query parameters for ruleType = "tag":
//		reference (int) - filter rate limit settings for this tag (ID)
//
func (api *API) ListRateLimits(rateLimitType string, params map[string]string) ([]RateLimit, error) {
	if _, ok := methods["listRateLimits"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listRateLimits")
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

	definition := methods["listRateLimits"]
	definition.Action = fmt.Sprintf(definition.Action, rateLimitType, page)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}
	var records []RateLimit
	for _, v := range *result.(*[]RateLimit) {
		records = append(records, v)
	}

	return records, nil
}

//
// CreateRateLimit creates a new rate limit setting for the passed subdomain (name) using the MYRA API
//
func (api *API) CreateRateLimit(ratelimit *RateLimit) (*RateLimit, error) {
	if _, ok := methods["createRateLimit"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "createRateLimit")
	}

	result, err := api.call(methods["createRateLimit"], ratelimit)
	if err != nil {
		return nil, err
	}
	return result.(*RateLimit), nil
}

//
// UpdateRateLimit updates the passed rate limit setting using the MYRA API
//
func (api *API) UpdateRateLimit(ratelimit *RateLimit) (*RateLimit, error) {
	if _, ok := methods["updateRateLimit"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "updateRateLimit")
	}

	result, err := api.call(methods["updateRateLimit"], ratelimit)
	if err != nil {
		return nil, err
	}
	return result.(*RateLimit), nil
}

//
// DeleteRateLimit deletes the passed rate limit setting using the MYRA API
//
func (api *API) DeleteRateLimit(ratelimit *RateLimit) (*RateLimit, error) {
	if _, ok := methods["deleteRateLimit"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteRateLimit")
	}

	result, err := api.call(methods["deleteRateLimit"], ratelimit)
	if err != nil {
		return nil, err
	}
	return result.(*RateLimit), nil
}
