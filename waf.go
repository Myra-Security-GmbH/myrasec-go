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
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	ExpireDate    *types.DateTime `json:"expireDate,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Direction     string          `json:"direction"`
	LogIdentifier string          `json:"logIdentifier"`
	RuleType      string          `json:"ruleType"`
	SubDomainName string          `json:"subDomainName"`
	Sort          int             `json:"sort"`
	Sync          bool            `json:"sync"`
	Template      bool            `json:"template"`
	ProcessNext   bool            `json:"processNext"`
	Enabled       bool            `json:"enabled"`
	Actions       []*WAFAction    `json:"actions"`
	Conditions    []*WAFCondition `json:"conditions"`
}

// WAFAction ...
type WAFAction struct {
	ID                int             `json:"id,omitempty"`
	Created           *types.DateTime `json:"created,omitempty"`
	Modified          *types.DateTime `json:"modified,omitempty"`
	ForceCustomValues bool            `json:"forceCustomValues"`
	AvailablePhases   int             `json:"availablePhases"`
	Name              string          `json:"name"`
	Type              string          `json:"type"`
	CustomKey         string          `json:"customKey"`
	Value             string          `json:"value"`
}

// WAFCondition ...
type WAFCondition struct {
	ID                int             `json:"id,omitempty"`
	Created           *types.DateTime `json:"created,omitempty"`
	Modified          *types.DateTime `json:"modified,omitempty"`
	ForceCustomValues bool            `json:"forceCustomValues"`
	AvailablePhases   int             `json:"availablePhases"`
	Alias             string          `json:"alias"`
	Category          string          `json:"category"`
	MatchingType      string          `json:"matchingType"`
	Name              string          `json:"name"`
	Key               string          `json:"key"`
	Value             string          `json:"value"`
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
	var conditions []WAFCondition
	conditions = append(conditions, *result.(*[]WAFCondition)...)

	return conditions, nil
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
	var actions []WAFAction
	actions = append(actions, *result.(*[]WAFAction)...)

	return actions, nil
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
	var rules []WAFRule
	rules = append(rules, *result.(*[]WAFRule)...)

	return rules, nil
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

	var rules []WAFRule
	rules = append(rules, *result.(*[]WAFRule)...)

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
