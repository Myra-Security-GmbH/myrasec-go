# Organization Contact
Organization contacts are the contact persons of your organization. The organization is resolved from the authenticated session, so contacts carry no domain or subdomain context.

```go
type OrganizationContact struct {
	ID                             int             `json:"id,omitempty"`
	Created                        *types.DateTime `json:"created,omitempty"`
	Modified                       *types.DateTime `json:"modified,omitempty"`
	UserID                         int             `json:"userId,omitempty"`
	Name                           string          `json:"name,omitempty"`
	Email                          string          `json:"email"`
	Phone                          string          `json:"phone,omitempty"`
	SecondaryPhone                 string          `json:"secondaryPhone,omitempty"`
	PreferredCommunicationLanguage string          `json:"preferredCommunicationLanguage,omitempty"`
	Comments                       string          `json:"comments,omitempty"`
	Types                          []string        `json:"types,omitzero"`
	ReceiveSSLReminders            bool            `json:"receiveSSLReminders"`
	ReceiveTrafficAlerts           *bool           `json:"receiveTrafficAlerts,omitempty"`
	TechnicalPriority              int             `json:"technicalPriority,omitempty"`
	EscalationChain                int             `json:"escalationPriority,omitempty"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is a unique identifier for the contact. Server-generated; required for updates and deletes, ignored during creation. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the object (`ISO 8601`). Read-only on create; must be echoed back on update so the server can detect concurrent modifications. |
| `UserID` | int | The identifier of the linked user of the organization, if the contact is tied to a user account. |
| `Name` | string | The display name of the contact. |
| `Email` | string | The contact's email address. |
| `Phone` | string | The contact's primary phone number. |
| `SecondaryPhone` | string | The contact's secondary phone number. |
| `PreferredCommunicationLanguage` | string | The preferred communication language. Valid values: `EN`, `DE`. |
| `Comments` | string | Free-form notes about the contact. |
| `Types` | []string | The contact-type keys assigned to this contact (multi-select). Valid keys come from the contact-type catalog (see below). Uses `omitzero`: omit (nil) to leave the stored types unchanged, or send an empty slice to remove all of them. |
| `ReceiveSSLReminders` | bool | Indicates whether the contact receives SSL expiry reminders. |
| `ReceiveTrafficAlerts` | *bool | Indicates whether the contact receives traffic alerts. It is a pointer on purpose: omitting it leaves the stored value unchanged. |
| `TechnicalPriority` | int | The `TechnicalPriority` ranking of the contact. Valid range: 1-6. |
| `EscalationChain` | int | The escalation priority (escalation chain). Valid range: 1-6. The REST key stays `escalationPriority`. Only stored when the contact also carries the `ESCALATION` type; submitting `types` without `ESCALATION` resets it. |

## Create
Creating contacts is a bulk operation: pass a slice of contacts. The API acknowledges the create with an empty body, so the returned slice is the input echoed back and carries **no** server-generated `id`, `created` or `modified` values. Call `ListOrganizationContactsContext` afterwards to obtain the persisted contacts with their ids.

Set `EscalationChain` only together with the `ESCALATION` contact type: the server stores the escalation priority solely when that type is present.

### Example
```go
receiveTrafficAlerts := true
_, err := api.CreateOrganizationContactsContext(ctx, []myrasec.OrganizationContact{
    {
        Name:                           "Alice Admin",
        Email:                          "alice@example.com",
        Phone:                          "+49111",
        PreferredCommunicationLanguage: "EN",
        Types:                          []string{"SALES", "ESCALATION"},
        ReceiveSSLReminders:            true,
        ReceiveTrafficAlerts:           &receiveTrafficAlerts,
        TechnicalPriority:              1,
        EscalationChain:                2,
    },
})
if err != nil {
    log.Fatal(err)
}

// The created contacts have no id yet; re-list to obtain them.
contacts, err := api.ListOrganizationContactsContext(ctx, nil)
if err != nil {
    log.Fatal(err)
}
```

## List
The listing operation returns the contacts of the authenticated organization.

### Example
```go
contacts, err := api.ListOrganizationContactsContext(ctx, nil)
if err != nil {
    log.Fatal(err)
}
```

It is possible to pass a map of parameters (`map[string]string`) such as `search`, `page` and `pageSize` to the `ListOrganizationContactsContext` function.

## Update
The update is a **full replace**: fields you omit are cleared server-side, and the `modified` value must be echoed back so the server can detect concurrent modifications. Treat it as a read-modify-write — fetch the contact, change the fields you want, and send the whole object back.

### Example
```go
// Fetch the current contacts and pick the one to change.
contacts, err := api.ListOrganizationContactsContext(ctx, nil)
if err != nil {
    log.Fatal(err)
}

contact := contacts[0]
contact.Name = "Alice Admin"

// contact still carries its ID and Modified, so the full-replace update keeps the
// remaining fields and passes the optimistic-lock check.
c, err := api.UpdateOrganizationContactContext(ctx, &contact)
if err != nil {
    log.Fatal(err)
}
```

## Delete
For deleting a contact it is only necessary to send the `id` attribute.

### Example
```go
contact := &myrasec.OrganizationContact{ID: 0000}

_, err := api.DeleteOrganizationContactContext(ctx, contact)
if err != nil {
    log.Fatal(err)
}
```

## Contact types
The contact-type catalog lists the valid keys for the `Types` multi-select field.

```go
type OrganizationContactType struct {
	ID    int    `json:"id,omitempty"`
	Key   string `json:"key"`
	Label string `json:"label"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | The unique identifier for the contact type. Server-generated and read-only. |
| `Key` | string | The machine-readable key used in a contact's `types` multi-select. |
| `Label` | string | The human-readable label of the contact type. |

### Example
```go
types, err := api.ListOrganizationContactTypesContext(ctx)
if err != nil {
    log.Fatal(err)
}

for _, t := range types {
    log.Println(t.Key, t.Label)
}
```
