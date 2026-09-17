# Organization Note
Each organization has a single free-form note. Both creating and editing the note go through the same upsert (save) route, and reading an organization without a saved note returns an empty note shape.

```go
type OrganizationNote struct {
	ID       int             `json:"id,omitempty"`
	Created  *types.DateTime `json:"created,omitempty"`
	Modified *types.DateTime `json:"modified,omitempty"`
	Notes    string          `json:"notes"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | The unique identifier for the note. Server-generated and read-only. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. Server-managed, read-only. |
| `Modified` | *types.DateTime | Identifies the version of the note (`ISO 8601`). Read-only on create; must be echoed back when saving an existing note so the server can detect concurrent modifications. |
| `Notes` | string | The free-form note text of the organization. |

## Get
Returns the note of the authenticated organization. When no note has been saved yet, an empty note shape is returned.

### Example
```go
note, err := api.GetOrganizationNoteContext(ctx)
if err != nil {
    log.Fatal(err)
}

log.Println(note.Notes)
```

## Save
Saving the note is an upsert: the same call creates the note when none exists and edits it otherwise. When editing an existing note, echo its `Modified` value back so the optimistic-lock check passes — treat it as a read-modify-write.

### Example
```go
// Fetch the current note, change the text, and save it back.
note, err := api.GetOrganizationNoteContext(ctx)
if err != nil {
    log.Fatal(err)
}

note.Notes = "remember the maintenance window"

// note still carries its Modified value, so the save passes the version check.
saved, err := api.UpdateOrganizationNoteContext(ctx, note)
if err != nil {
    log.Fatal(err)
}
```
