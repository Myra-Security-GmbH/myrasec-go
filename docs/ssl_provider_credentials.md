# SSL provider credentials

SSL provider credentials store the account of an enterprise certificate authority (Sectigo or D-Trust) that Myra uses to request [managed certificates](./ssl_certificate_request.md). A certificate request references the credentials through its `SSLProviderCredentialsID`. Let's Encrypt needs no credentials.

## Requirements

- The organization needs the `Myra-Certificate` feature, which is enabled by Myra support.
- Every operation requires a member of the root user group of the organization or an organization administrator. Members of sub groups receive HTTP 403.

```go
type SSLProviderCredentials struct {
	ID         int             `json:"id,omitempty"`
	Created    *types.DateTime `json:"created,omitempty"`
	Modified   *types.DateTime `json:"modified,omitempty"`
	Name       string          `json:"name"`
	Provider   string          `json:"provider"`
	Cert       string          `json:"cert,omitempty"`
	PrivateKey string          `json:"privateKey,omitempty"`
	Email      string          `json:"email,omitempty"`
	Endpoint   string          `json:"endpoint,omitempty"`
	EABKid     string          `json:"eabKid,omitempty"`
	EABHmac    string          `json:"eabHmac,omitempty"`
	Comment    string          `json:"comment,omitempty"`
}
```

| Field | Type | Description |
|---|---|---|
| `ID` | int | ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete credentials it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. It will be created by the server after creating new credentials. This value is only informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates. A deletion identifies the credentials by their `ID` only. This value is always a date type with an `ISO 8601` format. |
| `Name` | string | An arbitrary, user-defined label for the credentials. Required. |
| `Provider` | string | The certificate provider the credentials belong to. Valid values: `SECTIGO`, `DTRUST`. Use the constants `myrasec.SSLProviderSectigo` and `myrasec.SSLProviderDTrust`. Required. |
| `Cert` | string | The public certificate of the ACME account key pair in PEM format. Leave `Cert` and `PrivateKey` empty on creation to let the server generate a pair. Supplying one of the two without the other is rejected. An empty value on update keeps the stored certificate. |
| `PrivateKey` | string | The private key of the ACME account key pair in PEM format. Write-only: the API never returns it. A value sent on update by a customer administrator is discarded. Contact Myra support to replace the pair. |
| `Email` | string | The contact email address of the provider account. Optional for `SECTIGO`. Must be empty for `DTRUST`, a non-empty value is rejected. |
| `Endpoint` | string | The ACME directory URL of the provider product, for example `https://acme.sectigo.com/v2/OV`. Required for `SECTIGO`, since the region and the validation product (DV, OV, EV) are encoded in the URL and tied to the EAB credentials. Optional for `DTRUST`, which falls back to the platform default. |
| `EABKid` | string | The External Account Binding key identifier issued by the provider. Required. |
| `EABHmac` | string | The External Account Binding HMAC key issued by the provider. Required on creation. Write-only: the API never returns it. An empty value on update keeps the stored key, a non-empty value rotates it. |
| `Comment` | string | An optional free-text note. |

## Create
To create new credentials send a `SSLProviderCredentials` without the attributes `ID`, `Created` and `Modified`. Leave `Cert` and `PrivateKey` empty to let the server generate the ACME account key pair. The response never contains `EABHmac` or `PrivateKey`.

### Example (Sectigo)
```go
credentials, err := api.CreateSSLProviderCredentials(&myrasec.SSLProviderCredentials{
    Name:     "Sectigo OV",
    Provider: myrasec.SSLProviderSectigo,
    Endpoint: "https://acme.sectigo.com/v2/OV",
    Email:    "pki@example.com",
    EABKid:   os.Getenv("SECTIGO_EAB_KID"),
    EABHmac:  os.Getenv("SECTIGO_EAB_HMAC"),
})
if err != nil {
    log.Fatal(err)
}

log.Println(credentials.ID)
```

### Example (D-Trust)
```go
credentials, err := api.CreateSSLProviderCredentials(&myrasec.SSLProviderCredentials{
    Name:     "D-Trust",
    Provider: myrasec.SSLProviderDTrust,
    EABKid:   os.Getenv("DTRUST_EAB_KID"),
    EABHmac:  os.Getenv("DTRUST_EAB_HMAC"),
})
if err != nil {
    log.Fatal(err)
}
```

## List
The listing operation returns the credentials of the organization.

### Example
```go
list, err := api.ListSSLProviderCredentials(map[string]string{"provider": "SECTIGO"})
if err != nil {
    log.Fatal(err)
}

for _, c := range list {
    log.Println(c.ID, c.Name, c.Provider)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListSSLProviderCredentials` function.

| Name | Description | Default |
|---|---|---|
| `provider` | Restrict the result to one provider: `SECTIGO` or `DTRUST`. | null |
| `search` | Restrict the result to credentials whose name, email or EAB key identifier contains the search string. | null |
| `page` | Specify the page of the result. | 1 |
| `pageSize` | Specify the amount of results in the response. | 50 |

## Read
The read operation returns a single credentials object by its `ID`.

### Example
```go
credentials, err := api.GetSSLProviderCredentials(credentialsId)
if err != nil {
    log.Fatal(err)
}
```

## Update
Updating credentials needs the `ID` and `Modified` attributes to identify the object and verify the version. Because the API never returns the secrets, an object read from the API has empty `EABHmac`, `PrivateKey` and, when it was server-generated, `Cert` attributes. Empty values are omitted from the request and keep the stored secrets. Set `EABHmac` to rotate the HMAC key.

### Example
```go
credentials, err := api.GetSSLProviderCredentials(credentialsId)
if err != nil {
    log.Fatal(err)
}

credentials.EABKid = os.Getenv("SECTIGO_EAB_KID")
credentials.EABHmac = os.Getenv("SECTIGO_EAB_HMAC")

updated, err := api.UpdateSSLProviderCredentials(credentials)
if err != nil {
    log.Fatal(err)
}
```

## Delete
Deleting credentials detaches them from every managed certificate request that references them: the requests keep existing and their `SSLProviderCredentialsID` is cleared. Certificates already issued keep being served, renewals of those requests need new credentials.

### Example
```go
_, err := api.DeleteSSLProviderCredentials(&myrasec.SSLProviderCredentials{ID: credentialsId})
if err != nil {
    log.Fatal(err)
}
```

## List certificates
`ListSSLProviderCertificates` returns the certificates issued with the credentials as compact summaries without PEM data, private key and intermediates. The full certificate is available through `ListSSLCertificates` on the domain given by `DomainID`, see [SSL certificates](./ssl.md).

```go
type SSLCertificateSummary struct {
	ID                   int             `json:"id,omitempty"`
	Created              *types.DateTime `json:"created,omitempty"`
	Modified             *types.DateTime `json:"modified,omitempty"`
	Subject              string          `json:"subject"`
	SubjectAlternatives  []string        `json:"subjectAlternatives"`
	Algorithm            string          `json:"algorithm"`
	ValidFrom            *types.DateTime `json:"validFrom"`
	ValidTo              *types.DateTime `json:"validTo"`
	Fingerprint          string          `json:"fingerprint"`
	SerialNumber         string          `json:"serialNumber"`
	Wildcard             bool            `json:"wildcard"`
	ExtendedValidation   bool            `json:"extendedValidation"`
	Managed              bool            `json:"managed"`
	Multidomain          bool            `json:"multidomain"`
	SslConfigurationName string          `json:"sslConfigurationName,omitempty"`
	RequestID            int             `json:"requestId,omitempty"`
	DomainID             int             `json:"domainId,omitempty"`
	Subdomains           []string        `json:"subdomains,omitempty"`
}
```

| Field | Type | Description |
|---|---|---|
| `ID` | int | The identifier of the certificate, the same as in `ListSSLCertificates`. |
| `Subject` | string | The subject of the certificate. |
| `SubjectAlternatives` | []string | The names covered by the certificate. |
| `Algorithm` | string | The signature algorithm of the certificate. |
| `ValidFrom` | *types.DateTime | Time when the certificate starts to be valid. |
| `ValidTo` | *types.DateTime | Time when the certificate expires. |
| `Fingerprint` | string | Fingerprint of the certificate. |
| `SerialNumber` | string | Serial number of the certificate. |
| `Wildcard` | bool | True when the certificate is a wildcard certificate. |
| `ExtendedValidation` | bool | True when the browser handles the certificate as extended validation. |
| `Managed` | bool | True when the certificate is managed and renewed by Myra. |
| `Multidomain` | bool | True when the certificate spans more than one domain. |
| `SslConfigurationName` | string | The TLS profile applied to the certificate. |
| `RequestID` | int | The `ID` of the managed certificate request the certificate was issued for. |
| `DomainID` | int | The `ID` of the domain the certificate belongs to. |
| `Subdomains` | []string | The subdomains the certificate is assigned to. |

### Example
```go
certs, err := api.ListSSLProviderCertificates(credentialsId, nil)
if err != nil {
    log.Fatal(err)
}

for _, c := range certs {
    log.Println(c.RequestID, c.Subject, c.ValidTo)
}
```

The function accepts the parameters `search`, `page` and `pageSize` like `ListSSLProviderCredentials`.

## Error handling
See [Error handling](./ssl_certificate_request.md#error-handling) of the certificate requests. In addition to the causes listed there, HTTP 403 is returned when the API user is neither a member of the root user group nor an organization administrator.
