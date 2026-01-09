package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getWAFMethods returns WAF related API calls
func getWAFMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listWAFConditions": {
			Name:   "listWAFConditions",
			Action: "waf/conditions",
			Method: http.MethodGet,
			Result: []WAFCondition{},
		},
		"listWAFActions": {
			Name:   "listWAFActions",
			Action: "waf/actions",
			Method: http.MethodGet,
			Result: []WAFAction{},
		},
		"listWAFRules": {
			Name:   "listWAFRules",
			Action: "domain/%d/waf-rules",
			Method: http.MethodGet,
			Result: []WAFRule{},
		},
		"fetchWAFRule": {
			Name:   "fetchWAFRule",
			Action: "domain/waf-rules/%d",
			Method: http.MethodGet,
			Result: []WAFRule{},
		},
		"createWAFRule": {
			Name:   "createWAFRule",
			Action: "domain/%d/%s/waf-rules",
			Method: http.MethodPost,
			Result: WAFRule{},
		},
		"updateWAFRule": {
			Name:   "updateWAFRule",
			Action: "domain/%d/%s/waf-rules/%d",
			Method: http.MethodPut,
			Result: WAFRule{},
		},
		"deleteWAFRule": {
			Name:   "deleteWAFRule",
			Action: "domain/waf-rules/%d",
			Method: http.MethodDelete,
			Result: WAFRule{},
		},
	}
}

// WAFRule ...
type WAFRule struct {
	ID            int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a WAFRule it is necessary to add this attribute to your object."`
	Created       *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new WAFRule object. This value is informational so it is not necessary to add this attribute to any API call."`
	Modified      *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format."`
	ExpireDate    *types.DateTime `json:"expireDate,omitempty" jsonschema:"ExpireDate describes how long a WAFRule is valid, and when it will expire."`
	Name          string          `json:"name" jsonschema:"Identifies the WAF rule by its name."`
	Description   string          `json:"description" jsonschema:"The Description will explain what the WAFRule is for."`
	Direction     string          `json:"direction" jsonschema:"The direction can be in (for Request) or out (for Response)."`
	LogIdentifier string          `json:"logIdentifier" jsonschema:"A string to identify the matching rule in the access log."`
	Uuid          string          `json:"uuid,omitempty"`
	RuleType      string          `json:"ruleType"`
	SubDomainName string          `json:"subDomainName"`
	Sort          int             `json:"sort" jsonschema:"Defines the sorting of WAFRules."`
	Sync          bool            `json:"sync"`
	Template      bool            `json:"template"`
	ProcessNext   bool            `json:"processNext" jsonschema:"After a rule has been applied, the rule chain will be executed as determined."`
	Enabled       bool            `json:"enabled" jsonschema:"Describes if the rule is enabled or not."`
	Actions       []*WAFAction    `json:"actions" jsonschema:"List of WAF actions."`
	Conditions    []*WAFCondition `json:"conditions" jsonschema:"List of WAF conditions."`
}

// WAFAction ...
type WAFAction struct {
	ID                int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a WAF Action it is necessary to add this attribute to your object."`
	Created           *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new WAFRule action object. This value is informational so it is not necessary to add this attribute to any API call."`
	Modified          *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format. "`
	ForceCustomValues bool            `json:"forceCustomValues" jsonschema:"This attributes determines number of input fields when utilised (0=none, 1=value, 2=key+value)."`
	AvailablePhases   int             `json:"availablePhases" jsonschema:"This attributes determines the support for different phases (1=request, 2=response, 3=both)."`
	Name              string          `json:"name" jsonschema:"Display name of the action."`
	Type              string          `json:"type" jsonschema:"Type of the action."`
	CustomKey         string          `json:"customKey" jsonschema:"should be set by user in case forceCustomValues is true."`
	Value             string          `json:"value" jsonschema:"Default value for the action, typically empty string (has to be set by user when utilised)."`
}

// WAFCondition ...
type WAFCondition struct {
	ID                int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a WAF Condition it is necessary to add this attribute to your object."`
	Created           *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new WAFRule condition object. This value is informational so it is not necessary to add this attribute to any API call."`
	Modified          *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format."`
	ForceCustomValues bool            `json:"forceCustomValues" jsonschema:"This attributes determines number of input fields when utilised (0=none, 1=value, 2=key+value)."`
	AvailablePhases   int             `json:"availablePhases" jsonschema:"This attributes determines the support for different phases (1=request, 2=response, 3=both)."`
	Alias             string          `json:"alias" jsonschema:"Display name of the condition."`
	Category          string          `json:"category" jsonschema:"Category of the WAF confition."`
	MatchingType      string          `json:"matchingType" jsonschema:"Describes how the values have to match, possible values are EXACT, IREGEX, REGEX, PREFIX, SUFFIX, NOT EXACT, NOT IREGEX, NOT REGEX, NOT PREFIX, NOT SUFFIX ."`
	Name              string          `json:"name" jsonschema:"Type of the condition."`
	Key               string          `json:"key" jsonschema:"Should be set by user in case forceCustomValues is true."`
	Value             string          `json:"value" jsonschema:"Default value for the condition, typically empty string (has to be set by user when utilised)."`
}

// ListWAFConditions returns a list of available WAF conditions
func (api *API) ListWAFConditions() ([]WAFCondition, error) {
	if _, ok := methods["listWAFConditions"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWAFConditions")
	}

	result, err := api.call(methods["listWAFConditions"])
	if err != nil {
		return nil, err
	}

	return *result.(*[]WAFCondition), nil
}

// ListWAFActions returns a list of available WAF actions
func (api *API) ListWAFActions() ([]WAFAction, error) {
	if _, ok := methods["listWAFActions"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWAFActions")
	}

	result, err := api.call(methods["listWAFActions"])
	if err != nil {
		return nil, err
	}

	return *result.(*[]WAFAction), nil
}

// ListWAFRules returns a list of WAF rules.
func (api *API) ListWAFRules(domainId int, params map[string]string) ([]WAFRule, error) {
	if _, ok := methods["listWAFRules"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWAFRules")
	}

	definition := methods["listWAFRules"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]WAFRule), nil
}

// FetchWAFRule returns a single WAF rule for the given ID
func (api *API) FetchWAFRule(id int, params map[string]string) (*WAFRule, error) {
	if _, ok := methods["fetchWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "fetchWAFRule")
	}

	definition := methods["fetchWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	rules := *result.(*[]WAFRule)
	if len(rules) <= 0 {
		return nil, fmt.Errorf("unable to fetch WAF rule for passed id [%d]", id)
	}

	return &rules[0], nil
}

// CreateWAFRule creates a new WAF rule
func (api *API) CreateWAFRule(rule *WAFRule, domainId int, subDomainName string) (*WAFRule, error) {
	if _, ok := methods["createWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createWAFRule")
	}

	definition := methods["createWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}
	return result.(*WAFRule), nil
}

// UpdateWAFRule updates the passed WAF rule
func (api *API) UpdateWAFRule(rule *WAFRule, domainId int, subDomainName string) (*WAFRule, error) {
	if _, ok := methods["updateWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateWAFRule")
	}

	definition := methods["updateWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, rule.ID)

	result, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}
	return result.(*WAFRule), nil
}

// DeleteWAFRule deletes the passed WAF rule
func (api *API) DeleteWAFRule(rule *WAFRule) (*WAFRule, error) {
	if _, ok := methods["deleteWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteWAFRule")
	}

	definition := methods["deleteWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, rule.ID)

	_, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}
