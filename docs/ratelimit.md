# IP rate limit
The IP black/whitelist of Myra lets you grant or deny access from individual IP addresses or subnets.

```go
type RateLimit struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	Burst         int             `json:"burst"`
	Network       string          `json:"network"`
	SubDomainName string          `json:"subDomainName"`
	Timeframe     int             `json:"timeframe"`
	Type          string          `json:"type"`
	Value         int             `json:"value"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | Id is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a IP rate limit setting it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created will be created by the server after creating a new cache setting object. This value is only informational so it is not necessary to add this an attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add the modified timestamp for updates and deletes. |
| `Burst` | int | Burst defines how many requests a client can make in excess of the specified rate. |
| `Network` | string | Network in CIDR notation affected by the rate limiter. |
| `SubDomainName` | string |  |
| `Timeframe` | int | The affected timeframe in seconds for the rate limit. |
| `Type` | string | `tag` or `domain` |
| `Value` | int | Maximum amount of requests for the given network. |


## Create
To create a new IP rate limit setting, you need to send an RateLimit object without the attributes "id", "created" and "modified". All those attributes are generated by the server and returned to you after a successful insert.

### Example
```go
ratelimit := &myrasec.RateLimit{
    Burst:         50,
    Value:         100,
    Timeframe:     60,
    SubDomainName: "www.example.com.",
    Type:          "domain",
    Network:       "127.0.0.1/24",
}
rl, err := api.CreateRateLimit(ratelimit)
if err != nil {
    log.Fatal(err)
}
```


## Read
The listing operation returns a list of IP rate limit settings. The `rateLimitType`, required by this function, can be either `dns` or `tag`.

It is required to pass a map of parameters (`map[string]string`) to the `ListRateLimits` function.

| name | rateLimitType | description | default |
|---|---|---|---|
| `search` (string) | `dns` and `tag` | Filter by the specified search query | null |
| `page` | Specify the page of the result | 1 |
| `pageSize` | Specify the amount of results in the response | 50 |
| `reference` (int) | `dns` and `tag` | filter rate limit settings for this domain or tag (ID) | null |
| `subDomainName` (string) | `dns` | filter rate limit settings for this subdomain (name) | null |

### Example
```go
ratelimits, err := api.ListRateLimits("domain", map[string]string{"subDomainName": "www.example.com"})
if err != nil {
    log.Fatal(err)
}
```


## Update
Updating an IP rate limit setting is very similar to creating a new one. The main difference is that an update will need the "id" and "modified" attributes to identify the version of the object you are trying to update.

### Example
```go
ratelimit := &myrasec.RateLimit{
    ID:   0000,
    Modified: &types.DateTime{
        Time: modified,
    },
    Value:      200,
}

rl, err := api.UpdateRateLimit(ratelimit);
if err != nil {
    log.Fatal(err)
}
```


## Delete
To delete a IP rate limit setting you only need to send "id" and "modified" as body content.

### Example
```go
ratelimit := &myrasec.RateLimit{
    ID:   0000,
    Modified: &types.DateTime{
        Time: modified,
    },
}

rl, err := api.DeleteRateLimit(rateLimit);
if err != nil {
    log.Fatal(err)
}
```