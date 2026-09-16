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
	Types                          []string        `json:"types,omitempty"`
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
| `Modified` | *types.DateTime | Identifies the version of the object. Server-managed, read-only. |
| `UserID` | int | The identifier of the linked user of the organization, if the contact is tied to a user account. |
| `Name` | string | The display name of the contact. |
| `Email` | string | The contact's email address. |
| `Phone` | string | The contact's primary phone number. |
| `SecondaryPhone` | string | The contact's secondary phone number. |
| `PreferredCommunicationLanguage` | string | The preferred communication language. Valid values: `EN`, `DE`. |
| `Comments` | string | Free-form notes about the contact. |
| `Types` | []string | The contact-type keys assigned to this contact (multi-select). Valid keys come from the contact-type catalog (see below). |
| `ReceiveSSLReminders` | bool | Indicates whether the contact receives SSL expiry reminders. |
| `ReceiveTrafficAlerts` | *bool | Indicates whether the contact receives traffic alerts. It is a pointer on purpose: omitting it leaves the stored value unchanged. |
| `TechnicalPriority` | int | The technical contact priority. Valid range: 1-6. |
| `EscalationChain` | int | The escalation priority (escalation chain). Valid range: 1-6. The REST key stays `escalationPriority`. |

## Create
Creating contacts is a bulk operation: pass a slice of contacts. The generated `id`, `created` and `modified` attributes are returned after a successful insert.

### Example
```go
receiveTrafficAlerts := true
contact, err := api.CreateOrganizationContactsContext(ctx, []myrasec.OrganizationContact{
    {
        Name:                           "Alice Admin",
        Email:                          "alice@example.com",
        Phone:                          "+49111",
        PreferredCommunicationLanguage: "EN",
        Types:                          []string{"technical", "billing"},
        ReceiveSSLReminders:            true,
        ReceiveTrafficAlerts:           &receiveTrafficAlerts,
        TechnicalPriority:              1,
        EscalationChain:                2,
    },
})
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
Updating a contact requires the generated `id` to identify the object.

### Example
```go
contact := &myrasec.OrganizationContact{
    ID:    0000,
    Email: "alice@example.com",
    Name:  "Alice Admin",
    Types: []string{"billing"},
}

c, err := api.UpdateOrganizationContactContext(ctx, contact)
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
