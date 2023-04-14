# Tag settings

For tag settings you need the object from [Tag](./tag.md) and [Settings](./subdomain_settings.md)

## Create/Update/Delete
To create/update/delete a tag settings you need to send a `map[string]interface{}` with a tagId.

To create/update you add the specific attribute to the map with the required value.

To delete an attribute you have to add the attribute with `nil`.

Only attributes in the map will be touched in the api.
### Example
```go
tagId := 0000
settingsMap['only_https'] = true // update/create
settingsMap['myra_ssl_header'] = nil // delete
ts, err := api.UpdateTagSettingsPartial(settingsMap, tagId)
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