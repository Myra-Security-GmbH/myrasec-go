# SSL Configuration

```go
type SslConfiguration struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Ciphers   string `json:"ciphers"`
	Protocols string `json:"protocols"`
}
```

| Field | Type | Description |
|---|---|---|
| `ID` | int | Identifier of that configuration |
| `Name` | string | Unique string identifier of this configuration |
| `Ciphers` | string | List of ciphers for this configuration |
| `Protocols` | string | List of protocols for this configuration |

## List configurations
To get a list of all valid SSL configurations you can simply call this method:
```go
configurations, err := api.ListSSLConfigurations()
if err != nil {
    log.Fatal(err)
}
```