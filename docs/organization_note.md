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
| `Modified` | *types.DateTime | Identifies the version of the note. Server-managed, read-only. |
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
Saving the note is an upsert: the same call creates the note when none exists and edits it otherwise.

### Example
```go
note := &myrasec.OrganizationNote{
    Notes: "remember the maintenance window",
}

saved, err := api.UpdateOrganizationNoteContext(ctx, note)
if err != nil {
    log.Fatal(err)
}
```
