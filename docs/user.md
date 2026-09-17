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
    PrimaryPhone       string          `json:"primaryPhone,omitempty"`
    SecondaryPhone     string          `json:"secondaryPhone,omitempty"`
    PreferredCommunicationLanguage string `json:"preferredCommunicationLanguage,omitempty"`
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
| `ID` | int | ID is a unique identifier for the user. Server-generated and read-only. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the user object. Server-managed, read-only. |
| `Login` | string | The user's login name. Must be a valid email address (format: `user@example.com`). |
| `Email` | string | The user's contact email address. |
| `Firstname` | string | The user's given name. |
| `Lastname` | string | The user's family name. |
| `PrimaryPhone` | string | The user's primary phone number. |
| `SecondaryPhone` | string | The user's secondary phone number. |
| `PreferredCommunicationLanguage` | string | The user's preferred communication language. |
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
| `ID` | int | ID is a unique identifier for the role assignment. Server-generated and read-only. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the role assignment. Server-managed, read-only. |
| `GroupID` | int | The identifier of the group this role applies to. |
| `GroupName` | string | The display name of the group this role applies to. |
| `Role` | string | The role identifier. Allowed values: `ADMINISTRATOR`, `USER`. |

## Get the authenticated user
Returns the User object representing the account currently used to access the API. The response is sparse (typically only `ID`, `Login`, `Created`, `Modified`); other fields are populated by endpoints that return full user objects.

### Example
```go
me, err := api.MeContext(ctx)
if err != nil {
    log.Fatal(err)
}

log.Println(me.ID, me.Login)
```

## List
The listing operation returns the users visible to the authenticated account. An account holding the `ADMINISTRATOR` role in at least one group (root or sub group) gets every user of its organization, any other account gets only itself. Agent accounts cannot call this endpoint while in agent mode (the API answers `403 Forbidden`).

The list carries the base user fields only: `Admin`, `RootAdmin`, `Roles` and `RootGroupRoles` are not populated by this endpoint.

### Example
```go
users, err := api.ListUsersContext(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, u := range users {
    log.Println(u.ID, u.Login, u.Firstname, u.Lastname)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListUsersContext` function.

| name | description | default |
|---|---|---|
| `search` | Restrict the result to users whose login, email, first name, last name or full name contains the search term. | null |
| `sort` | Sort the result by a user attribute, `field:direction` with `asc` or `desc`, e.g. `login:asc`. Several fields can be separated by commas. A value without the direction is ignored, a field that is not a user attribute is rejected. | null |
| `page` | Specify the page of the result. | 1 |
| `pageSize` | Specify the amount of results in the response. | 50 |

```go
users, err := api.ListUsersContext(ctx, map[string]string{
    myrasec.ParamSearch:   "example.com",
    "sort":                "lastname:asc,firstname:asc",
    myrasec.ParamPage:     "1",
    myrasec.ParamPageSize: "25",
})
```

## Update
Updates a user identified by its `ID`. Profile fields such as `PrimaryPhone`, `SecondaryPhone` and `PreferredCommunicationLanguage` are writable through this call.

The update is a **full replace**: fields you omit are cleared server-side, and the `Modified` value must be echoed back so the server can detect concurrent modifications. Treat it as a read-modify-write — fetch the user, change the fields you want, and send the whole object back. Because it is a full replace, the `Active`, `Locked` and `Deleted` flags are always written: set `Active` to `false` to deactivate a user.

### Example
```go
// Fetch the user via the list endpoint, then change the fields to update.
users, err := api.ListUsersContext(ctx, map[string]string{"search": "user@example.com"})
if err != nil {
    log.Fatal(err)
}

user := users[0]
user.PrimaryPhone = "+49111"
user.SecondaryPhone = "+49222"
user.PreferredCommunicationLanguage = "EN"

// user still carries its ID and Modified, so the full-replace update keeps the
// remaining fields and passes the optimistic-lock check.
updated, err := api.UpdateUserContext(ctx, &user)
if err != nil {
    log.Fatal(err)
}
```

## Use as a return type
Several endpoints return User objects with the fields populated as documented above. For example, [ListUsersFromGroup](./usergroup.md) returns a list of users that are members of a given group, including the optional fields like `Firstname`, `OrganizationID` and (when `includeRoles=true` is passed) `Roles`.

```go
users, err := api.ListUsersFromGroupContext(ctx, groupId, map[string]string{
    "includeRoles": "true",
})
if err != nil {
    log.Fatal(err)
}

for _, u := range users {
    log.Println(u.ID, u.Login, u.Firstname, u.Lastname)
}
```
