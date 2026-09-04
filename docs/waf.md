# WAF Rule
```go
type WAFRule struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	ExpireDate    *types.DateTime `json:"expireDate,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Direction     string          `json:"direction"`
	LogIdentifier string          `json:"logIdentifier"`
	Uuid          string          `json:"uuid,omitempty"`
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
```
| Field | Type | Description |
| --- | --- | --- |
| `ID` | int | ID is a unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a WAFRule it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Created will be created by the server after creating a new WAFRule object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an `ISO 8601` format. |
| `ExpireDate` | *types.DateTime | ExpireDate describes how long a WAFRule is valid and when it will expire. |
| `Name` | string | Identifies the WAF rule by its name. |
| `Description` | string | The Description will explain what the WAFRule is for. |
| `Direction` | string | The direction can be `in` (for Request) or `out` (for Response). |
| `LogIdentifier` | string | A string to identify the matching rule in the access log. |
| `Uuid` | string | System-assigned unique string identifier (UUID). Read-only. |
| `RuleType` | string | The type classification of the rule (e.g., 'custom_rule'). |
| `SubDomainName` | string | The FQDN of the subdomain this rule belongs to. Immutable; usually inferred from the URL parameter. |
| `Sort` | int | Defines the sorting of WAFRules. |
| `Sync` | bool | If true, the rule will be propagated to the edge nodes after create/update. Typically set to true. |
| `Template` | bool | If true, this rule is a template and not directly applied to traffic. |
| `ProcessNext` | bool | After a rule has been applied, the rule chain will be executed as determined. |
| `Enabled` | bool | Describes if the rule is enabled or not. |
| `Actions` | []*WAFAction | List of WAF actions. |
| `Conditions` | []*WAFCondition | List of WAF conditions. |

```go
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
```
| Field | Type | Description |
| --- | --- | --- |
| `ID` | int | ID is a unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a WAF Action it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Created will be created by the server after creating a new WAFRule action object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an `ISO 8601` format. |
| `ForeceCustomValues` | bool | This attributes determines number of input fields when utilised (0=none, 1=value, 2=key+value). |
| `AvailablePhases` | int | This attributes determines the support for different phases (1=request, 2=response, 3=both). |
| `Name` | string | Display name of the action. |
| `Type` | string | [Type of the action.](./waf_action.md) |
| `CustomKey` | string | should be set by user in case `forceCustomValues` is `true`. |
| `Value` | string | Default value for the action, typically empty string (has to be set by user when utilised). |

```go
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
```
| Field | Type | Description |
| --- | --- | --- |
| `ID` | int | ID is a unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a WAF Condition it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Created will be created by the server after creating a new WAFRule condition object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an `ISO 8601` format. |
| `ForeceCustomValues` | bool | This attributes determines number of input fields when utilised (0=none, 1=value, 2=key+value). |
| `AvailablePhases` | int | This attributes determines the support for different phases (1=request, 2=response, 3=both). |
| `Alias` | string | Display name of the condition. |
| `Category` | string | Category of the WAF condition. |
| `MatchingType` | string | Describes how the values have to match, possible values are `EXACT`, `IREGEX`, `REGEX`, `PREFIX`, `SUFFIX`, `NOT EXACT`, `NOT IREGEX`, `NOT REGEX`, `NOT PREFIX`, `NOT SUFFIX`. |
| `Name` | string | [Type of the condition.](./waf_condition.md) |
| `Key` | string | Should be set by user in case `forceCustomValues` is `true`. |
| `Value` | string | Default value for the condition, typically empty string (has to be set by user when utilised). |

## Create
To create a new WAFRule it is necessary to send a WAFRule object without the attributes "id" and "modified". Both attributes will be generated by the server and returned after a successful insert is done.

### Example
```go
newWAFRule := &myrasec.WAFRule{
    Name: "RuleName",
    Description: "Example WAFRule",
    Direction: "in",
    LogIdentifier: "Example-Log",
    Sort: 1,
    Sync: true,
    ProcessNext: true,
    Enabled: true,
    Actions: []*myrasec.WAFAction{},
    Conditions: []*myrasec.WAFCondition{},
}

t, err := api.CreateWAFRuleContext(ctx, newWAFRule, domainId, "www.example.com")
if err != nil {
    log.Fatal(err)
}
```

## List
The listing operation returns a list of WAFRules for the given ID.

### Example
```go
rules, err := api.ListWAFRulesContext(ctx, domainId, nil)
if err != nil {
    log.Fatal(err)
}
```

WAFRules can also be fetched for a single subdomain.

### Example

```go
rules, err := api.ListWAFRulesContext(ctx, domainId, map[string]string{
	"subDomain": "www.example.com",
})
if err != nil {
    log.Fatal(err)
}
```

## Read
The read operation returns an object of WAFRule for the given ruleId

### Example
```go
rule, err := api.FetchWAFRuleContext(ctx, ruleId)
if err != nil {
    log.Fatal(err)
}
```

An easy way to access the actions/conditions of a WAF Rule is to range over them:

### Example

```go
rule, err := api.FetchWAFRuleContext(ctx, ruleId)
if err != nil {
	log.Fatal(err)
}

for _, action := range rule.Actions {
	// ...
}

for _, condition := range rule.Conditions {
	// ...
}
```

## Update
Updating a WAFRule is very similar to creating a new one. The main difference is that an update will need the generated "id" and "modified" attribute to identify the object you are trying to update.

### Example
```go
rule := &myrasec.WAFRule{
    ID: 0000,
    Modified: &types.DateTime{
        Time: modified
    },
    Name: "RuleName",
    Description: "Updated WAFRule",
    Direction: "in",
    LogIdentifier: "Example-Log",
    Sort: 1,
    Sync: true,
    ProcessNext: true,
    Enabled: true,
    Actions: []*myrasec.WAFAction{},
    Conditions: []*myrasec.WAFCondition{},
}

updated, err := api.UpdateWAFRuleContext(ctx, rule, domainId, "www.example.com")
if err != nil {
    log.Fatal(err)
}
```

## Delete
For deleting a WAFRule it is only necessary to send the "id" and "modified" attribute as body content.

### Example
```go
rule := &myrasec.WAFRule{
    ID: 0000,
    Modified: &types.DateTime{
        Time: modified
    }
}

t, err := api.DeleteWAFRuleContext(ctx, rule)
if err != nil {
    log.Fatal(err)
}
```
