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

| Field | Type | Description |
|---|---|---|
| `ID` | int | Id is an unique identifier for an object. |
| `Created` | *types.DateTime | Created will be created by the server after creating a new cache setting object. This value is only informational so it is not necessary to add this an attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. |
| `Network` | string | The nework (CIDR notation) |
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

| Name | Description | Default |
|---|---|---|
| `search` | Filter by the specified search query | null |
| `page` | Specify the page of the result | 1 |
| `pageSize` | Specify the amount of results in the response | 50 |
| `enabled` | Return only enabled IP ranges | null |