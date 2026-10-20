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
ts, err := api.UpdateTagSettingsPartialContext(ctx, settingsMap, tagId)
if err != nil {
    log.Fatal(err)
}
```

## Read
To list all tag settings you need to send a tagId.

### Example
```go
tagId := 0000
ts, err := api.ListTagSettingsContext(ctx, tagId)
if err != nil {
    log.Fatal(err)
}
```

## Private key settings

`myra_ssl_certificate_key` is write-only. Reads return `["<redacted>"]` for a configured key and retain the empty or inherited representation when no key is configured. `RedactedSettingsPrivateKey` identifies this marker. It contains no key material and cannot be used as a private key.

Sending the marker back preserves the stored key, so a read followed by an update of another setting does not replace it. To replace a key, send the new value. To remove the stored setting, use the partial-update method with `"myra_ssl_certificate_key": nil`. An explicit empty slice in a partial update creates an empty override. Omitting the field preserves it. The typed `Settings` struct uses `omitempty`, so an empty slice there is omitted rather than clearing the setting.
