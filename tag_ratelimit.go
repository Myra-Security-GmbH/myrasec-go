package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// getTagRateLimitMethods returns Tag Rate Limit related API calls
//
func getTagRateLimitMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listTagRateLimits": {
			Name:   "listTagRateLimits",
			Action: "tag/%d/ratelimits",
			Method: http.MethodGet,
			Result: []TagRateLimit{},
		},
		"createTagRateLimit": {
			Name:   "createTagRateLimit",
			Action: "tag/%d/ratelimits",
			Method: http.MethodPost,
			Result: TagRateLimit{},
		},
		"updateTagRateLimit": {
			Name:   "updateTagRateLimit",
			Action: "tag/%d/ratelimits/%d",
			Method: http.MethodPut,
			Result: TagRateLimit{},
		},
		"deleteTagRateLimit": {
			Name:   "deleteTagRateLimit",
			Action: "tag/%d/ratelimits/%d",
			Method: http.MethodDelete,
			Result: TagRateLimit{},
		},
	}
}

//
// TagRateLimit ...
//
type TagRateLimit struct {
	ID        int             `json:"id,omitempty"`
	Created   *types.DateTime `json:"created,omitempty"`
	Modified  *types.DateTime `json:"modified,omitempty"`
	Network   string          `json:"network"`
	TagId     int             `json:"tagId"`
	Type      string          `json:"type"`
	Burst     int             `json:"burst"`
	Timeframe int             `json:"timeframe"`
	Value     int             `json:"value"`
}

//
// ListTagRateLimits returns a slice containing all visible rate limits for a tag
//
func (api *API) ListTagRateLimits(tagId int, params map[string]string) ([]TagRateLimit, error) {
	if _, ok := methods["listTagRateLimits"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagRateLimits")
	}

	definition := methods["listTagRateLimits"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var rateLimits []TagRateLimit
	rateLimits = append(rateLimits, *result.(*[]TagRateLimit)...)

	return rateLimits, nil
}

//
// CreateTagRateLimit creates a new rate limit for passed tag
//
func (api *API) CreateTagRateLimit(rateLimit *TagRateLimit, tagId int) (*TagRateLimit, error) {
	if _, ok := methods["createTagRateLimit"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagRateLimit")
	}

	definition := methods["createTagRateLimit"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, rateLimit)
	if err != nil {
		return nil, err
	}
	return result.(*TagRateLimit), nil
}

//
// UpdateTagRateLimit updates a new rate limit for passed tag
//
func (api *API) UpdateTagRateLimit(rateLimit *TagRateLimit, tagId int) (*TagRateLimit, error) {
	if _, ok := methods["updateTagRateLimit"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagRateLimit")
	}

	definition := methods["updateTagRateLimit"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, rateLimit.ID)

	result, err := api.call(definition, rateLimit)
	if err != nil {
		return nil, err
	}
	return result.(*TagRateLimit), nil
}

//
// DeleteTagRateLimit deletes a new rate limit for passed tag
//
func (api *API) DeleteTagRateLimit(rateLimit *TagRateLimit, tagId int) (*TagRateLimit, error) {
	if _, ok := methods["deleteTagRateLimit"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTagRateLimit")
	}

	definition := methods["deleteTagRateLimit"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, rateLimit.ID)

	_, err := api.call(definition)
	if err != nil {
		return nil, err
	}
	return rateLimit, nil
}
