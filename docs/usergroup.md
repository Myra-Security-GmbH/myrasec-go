# User Group

```go
type UserGroup struct {
    ID           int             `json:"id,omitempty"`
    Created      *types.DateTime `json:"created,omitempty"`
    Modified     *types.DateTime `json:"modified,omitempty"`
    Name         string          `json:"name,omitempty"`
    Parent       int             `json:"parent,omitempty"`
    Children     []UserGroup     `json:"children,omitempty"`
    Roles        []string        `json:"roles,omitempty"`
    MembersCount int             `json:"membersCount,omitempty"`
    Type         string          `json:"type,omitempty"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is a unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a UserGroup it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Created will be created by the server after creating a new UserGroup object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an `ISO 8601` format. |
| `Name` | string | The display name of the group. |
| `Parent` | int | The identifier of the parent group. Zero or omitted means the group is a root group. |
| `Children` | []UserGroup | The immediate child groups nested under this group. Read-only, populated by the server on list and read responses. |
| `Roles` | []string | The role identifiers the requesting user holds on this group. Read-only. |
| `MembersCount` | int | The number of users currently assigned to this group. Read-only. |
| `Type` | string | Distinguishes between a regular user group and an agent group. Allowed values: `USER`, `AGENT`. Defaults to `USER` on creation when omitted. |

```go
type GroupRole struct {
    ID       int             `json:"id,omitempty"`
    Created  *types.DateTime `json:"created,omitempty"`
    Modified *types.DateTime `json:"modified,omitempty"`
    UserID   int             `json:"userId"`
    Role     string          `json:"role"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is a unique identifier for the role assignment. Server-generated and read-only. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the role assignment. Server-managed, read-only. |
| `UserID` | int | The identifier of the user being assigned to the group. Required. |
| `Role` | string | The role identifier granted to the user within the group. Allowed values: `ADMINISTRATOR`, `USER`. Required. |

## Create
To create a new user group it is necessary to send a UserGroup object with at least a name. Set `Parent` to nest the new group under an existing group; omit it to create a root group.

### Example
```go
group := &myrasec.UserGroup{
    Name: "Engineering",
}

g, err := api.CreateUserGroupContext(ctx, group)
if err != nil {
    log.Fatal(err)
}

log.Println(g.ID, g.Name)
```

## List
The listing operation returns the user groups visible to the authenticated account. Root groups are returned with their child groups nested under `Children`.

### Example
```go
groups, err := api.ListUserGroupsContext(ctx, nil)
if err != nil {
    log.Fatal(err)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListUserGroupsContext` function.

| name | description | default |
|---|---|---|
| `role` | Restrict the result to groups where the requesting user holds the given role (`ADMINISTRATOR` or `USER`). | null |
| `page` | Specify the page of the result. | 1 |
| `pageSize` | Specify the amount of results in the response. | 50 |

## Read
The read operation returns a single UserGroup by its ID.

### Example
```go
group, err := api.GetUserGroupContext(ctx, groupId)
if err != nil {
    log.Fatal(err)
}
```

## Update
Updating a user group needs the `ID` and `Modified` attributes to identify the object and verify the version.

### Example
```go
group := &myrasec.UserGroup{
    ID: 0000,
    Modified: &types.DateTime{
        Time: modified,
    },
    Name: "Engineering EU",
}

g, err := api.UpdateUserGroupContext(ctx, group)
if err != nil {
    log.Fatal(err)
}
```

## Delete
For deleting a user group it is only necessary to send the `ID` and `Modified` attributes as body content.

### Example
```go
group := &myrasec.UserGroup{
    ID: 0000,
    Modified: &types.DateTime{
        Time: modified,
    },
}

g, err := api.DeleteUserGroupContext(ctx, group)
if err != nil {
    log.Fatal(err)
}
```

## List members of a group
Returns the users that are members of the given group.

### Example
```go
users, err := api.ListUsersFromGroupContext(ctx, groupId, nil)
if err != nil {
    log.Fatal(err)
}
```

| name | description | default |
|---|---|---|
| `search` | Filter by the specified search query. | null |
| `includeRoles` | If `true`, the response includes each user's role assignments on the group. | false |
| `sort` | Sort field (e.g. `login`, `firstname`). | null |
| `page` | Specify the page of the result. | 1 |
| `pageSize` | Specify the amount of results in the response. | 50 |
| `language` | Locale for translated fields. | null |

## Add a user to a group
Assigns an existing user to the group with the given role.

### Example
```go
role := &myrasec.GroupRole{
    UserID: 0000,
    Role:   myrasec.GroupRoleUser,
}

r, err := api.AddUserToGroupContext(ctx, role, groupId)
if err != nil {
    log.Fatal(err)
}

log.Println(r.ID, r.Role)
```

## Remove a user from a group
Removes the passed user from the group identified by `groupId`.

### Example
```go
u, err := api.RemoveUserFromGroupContext(ctx, user, groupId)
if err != nil {
    log.Fatal(err)
}
```
