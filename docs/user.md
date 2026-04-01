# User

```go
type User struct {
	ID       int             `json:"id,omitempty"`
	Created  *types.DateTime `json:"created,omitempty"`
	Modified *types.DateTime `json:"modified,omitempty"`
	Login    string          `json:"login,omitempty"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is an unique identifier for an object. This value is always a number type and is server-generated. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. This value is server-managed and read-only. |
| `Modified` | *types.DateTime | Identifies the version of the object. This value is always a date type with an `ISO 8601` format. |
| `Login` | string | The user's login name. Must be a valid email address (e.g., `user@example.com`). |

## Read
The `Me` function returns the information of the currently authenticated user.

### Example
```go
user, err := api.Me()
if err != nil {
    log.Fatal(err)
}

log.Println(user.ID, user.Login)
```
