package myrasec

import (
	"context"
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

// WAFRule represents a single Web Application Firewall rule for a specific subdomain.
// It defines a set of conditions and the actions to take when those conditions are met.
type WAFRule struct {
	// ID is the unique identifier for the WAF rule.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the rule. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the rule was added.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// ExpireDate defines when this rule automatically becomes invalid.
	// If nil, the rule never expires.
	ExpireDate *types.DateTime `json:"expireDate,omitempty" jsonschema:"The timestamp (ISO 8601) when this rule expires. If null, the rule remains valid indefinitely."`

	// Name is a unique label for the rule.
	Name string `json:"name" jsonschema:"A descriptive name to identify the WAF rule."`

	// Description provides further details about the rule's purpose.
	Description string `json:"description" jsonschema:"A detailed description explaining the purpose of this WAF rule."`

	// Direction specifies whether the rule applies to incoming requests or outgoing responses.
	// Valid values: 'in', 'out'.
	Direction string `json:"direction" jsonschema:"The traffic direction to inspect. Valid values: 'in' (incoming request) or 'out' (outgoing response)."`

	// LogIdentifier is a custom tag string used to find matches in the access logs.
	LogIdentifier string `json:"logIdentifier" jsonschema:"A custom string identifier used to tag and find rule matches in the access logs."`

	// UUID is a system-assigned unique identifier string.
	// Read-only.
	Uuid string `json:"uuid,omitempty" jsonschema:"System-assigned unique string identifier (UUID). Read-only."`

	// RuleType identifies the category or logic type of the rule.
	// Typically 'custom_rule' for user-defined rules.
	RuleType string `json:"ruleType" jsonschema:"The type classification of the rule (e.g., 'custom_rule')."`

	// SubDomainName is the FQDN this rule belongs to.
	// Usually set via URL context.
	SubDomainName string `json:"subDomainName" jsonschema:"The FQDN of the subdomain this rule belongs to. Immutable; usually inferred from the URL parameter."`

	// Sort defines the execution order of rules.
	// Lower numbers are processed first.
	Sort int `json:"sort" jsonschema:"The execution order/priority. Lower numbers are processed first."`

	// Sync indicates if the rule is synchronized to the edge nodes.
	Sync bool `json:"sync" jsonschema:"Indicates synchronization status with edge nodes."`

	// Template indicates if this rule serves as a template for others.
	Template bool `json:"template" jsonschema:"If true, this rule is a template and not directly applied to traffic."`

	// ProcessNext controls the rule chain execution flow.
	// If true, subsequent rules are evaluated even if this rule matches.
	ProcessNext bool `json:"processNext" jsonschema:"Flow control: If true, the system continues to process the next rule in the chain after a match. If false, processing stops here."`

	// Enabled controls whether the rule is currently active.
	Enabled bool `json:"enabled" jsonschema:"Indicates if the WAF rule is currently active (enabled) or ignored."`

	// Actions defines what happens when the conditions are met.
	Actions []*WAFAction `json:"actions" jsonschema:"List of actions to execute when the rule conditions are met (e.g., Block, Log)."`

	// Conditions defines the logical checks required to trigger the rule.
	Conditions []*WAFCondition `json:"conditions" jsonschema:"List of logical conditions that must be satisfied for the rule to trigger."`
}

// WAFAction represents an operation executed when a WAF rule triggers.
type WAFAction struct {
	// ID is the unique identifier for the action.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the action. Server-generated; required for updates/deletes, ignored during creation."`

	// Created indicates when the action was added.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601). Required for updates and deletes."`

	// Type defines the specific action to perform (e.g., 'block', 'allow', 'log').
	Type string `json:"type" jsonschema:"The type of action to execute (e.g., 'block', 'log', 'allow')."`

	// Name is the display name of the action type.
	// Read-only; derived from Type.
	Name string `json:"name" jsonschema:"The display name of the action. Read-only."`

	// CustomKey is an optional configuration key for the action.
	// Usage depends on ForceCustomValues.
	CustomKey string `json:"customKey" jsonschema:"Optional configuration key. Usage depends on the specific action type."`

	// Value is the configuration value for the action.
	// Usage depends on ForceCustomValues.
	Value string `json:"value" jsonschema:"The configuration value for the action. Required for certain action types."`

	// ForceCustomValues indicates input requirements for this action type.
	// 0=none, 1=value, 2=key+value. Read-only metadata.
	ForceCustomValues bool `json:"forceCustomValues" jsonschema:"readOnly=true,description=Metadata indicating input requirements: 0=none, 1=value, 2=key+value. Read-only."`

	// AvailablePhases indicates in which request phases this action is valid.
	// 1=request, 2=response, 3=both. Read-only metadata.
	AvailablePhases int `json:"availablePhases" jsonschema:"readOnly=true,description=Metadata indicating valid phases: 1=request, 2=response, 3=both. Read-only."`
}

// WAFCondition represents a logical check within a WAF rule.
type WAFCondition struct {
	// ID is the unique identifier for the condition.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the condition. Server-generated; required for updates/deletes, ignored during creation."`

	// Created indicates when the condition was added.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601). Required for updates and deletes."`

	// Name identifies the type of check (e.g., 'url', 'ip', 'header').
	Name string `json:"name" jsonschema:"The type of the condition (e.g., 'url', 'ip', 'user_agent')."`

	// MatchingType defines the comparison operator.
	// Valid values: EXACT, IREGEX, REGEX, PREFIX, SUFFIX and their NOT variants.
	MatchingType string `json:"matchingType" jsonschema:"The comparison operator. Valid values: 'EXACT', 'IREGEX' (case-insensitive regex), 'REGEX', 'PREFIX', 'SUFFIX', 'NOT EXACT', 'NOT IREGEX', 'NOT REGEX', 'NOT PREFIX', 'NOT SUFFIX'."`

	// Key specifies the target of the check (e.g., the specific header name).
	// Usage depends on the condition Name.
	Key string `json:"key" jsonschema:"The specific target parameter (e.g., the name of the Header to check). Usage depends on condition type."`

	// Value specifies the pattern or content to match against.
	Value string `json:"value" jsonschema:"The value or pattern to match against."`

	// Alias is a human-readable label for the condition type.
	// Read-only; derived from Name.
	Alias string `json:"alias" jsonschema:"The display label for the condition type. Read-only."`

	// Category groups conditions types.
	// Read-only.
	Category string `json:"category" jsonschema:"The category of the condition. Read-only."`

	// ForceCustomValues indicates input requirements for this condition type.
	// 0=none, 1=value, 2=key+value. Read-only metadata.
	ForceCustomValues bool `json:"forceCustomValues" jsonschema:"Metadata indicating input requirements: 0=none, 1=value, 2=key+value. Read-only."`

	// AvailablePhases indicates in which request phases this condition is valid.
	// 1=request, 2=response, 3=both. Read-only metadata.
	AvailablePhases int `json:"availablePhases" jsonschema:"Metadata indicating valid phases: 1=request, 2=response, 3=both. Read-only."`
}

// ListWAFConditionsContext returns a list of available WAF conditions
func (api *API) ListWAFConditionsContext(ctx context.Context) ([]WAFCondition, error) {
	if _, ok := api.methods["listWAFConditions"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWAFConditions")
	}

	result, err := api.call(ctx, api.methods["listWAFConditions"])
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]WAFCondition)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListWAFConditions is equivalent to ListWAFConditionsContext with context.Background().
//
// Deprecated: use ListWAFConditionsContext.
func (api *API) ListWAFConditions() ([]WAFCondition, error) {
	return api.ListWAFConditionsContext(context.Background())
}

// ListWAFActionsContext returns a list of available WAF actions
func (api *API) ListWAFActionsContext(ctx context.Context) ([]WAFAction, error) {
	if _, ok := api.methods["listWAFActions"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWAFActions")
	}

	result, err := api.call(ctx, api.methods["listWAFActions"])
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]WAFAction)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListWAFActions is equivalent to ListWAFActionsContext with context.Background().
//
// Deprecated: use ListWAFActionsContext.
func (api *API) ListWAFActions() ([]WAFAction, error) {
	return api.ListWAFActionsContext(context.Background())
}

// ListWAFRulesContext returns a list of WAF rules.
func (api *API) ListWAFRulesContext(ctx context.Context, domainId int, params map[string]string) ([]WAFRule, error) {
	if _, ok := api.methods["listWAFRules"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listWAFRules")
	}

	definition := api.methods["listWAFRules"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]WAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListWAFRules is equivalent to ListWAFRulesContext with context.Background().
//
// Deprecated: use ListWAFRulesContext.
func (api *API) ListWAFRules(domainId int, params map[string]string) ([]WAFRule, error) {
	return api.ListWAFRulesContext(context.Background(), domainId, params)
}

// FetchWAFRuleContext returns a single WAF rule for the given ID
func (api *API) FetchWAFRuleContext(ctx context.Context, id int, params map[string]string) (*WAFRule, error) {
	if _, ok := api.methods["fetchWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "fetchWAFRule")
	}

	definition := api.methods["fetchWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	rulesPtr, ok := result.(*[]WAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	rules := *rulesPtr
	if len(rules) <= 0 {
		return nil, fmt.Errorf("unable to fetch WAF rule for passed id [%d]", id)
	}

	return &rules[0], nil
}

// FetchWAFRule is equivalent to FetchWAFRuleContext with context.Background().
//
// Deprecated: use FetchWAFRuleContext.
func (api *API) FetchWAFRule(id int, params map[string]string) (*WAFRule, error) {
	return api.FetchWAFRuleContext(context.Background(), id, params)
}

// CreateWAFRuleContext creates a new WAF rule
func (api *API) CreateWAFRuleContext(ctx context.Context, rule *WAFRule, domainId int, subDomainName string) (*WAFRule, error) {
	if _, ok := api.methods["createWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createWAFRule")
	}

	definition := api.methods["createWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(ctx, definition, rule)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*WAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateWAFRule is equivalent to CreateWAFRuleContext with context.Background().
//
// Deprecated: use CreateWAFRuleContext.
func (api *API) CreateWAFRule(rule *WAFRule, domainId int, subDomainName string) (*WAFRule, error) {
	return api.CreateWAFRuleContext(context.Background(), rule, domainId, subDomainName)
}

// UpdateWAFRuleContext updates the passed WAF rule
func (api *API) UpdateWAFRuleContext(ctx context.Context, rule *WAFRule, domainId int, subDomainName string) (*WAFRule, error) {
	if _, ok := api.methods["updateWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateWAFRule")
	}

	definition := api.methods["updateWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, rule.ID)

	result, err := api.call(ctx, definition, rule)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*WAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateWAFRule is equivalent to UpdateWAFRuleContext with context.Background().
//
// Deprecated: use UpdateWAFRuleContext.
func (api *API) UpdateWAFRule(rule *WAFRule, domainId int, subDomainName string) (*WAFRule, error) {
	return api.UpdateWAFRuleContext(context.Background(), rule, domainId, subDomainName)
}

// DeleteWAFRuleContext deletes the passed WAF rule
func (api *API) DeleteWAFRuleContext(ctx context.Context, rule *WAFRule) (*WAFRule, error) {
	if _, ok := api.methods["deleteWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteWAFRule")
	}

	definition := api.methods["deleteWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, rule.ID)

	_, err := api.call(ctx, definition, rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// DeleteWAFRule is equivalent to DeleteWAFRuleContext with context.Background().
//
// Deprecated: use DeleteWAFRuleContext.
func (api *API) DeleteWAFRule(rule *WAFRule) (*WAFRule, error) {
	return api.DeleteWAFRuleContext(context.Background(), rule)
}
