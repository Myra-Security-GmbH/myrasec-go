# User

```go
type User struct {
    ID                 int             `json:"id,omitempty"`
    Created            *types.DateTime `json:"created,omitempty"`
    Modified           *types.DateTime `json:"modified,omitempty"`
    Login              string          `json:"login,omitempty"`
    Email              string          `json:"email,omitempty"`
    Firstname          string          `json:"firstname,omitempty"`
    Lastname           string          `json:"lastname,omitempty"`
    OrganizationID     int             `json:"organizationId,omitempty"`
    OrganizationName   string          `json:"organizationName,omitempty"`
    Active             bool            `json:"active,omitempty"`
    Locked             bool            `json:"locked,omitempty"`
    Deleted            bool            `json:"deleted,omitempty"`
    Agent              types.Bool      `json:"agent,omitempty"`
    TfaEnabled         bool            `json:"tfaEnabled,omitempty"`
    TfaRequired        bool            `json:"tfaRequired,omitempty"`
    IsIndirectCustomer bool            `json:"isIndirectCustomer,omitempty"`
    Admin              bool            `json:"admin,omitempty"`
    RootAdmin          bool            `json:"rootAdmin,omitempty"`
    Roles              []UserRole      `json:"roles,omitempty"`
    RootGroupRoles     []UserRole      `json:"rootGroupRoles,omitempty"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is an unique identifier for the user. Server-generated and read-only. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the user object. Server-managed, read-only. |
| `Login` | string | The user's login name. Must be a valid email address (format: `user@example.com`). |
| `Email` | string | The user's contact email address. |
| `Firstname` | string | The user's given name. |
| `Lastname` | string | The user's family name. |
| `OrganizationID` | int | The unique identifier of the organization the user belongs to. |
| `OrganizationName` | string | The display name of the user's organization. |
| `Active` | bool | Indicates whether the user account is currently enabled. |
| `Locked` | bool | Indicates whether the user account is locked (e.g. after failed login attempts). |
| `Deleted` | bool | Indicates whether the user has been soft-deleted. |
| `Agent` | types.Bool | Indicates whether the user has agent (support staff) privileges. The API sends this flag as a string, see below. |
| `TfaEnabled` | bool | Indicates whether two-factor authentication is currently active for this user. |
| `TfaRequired` | bool | Indicates whether two-factor authentication is required for this user. |
| `IsIndirectCustomer` | bool | Indicates whether the user belongs to an indirect-customer organization. |
| `Admin` | bool | Indicates whether the user is an administrator within their organization. |
| `RootAdmin` | bool | Indicates whether the user is a root administrator with platform-wide access. |
| `Roles` | []UserRole | The user's role assignments across all groups they are a member of. |
| `RootGroupRoles` | []UserRole | The user's role assignments restricted to root (top-level) groups. |

```go
type Bool bool
```

`types.Bool` is a `bool` and can be used as such (`if user.Agent { ... }`, or `bool(user.Agent)` where a plain `bool` is required). It only differs in the way it is decoded: the API sends the agent flag as a string (`""` for false, `"1"` for true) instead of a JSON boolean. Both encodings are accepted, as are `null` and the numbers `0` and `1`; anything else is a decoding error. The flag marshals as a plain JSON boolean.

```go
type UserRole struct {
    ID        int             `json:"id,omitempty"`
    Created   *types.DateTime `json:"created,omitempty"`
    Modified  *types.DateTime `json:"modified,omitempty"`
    GroupID   int             `json:"groupId,omitempty"`
    GroupName string          `json:"groupName,omitempty"`
    Role      string          `json:"role,omitempty"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is an unique identifier for the role assignment. Server-generated and read-only. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the role assignment. Server-managed, read-only. |
| `GroupID` | int | The identifier of the group this role applies to. |
| `GroupName` | string | The display name of the group this role applies to. |
| `Role` | string | The role identifier. Allowed values: `ADMINISTRATOR`, `USER`. |

## Get the authenticated user
Returns the User object representing the account currently used to access the API. The response is sparse (typically only `ID`, `Login`, `Created`, `Modified`); other fields are populated by endpoints that return full user objects.

### Example
```go
me, err := api.Me()
if err != nil {
    log.Fatal(err)
}

log.Println(me.ID, me.Login)
```

## Use as a return type
Several endpoints return User objects with the fields populated as documented above. For example, [ListUsersFromGroup](./usergroup.md) returns a list of users that are members of a given group, including the optional fields like `Firstname`, `OrganizationID` and (when `includeRoles=true` is passed) `Roles`.

```go
users, err := api.ListUsersFromGroup(groupId, map[string]string{
    "includeRoles": "true",
})
if err != nil {
    log.Fatal(err)
}

for _, u := range users {
    log.Println(u.ID, u.Login, u.Firstname, u.Lastname)
}
```
