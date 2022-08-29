# Tag settings

For tag settings you need the object from [Tag](./tag.md) and [Settings](./subdomain_settings.md)

## Create/Update
To create/update a new tag settings you need to send a Settings object with a tagId.

To create a new tag settings, you need to leave the attributes "id", "created", "modified" empty. To update you have to set these attributes.

### Example
```go
// Create tag settings with default values
tagId := 0000
settings := &myrasec.Settings{}
ts, err := api.UpdateTagSettings(settings, tagId)
if err != nil {
    log.Fatal(err)
}

// Update tag settings
ts.AccessLog = true
updated, err := api.UpdateTagSettings(ts, tagId)
if err != nil {
    log.Fatal(err)
}
```

## Read
To list all tag settings you need to send a tagId.

### Example
```go
tagId := 0000
ts, err := api.ListTagSettings(tagId)
if err != nil {
    log.Fatal(err)
}
```