# IP range

```go
type IPRange struct {
	ID        int             `json:"id,omitempty"`
	Created   *types.DateTime `json:"created,omitempty"`
	Modified  *types.DateTime `json:"modified,omitempty"`
	Network   string          `json:"network"`
	ValidFrom *types.DateTime `json:"validFrom,omitempty"`
	ValidTo   *types.DateTime `json:"validTo,omitempty"`
	Enabled   bool            `json:"enabled,omitempty"`
	Comment   string          `json:"comment,omitempty"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is an unique identifier for an object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. |
| `Modified` | *types.DateTime | Identifies the version of the object. This value is always a date type with an `ISO 8601` format. |
| `Network` | string | The network (CIDR notation). |
| `ValidFrom` | *types.DateTime | |
| `ValidTo` | *types.DateTime | |
| `Enabled` | bool | |
| `Comment` | string | |

## Read
The listing operation returns a list of IP ranges.

### Example
```go
filters, err := api.ListIPRanges(nil)
if err != nil {
    log.Fatal(err)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListIPRanges` function.

| name | description | default |
|---|---|---|
| `search` | Filter by the specified search query | null |
| `ipVersion` | Filter for the specified IP version. Possible values are `ipv4` and `ipv6` | null |
| `page` | Specify the page of the result | 1 |
| `pageSize` | Specify the amount of results in the response | 50 |
| `enabled` | Return only enabled IP ranges | null |