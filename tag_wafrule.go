package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getTagWAFRuleMethods returns Tag WAF rule related API calls
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

// TagWAFRule ...
type TagWAFRule struct {
	ID            int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a TagWAFRule it is necessary to add this attribute to your object."`
	Created       *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new TagWAFRule object. This value is informational so it is not necessary to add this attribute to any API call."`
	Modified      *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format."`
	ExpireDate    *types.DateTime `json:"expireDate,omitempty" jsonschema:"ExpireDate describes how long a TagWAFRule is valid, and when it will expire."`
	Name          string          `json:"name" jsonschema:"Identifies the TagWAFRule by its name."`
	Description   string          `json:"description" jsonschema:"The Description will explain what the TagWAFRule is for."`
	Direction     string          `json:"direction" jsonschema:"The direction can be in (for Request) or out (for Response)."`
	LogIdentifier string          `json:"logIdentifier" jsonschema:"A string to identify the matching rule in the access log."`
	Uuid          string          `json:"uuid,omitempty"`
	Sort          int             `json:"sort" jsonschema:"Defines the sorting of TagWAFRules."`
	Sync          bool            `json:"sync"`
	ProcessNext   bool            `json:"processNext" jsonschema:"After a rule has been applied, the rule chain will be executed as determined."`
	Enabled       bool            `json:"enabled" jsonschema:"Describes if the rule is enabled or not."`
	Actions       []*WAFAction    `json:"actions" jsonschema:"List of WAF actions."`
	Conditions    []*WAFCondition `json:"conditions" jsonschema:"List of WAF conditions."`
	TagId         int             `json:"tagId" jsonschema:"The related TagId."`
}

// GetTagWAFRule returns a single tag for the given identifier
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

// ListTagWAFRules returns a slice containing all visible tags
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

	return *result.(*[]TagWAFRule), nil
}

// CreateTagWAFRule creates a new tag using the MYRA API
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

// UpdateTagWAFRule updates the passed tag using the MYRA API
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

// DeleteTagWAFRule deletes the passed tag using the MYRA API
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
