# Permission

```go
type ObjectPermission struct {
    ID                   int                `json:"id,omitempty"`
    Created              *types.DateTime    `json:"created,omitempty"`
    Modified             *types.DateTime    `json:"modified,omitempty"`
    Action               string             `json:"action,omitempty"`
    ObjectType           string             `json:"objectType,omitempty"`
    ObjectPermissionType string             `json:"objectPermissionType,omitempty"`
    ObjectInstance       int                `json:"objectInstance,omitempty"`
    Name                 string             `json:"name,omitempty"`
    Scopes               []int              `json:"scopes,omitempty"`
    Parents              []ObjectPermission `json:"parents,omitempty"`
    Value                string             `json:"value,omitempty"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is a unique identifier for the permission. Server-generated; required for delete operations, ignored on create. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the permission. Server-managed, read-only. |
| `Action` | string | The action being granted. Allowed values: `READ`, `CREATE`, `EDIT`, `SWITCH`, `ADMIN`, `AGENT_MODE`, `USER_MODE`, `CACHE_CLEAR`, `REVIEW`, `PUBLISH`, `LOADBALANCE_READ`, `LOADBALANCE_EDIT`. |
| `ObjectType` | string | The entity type the permission applies to (e.g. `Domain`, `DnsRecord`). |
| `ObjectPermissionType` | string | The recipient kind of the permission. Allowed values: `USER`, `GROUP`. |
| `ObjectInstance` | int | The identifier of a specific object instance the permission is scoped to. Omit to apply to all instances of the object type within the visible scope. |
| `Name` | string | An optional human-readable label for the permission. |
| `Scopes` | []int | Identifiers of additional permissions that act as scope restrictions for this permission. |
| `Parents` | []ObjectPermission | Permissions this one inherits from. Read-only, populated by the server. |
| `Value` | string | An optional auxiliary value carried with the permission. |

```go
type PermissionCheckResult struct {
    IsAuthorized bool `json:"isAuthorized"`
}
```
| Field | Type | Description|
|---|---|---|
| `IsAuthorized` | bool | True if the requesting user is allowed to perform the requested action on the requested target. |

## Action and recipient constants
The package exposes string constants matching the API's allowed values.

```go
myrasec.PermissionActionRead              // "READ"
myrasec.PermissionActionCreate            // "CREATE"
myrasec.PermissionActionEdit              // "EDIT"
myrasec.PermissionActionSwitch            // "SWITCH"
myrasec.PermissionActionAdmin             // "ADMIN"
myrasec.PermissionActionAgentMode         // "AGENT_MODE"
myrasec.PermissionActionUserMode          // "USER_MODE"
myrasec.PermissionActionCacheClear        // "CACHE_CLEAR"
myrasec.PermissionActionReview            // "REVIEW"
myrasec.PermissionActionPublish           // "PUBLISH"
myrasec.PermissionActionLoadbalanceRead   // "LOADBALANCE_READ"
myrasec.PermissionActionLoadbalanceEdit   // "LOADBALANCE_EDIT"

myrasec.PermissionTypeUser                // "USER"
myrasec.PermissionTypeGroup               // "GROUP"
```

## List the permissions of the authenticated user

### Example
```go
permissions, err := api.ListMyPermissionsContext(ctx, nil)
if err != nil {
    log.Fatal(err)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListMyPermissionsContext` function.

| name | description | default |
|---|---|---|
| `actions` | Comma-separated list of actions to filter by (e.g. `READ,EDIT`). | null |
| `objects` | Comma-separated list of object types to filter by (e.g. `Domain,DnsRecord`). | null |
| `language` | Locale for translated fields. | null |

## Check a permission for the authenticated user
Asks the API whether the current user is allowed to perform the given action on the referenced object type or instance.

### Example
```go
check := &myrasec.ObjectPermission{
    Action:         myrasec.PermissionActionEdit,
    ObjectType:     "Domain",
    ObjectInstance: 12345,
}

result, err := api.CheckMyPermissionContext(ctx, check)
if err != nil {
    log.Fatal(err)
}

if result.IsAuthorized {
    log.Println("user can edit domain 12345")
}
```

## List the permissions of a specific user

### Example
```go
permissions, err := api.ListUserPermissionsContext(ctx, userId, nil)
if err != nil {
    log.Fatal(err)
}
```

## Grant a permission to a user

### Example
```go
permission := &myrasec.ObjectPermission{
    Action:         myrasec.PermissionActionRead,
    ObjectType:     "Domain",
    ObjectInstance: 12345,
}

p, err := api.AddUserPermissionContext(ctx, permission, userId)
if err != nil {
    log.Fatal(err)
}

log.Println(p.ID)
```

## Revoke a permission from a user
Pass the permission obtained from `ListUserPermissionsContext` (or stored after `AddUserPermissionContext`). `ID` and `Modified` are required for the revoke.

### Example
```go
p, err := api.RemoveUserPermissionContext(ctx, permission, userId)
if err != nil {
    log.Fatal(err)
}
```

## List the permissions of a user group

### Example
```go
permissions, err := api.ListUserGroupPermissionsContext(ctx, groupId, nil)
if err != nil {
    log.Fatal(err)
}
```

## Grant a permission to a user group

### Example
```go
permission := &myrasec.ObjectPermission{
    Action:         myrasec.PermissionActionEdit,
    ObjectType:     "Domain",
    ObjectInstance: 12345,
}

p, err := api.AddUserGroupPermissionContext(ctx, permission, groupId)
if err != nil {
    log.Fatal(err)
}
```

## Revoke a permission from a user group

### Example
```go
p, err := api.RemoveUserGroupPermissionContext(ctx, permission, groupId)
if err != nil {
    log.Fatal(err)
}
```
