# WAF Rule
```go
type Rule struct {
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
```
| Field | Type | Description |
| --- | --- | --- |
| `ID` | int | Id is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a WAFRule it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created will be created by the server after creating a new WAFRule object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. |
| `ExpireDate` | *types.DateTime | ExpireDate describes how long a WAFRule is valid, and when it will expire |
| `Name` | string | Identifies the tag by its name. |
| `Description` | string | The Description will explain what the WAFRule is for |
| `Direction` | string | The direction can be `in` or `out` |
| `LogIdentifier` | string | A comment to identify the matching rule in the access log. |
| `Sort` | int | Defines the sorting of WAFRules |
| `Sync` | bool | ... |
| `ProcessNext` | bool | After a rule has been applied, the rule chain will be executed as determined. |
| `Enabled` | bool | Describes if the rule is enabled or not |
| `Actions` | []WAFAction | List of WAF actions |
| `Conditions` | []WAFCondition | List of WAF conditions |

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
| `ID` | int | ID of the WAFAction |
| `Created` | *types.DateTime | Created will be created by the server after creating a new WAFRule action object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. |
| `ForeceCustomValues` | bool | This attributes determines number of input fields when utilised (0=none, 1=value, 2=key+value) |
| `AvailablePhases` | int | This attributes determines the support for different phases (1=request, 2=response, 3=both) |
| `Name` | string | Display name of the action |
| `Type` | string | [Type of the action](./waf_action.md) |
| `CustomKey` | string | should be set by user in case `forceCustomValues` is `true` |
| `Value` | string | Default value for the action, typically empty string (has to be set by user when utilised) |

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
| `ID` | int | ID of the WAFCondition |
| `Created` | *types.DateTime | Created will be created by the server after creating a new WAFRule condition object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. |
| `ForeceCustomValues` | bool | This attributes determines number of input fields when utilised (0=none, 1=value, 2=key+value) |
| `AvailablePhases` | int | This attributes determines the support for different phases (1=request, 2=response, 3=both) |
| `Alias` | string | Display name of the condition |
| `Category` | string | Category of the WAF confition |
| `MatchingType` | string | Describes how the values have to match, possible values are `EXACT`, `IREGEX`, `REGEX`, `PREFIX`, `SUFFIX` |
| `Name` | string | [Type of the condition](./waf_condition.md) |
| `Key` | string | should be set by user in case `forceCustomValues` is `true` |
| `Value` | string | Default value for the condition, typically empty string (has to be set by user when utilised) |

## Create
To create a new WAFRule it is necccessary to sent a WAFRule object without the attributes "id" and "modified". Both attributes will be generated by the server and returned after a successful insert is done.

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
    Actions: []WAFAction{},
    Contition: []WAFCondition{},
    Id: 12
}

t, err := api.CreateWAFRule(newWAFRule, domainId, "www.example.com")
if err != nil {
    log.Fatal(err)
}
```

## List
The listing operation returns a list of WAFRules for the given ID.

### Example
```go
rules, err := api.ListWAFRules(domainId, nil)
if err != nil {
    log.Fatal(err)
}
```

## Read
The read operation returns an object of WAFRule for the given ruleId
```go
rule, err := api.FetchWAFRule(ruleId)
if err != nil {
    log.Fatal(err)
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
    Actions: []WAFAction{},
    Contition: []WAFCondition{},
    Id: 12
}

updated, err := api.UpdateWAFRule(rule, domainId, "www.example.com")
if err != nil {
    log.Fatal(err)
}
```

## Delete
For deleting a WAFRule it is only neccessary to send the "id" and "modified" attribute as body content.

### Example
```go
rule := &myrasec.WAFRule{
    ID: 0000,
    Modified: &types.DateTime{
        Time: modified
    }
}

t, err := api.DeleteWAFRule(rule)
if err != nil {
    log.Fatal(err)
}
```
