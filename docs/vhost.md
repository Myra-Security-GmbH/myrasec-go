# VHost

```go
type VHost struct {
	ID         int    `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Value      string `json:"value,omitempty"`
	DomainName string `json:"domainName,omitempty"`
	Access     bool   `json:"access"`
	Paused     bool   `json:"paused"`
}
```

| Field | Type | Description|
|---|---|---|
| `ID` | int | ID is a unique identifier for an object. This value is always a number type and is server-generated. |
| `Label` | string | A descriptive name for the VHost without trailing dot (e.g., `shop.example.com`). |
| `Value` | string | The specific hostname or FQDN associated with this VHost (e.g., `shop.example.com.`). |
| `DomainName` | string | The Fully Qualified Domain Name (FQDN) of the parent domain. |
| `Access` | bool | Indicates if the VHost is configured to be accessible (active). |
| `Paused` | bool | Indicates if the VHost is currently paused. If true, traffic processing is suspended. |

## List all subdomains
The listing operation returns a list of all subdomains (VHosts) across all domains.

### Example
```go
subdomains, err := api.ListAllSubdomains(map[string]string{})
if err != nil {
    log.Fatal(err)
}

for _, s := range subdomains {
    log.Println(s.ID, s.Label, s.DomainName)
}
```

## List subdomains for a domain
Returns a list of subdomains (VHosts) for a specific domain, identified by domain ID.

### Example
```go
domainId := 1234

subdomains, err := api.ListAllSubdomainsForDomain(domainId, map[string]string{})
if err != nil {
    log.Fatal(err)
}

for _, s := range subdomains {
    log.Println(s.ID, s.Label, s.Value)
}
```
