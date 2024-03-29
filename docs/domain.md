# Domain

```go
type Domain struct {
	ID          int             `json:"id,omitempty"`
	Created     *types.DateTime `json:"created,omitempty"`
	Modified    *types.DateTime `json:"modified,omitempty"`
	Name        string          `json:"name"`
	AutoUpdate  bool            `json:"autoUpdate"`
	AutoDNS     bool            `json:"autoDns"`
	Paused      bool            `json:"paused"`
	PausedUntil *types.DateTime `json:"pausedUntil,omitempty"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | Id is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a Domain it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created will be created by the server after creating a new Domain object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. |
| `Name` | string | Identifies the domain by its name. The value cannot be changed after creation. To change a typo you need to remove and recreate the domain. |
| `AutoUpdate` | bool | Shows if the current domain has autoUpdate activated. If autoUpdate is deactivated changes on your configuration are not deployed until you reactivate autoUpdate. This is primary used to change a lot of settings at once to prevent Myra to deploy a half done configuration. In some cases Myra support also deactivates this option to prevent Myra system from removing special configuration settings. Please note that turning autoUpdate off is not correlated to database transactions. This means that any changes are saved but not deployed. |
| `AutoDNS` | bool | If autoDns flag is set while creating a new domain Myra tries to get a list of subDomains for this domain. Depending on your DNS provider configuration this may fail or return a incomplete list. For best results Myra recomments to use the subDomain API to create DNS records. |
| `Paused` | bool | Shows if the domain is currently in pause mode. |
| `PausedUntil` | *types.DateTime | Shows the date when Myra protection will be reactivated automatically. |


## Create
To create a new domain it is necessary to send a Domain object without the attributes "id" and "modified".
Both attributes will be generated by the server and returned after a successful insert is done.

### Example
```go
newDomain := &myrasec.Domain{
    Name: "example.com",
}

d, err := api.CreateDomain(newDomain)
if err != nil {
    panic(err)
}
```


## List
The listing operation returns a list of domains. The list contains domains for the account you are accessing the API with, and also all foreign domains you are allowed to manage.

### Example
```go
domains, err := api.ListDomains(nil)
if err != nil {
    panic(err)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListDomains` function.

| name | description | default |
|---|---|---|
| `search` | Filter by the specified search query | null |
| `page` | Specify the page of the result | 1 |
| `pageSize` | Specify the amount of results in the response | 50 |

## Read
The read operation returns a single domain by it's ID
```go
domain, err := api.GetDomain(domainId)
if err != nil {
    panic(err)
}
```

## Update
Updating a domain is very similar to creating a new one. The main difference is that an update will need
the generated "id" and "modified" attributes to identify the object you are trying to update.

***Updating a domain allows you only to change the autoUpdate flag. All other values are ignored.***

### Example
```go
domain := &myrasec.Domain{
    ID:   0000,
    Name: "example.com",
    Modified: &types.DateTime{
        Time: modified,
    },
    AutoUpdate: false,
}
d, err := api.UpdateDomain(domain)
if err != nil {
    panic(err)
}
```


## Delete
For deleting a domain it is only necessary to send the "id" and "modified" attributes as body content.

***Removing a domains also removes all configurations on Myra!
This could lead to an outage of your online presence. Please make sure that you are prepared for it. If unsure contact Myra support.***

### Example
```go
domain := &myrasec.Domain{
    ID:   0000,
    Name: "example.com",
}
d, err := api.DeleteDomain(domain)
if err != nil {
    panic(err)
}
```

## Fetch a domanin by a domain name
It is possible to fetch a single domain struct by passing the domain name to the `FetchDomain` function. 

### Example
```go
domain, err := api.FetchDomain("example.com")
if err != nil {
    panic(err)
}
```

## Fetch a domain by a subdomain name
It is possible to fetch a single domain struct by passing a subdomain name to the `FetchDomainForSubdomainName` function. 

### Example
```go
domain, err := api.FetchDomainForSubdomainName("www.example.com")
if err != nil {
    panic(err)
}
```