package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// getTagWAFRuleMethods returns Tag WAF rule related API calls
//
func getTagWAFRuleMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getTagWAFRule": {
			Name:               "getTagWAFRule",
			Action:             "tag/%d/waf-rules/%d",
			Method:             http.MethodGet,
			Result:             TagWAFRule{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listTagWAFRules": {
			Name:   "listTagWAFRules",
			Action: "tag/%d/waf-rules",
			Method: http.MethodGet,
			Result: []TagWAFRule{},
		},
		"createTagWAFRule": {
			Name:   "createTagWAFRule",
			Action: "tag/%d/waf-rules",
			Method: http.MethodPost,
			Result: TagWAFRule{},
		},
		"updateTagWAFRule": {
			Name:   "updateTagWAFRule",
			Action: "tag/%d/waf-rules/%d",
			Method: http.MethodPut,
			Result: TagWAFRule{},
		},
		"deleteTagWAFRule": {
			Name:   "deleteTagWAFRule",
			Action: "tag/%d/waf-rules/%d",
			Method: http.MethodDelete,
			Result: TagWAFRule{},
		},
	}
}

//
// TagWAFRule ...
//
type TagWAFRule struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	ExpireDate    *types.DateTime `json:"expireDate,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Direction     string          `json:"direction"`
	LogIdentifier string          `json:"logIdentifier"`
	Sort          int             `json:"sort"`
	Sync          bool            `json:"sync"`
	ProcessNext   bool            `json:"processNext"`
	Enabled       bool            `json:"enabled"`
	Actions       []*WAFAction    `json:"actions"`
	Conditions    []*WAFCondition `json:"conditions"`
	TagId         int             `json:"tagId"`
}

//
// GetTagWAFRule returns a single tag for the given identifier
//
func (api *API) GetTagWAFRule(tagId int, ruleId int) (*TagWAFRule, error) {
	if _, ok := methods["getTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getTagWAFRule")
	}

	definition := methods["getTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, ruleId)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*TagWAFRule), nil
}

//
// ListTagWAFRules returns a slice containing all visible tags
//
func (api *API) ListTagWAFRules(tagId int, params map[string]string) ([]TagWAFRule, error) {
	if _, ok := methods["listTagWAFRules"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagWAFRules")
	}

	definition := methods["listTagWAFRules"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var tags []TagWAFRule
	tags = append(tags, *result.(*[]TagWAFRule)...)

	return tags, nil
}

//
// CreateTagWAFRule creates a new tag using the MYRA API
//
func (api *API) CreateTagWAFRule(rule *TagWAFRule, tagId int) (*TagWAFRule, error) {
	if _, ok := methods["createTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagWAFRule")
	}

	definition := methods["createTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}

	return result.(*TagWAFRule), nil
}

//
// UpdateTagWAFRule updates the passed tag using the MYRA API
//
func (api *API) UpdateTagWAFRule(rule *TagWAFRule) (*TagWAFRule, error) {
	if _, ok := methods["updateTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagWAFRule")
	}

	definition := methods["updateTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, rule.TagId, rule.ID)

	result, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}

	return result.(*TagWAFRule), nil
}

//
// DeleteTagWAFRule deletes the passed tag using the MYRA API
//
func (api *API) DeleteTagWAFRule(rule *TagWAFRule) (*TagWAFRule, error) {
	if _, ok := methods["deleteTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTagWAFRule")
	}

	definition := methods["deleteTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, rule.TagId, rule.ID)

	_, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}

	return rule, nil
}
