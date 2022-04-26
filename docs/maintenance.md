# Maintenance
The maintenance functionality allows you to set a maintenance page which is directly served from Myra
servers. This is useful when you are going to maintain your servers and remove all load from your server. 

```go
type Maintenance struct {
	ID          int             `json:"id,omitempty"`
	Created     *types.DateTime `json:"created,omitempty"`
	Modified    *types.DateTime `json:"modified,omitempty"`
	Start       *types.DateTime `json:"start,omitempty"`
	End         *types.DateTime `json:"end,omitempty"`
	Active      bool            `json:"active"`
	Content     string          `json:"content"`
	ContentFrom string          `json:"contentFrom,omitempty"`
    FQDN        string          `json:"fqdn"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | Id is an unique identifier for an object. This value is always a number type and cannot be set, while inserting a new object. To update or delete a Maintenance it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an ISO8601 format. It will be created by the server after creating a new Maintenance object. This value is only informational so it is not necessary to add this an attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletions. This value is always a date type with an ISO8601 format. |
| `Start` | string | Start is a date type attribute with an ISO8601. This attribute shows the start date for a maintenance. This date have to be lower than end or null to start now. |
| `End` | string | End is a date type attribute with an ISO8601 and shows the end date for a maintenance. This date have to be higher than start or null to end now. |
| `Active` | string | This information shows if this a maintenance page is currently active. You cannot set this attribute directly instead you have to set start and end attribute to activate maintenance. |
| `Content` | string | HTML content to show as maintenance page. Please note that it is not possible to include resources from the domain you have set to maintenance mode. If your maintenance page contains images use a different domain or use inline base64 encoded images. |
| `ContentFrom` | string | This property can be used instead of the property content to reference an existing maintenance page’s content. Instead of sending the actual content, specify a valid FQDN here. This will copy the content from the referenced maintenance page to the newly created. |
| `FQDN` | string | Shows a FQDN (fully qualified domain name) for a maintenance. This attribute shows the domain to handle maintenance for. |
