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

// TagWAFRule represents a Web Application Firewall rule linked to a specific Tag.
// It allows applying WAF logic (conditions and actions) to all domains associated with that tag.
type TagWAFRule struct {
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

	// Sort defines the execution order of rules.
	// Lower numbers are processed first.
	Sort int `json:"sort" jsonschema:"The execution order/priority. Lower numbers are processed first."`

	// Sync indicates if the rule is synchronized to the edge nodes.
	Sync bool `json:"sync" jsonschema:"Indicates synchronization status with edge nodes."`

	// ProcessNext controls the rule chain execution flow.
	// If true, subsequent rules are evaluated even if this rule matches.
	ProcessNext bool `json:"processNext" jsonschema:"Flow control: If true, the system continues to process the next rule in the chain after a match. If false, processing stops here."`

	// Enabled controls whether the rule is currently active.
	Enabled bool `json:"enabled" jsonschema:"Indicates if the WAF rule is currently active (enabled) or ignored."`

	// Actions defines what happens when the conditions are met.
	Actions []*WAFAction `json:"actions" jsonschema:"List of actions to execute when the rule conditions are met (e.g., Block, Log, Allow)."`

	// Conditions defines the logical checks (e.g., IP match, Header match) required to trigger the rule.
	Conditions []*WAFCondition `json:"conditions" jsonschema:"List of logical conditions that must be satisfied for the rule to trigger."`

	// TagId is the ID of the Tag this rule belongs to.
	TagId int `json:"tagId" jsonschema:"The ID of the parent Tag to which this WAF rule is attached."`
}

// GetTagWAFRule returns a single tag for the given identifier
func (api *API) GetTagWAFRule(tagId int, ruleId int) (*TagWAFRule, error) {
	if _, ok := api.methods["getTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getTagWAFRule")
	}

	definition := api.methods["getTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, ruleId)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*TagWAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// ListTagWAFRules returns a slice containing all visible tags
func (api *API) ListTagWAFRules(tagId int, params map[string]string) ([]TagWAFRule, error) {
	if _, ok := api.methods["listTagWAFRules"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagWAFRules")
	}

	definition := api.methods["listTagWAFRules"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]TagWAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateTagWAFRule creates a new tag using the MYRA API
func (api *API) CreateTagWAFRule(rule *TagWAFRule, tagId int) (*TagWAFRule, error) {
	if _, ok := api.methods["createTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagWAFRule")
	}

	definition := api.methods["createTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*TagWAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateTagWAFRule updates the passed tag using the MYRA API
func (api *API) UpdateTagWAFRule(rule *TagWAFRule) (*TagWAFRule, error) {
	if _, ok := api.methods["updateTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagWAFRule")
	}

	definition := api.methods["updateTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, rule.TagId, rule.ID)

	result, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*TagWAFRule)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteTagWAFRule deletes the passed tag using the MYRA API
func (api *API) DeleteTagWAFRule(rule *TagWAFRule) (*TagWAFRule, error) {
	if _, ok := api.methods["deleteTagWAFRule"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTagWAFRule")
	}

	definition := api.methods["deleteTagWAFRule"]
	definition.Action = fmt.Sprintf(definition.Action, rule.TagId, rule.ID)

	_, err := api.call(definition, rule)
	if err != nil {
		return nil, err
	}

	return rule, nil
}
