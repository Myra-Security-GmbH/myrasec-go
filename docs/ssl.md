# SSL certificate

```go
type Certificate struct {
	ID           int             `json:"id,omitempty"`
	Created      *types.DateTime `json:"created,omitempty"`
	Modified     *types.DateTime `json:"modified,omitempty"`
	Subject      string          `json:"subject"`
	Algorithm    string          `json:"algorithm"`
	ValidFrom    *types.DateTime `json:"validFrom"`
	ValidTo      *types.DateTime `json:"validTo"`
	Fingerprint  string          `json:"fingerprint"`
	SerialNumber string          `json:"serialNumber"`
	Cert         string          `json:"cert,omitempty"`
}
```

| Field | Type | Description |  
|---|---|---|  
| `ID` | int | ID is a unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a Certificate it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. It will be created by the server after creating a new Certificate object. This value is only informational so it is not necessary to add this an attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletions. This value is always a date type with an `ISO 8601` format. |  
| `Subject` | string | Shows the subject of the uploaded certificate. |  
| `Algorithm` | string | Contains the signature algorithm. |  
| `ValidFrom` | *types.DateTime | Time when the certificate starts to be valid. This property is a date type with an `ISO 8601` format. |  
| `ValidTo` | *types.DateTime | Time when the certificate expires. This property is a date type with an `ISO 8601` format. |  
| `Fingerprint` | string | Fingerprint of the certificate. |  
| `SerialNumber` | string | Serial number of the certificate. |  
| `Cert` | string | Cert contains the certificate. |  


```go
type SSLCertificate struct {
	*Certificate
	SubjectAlternatives  []string          `json:"subjectAlternatives"`
	Intermediates        []SSLIntermediate `json:"intermediates,omitempty"`
	Wildcard             bool              `json:"wildcard"`
	ExtendedValidation   bool              `json:"extendedValidation"`
	Subdomains           []string          `json:"subdomains,omitempty"`
	Key                  string            `json:"key,omitempty"`
	CertRefreshForced    bool              `json:"certRefreshForced,omitempty"`
	CertToRefresh        int               `json:"certToRefresh,omitempty"`
	SslConfigurationName string            `json:"sslConfigurationName,omitempty"`
    Managed              bool              `json:"managed"`
}
```

| Field | Type | Description |  
|---|---|---|  
| `SubjectAlternatives` | string | Contains a list of subdomains which can be validated using this certificate. This list also contains the CN of the subject. |  
| `Intermediates` | []SSLIntermediate | The intermediate certificates of the chain of trust. They are sent on create and update and validated when the certificate is saved (see [Validation: errors and warnings](#validation-errors-and-warnings)). A broken intermediate rejects the save, e.g. one with the wrong key or an invalid signature (`INTERMEDIATE_WRONG_KEY`, `INTERMEDIATE_SIGNATURE_INVALID`), one that is not a CA or one that has expired. A missing intermediate is a warning (`CHAIN_INTERMEDIATE_MISSING`). Certificates that do not belong to the chain are discarded with a warning (`CHAIN_UNRELATED_CERTIFICATES`) and the save goes through; only the certificates of the chain are stored. The API sorts them by subject and issuer. |  
| `Wildcard` | bool | This property shows whether the certificate is valid for multiple subdomains of a domain. The certificate needs to have a *.domain.tld subject to return true. |  
| `ExtendedValidation` | bool | True if the browser handles the certificate as extended validation. We use the OIDs from Google Chrome™ to measure the extended validation level. |  
| `Subdomains` | []string | A list of subdomains assigned to this certificate. |  
| `Key` | string | The unencrypted private key that matches your certificate. |  
| `CertRefreshForced` | bool | Every time a certificate is refreshed with another non-matching certificate the operation is interrupted with an error. Setting certRefreshForced will ignore such errors and refresh the certificate anyway. Please use it only, if you are sure you can ignore an error when refreshing a certificate. |  
| `CertToRefresh` | int | This property allows you to update an already existing certificate with a new one without changing IP addresses, the value has to be the ID of the cert that should be refreshed. |
| `SslConfigurationName` | string | This property allows you to set a specific ssl configuration. default `Myra-Global-TLS-Default`, valid values are `Myra-Global-TLS-Default`, `2023-mozilla-intermediate`, `2023-mozilla-modern`. |
| `Managed` | bool | Indicates whether this certificate is managed by Myra or not. Managed certificates are issued and renewed through [SSL certificate requests](./ssl_certificate_request.md). |


```go
type SSLIntermediate struct {
	*Certificate
	Issuer string `json:"issuer"`
}
```

| Field | Type | Description |  
|---|---|---|  
| `Issuer` | string |  |  


## Create
To upload and assign a new certificate send a SSLCertificate object without the attributes "id", "created" and "modified". All attributes are generated by the server and will be returned after that.

### Example
```go
crt := `
-----BEGIN CERTIFICATE-----
MIIB5jCCAUegAwIBAgIBATAKB........
.................................
.................................
.................................
......AQ==
-----END CERTIFICATE-----```
key := `
-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIAcIS..................
.................................
.................................
......oO1fcQw==
-----END EC PRIVATE KEY-----```

ssl := &myrasec.SSLCertificate{
    Certificate: &myrasec.Certificate{
        Cert: crt,
    },
    Key: key,
    Subdomains: []string{
        "www.example.com",
    },
}
cert, err := api.CreateSSLCertificateContext(ctx, ssl, "example.com")
if err != nil {
    log.Fatal(err)
}
```

### Validation: errors and warnings
The API checks every uploaded certificate, key and chain when it is saved (create, renewal with `certToRefresh` and an update of the intermediates). The SDK relies on this save-time check; it has no separate preview call.

- A blocking result rejects the request and nothing is stored. The error is an `*APIError` whose `Violations` carry, per result, the field (`Path`, e.g. `cert`, `key`, `intermediates[0].cert`), the stable `Code` (e.g. `LEAF_EXPIRED`), the `Translations` (`en`, `de`) and the `Parameters` of the message. `err.Error()` renders each as `path: message [CODE]`. The warnings of the same check come with it in `Warnings`, so a rejected save shows every result at once.
- A warning (e.g. `CHAIN_INTERMEDIATE_MISSING`, `CHAIN_UNRELATED_CERTIFICATES`, `LEAF_EXPIRES_SOON`) does not fail the request. `CreateSSLCertificateWithWarningsContext` and `UpdateSSLCertificateWithWarningsContext` return the warnings next to the saved certificate; `CreateSSLCertificateContext` and `UpdateSSLCertificateContext` drop them. Warnings come only with the response of the save that triggered them; reading a certificate (`GetSSLCertificateContext`, `ListSSLCertificatesContext`) never returns them.

```go
cert, warnings, err := api.CreateSSLCertificateWithWarningsContext(ctx, ssl, domainId)
if err != nil {
    var apiErr *myrasec.APIError
    if errors.As(err, &apiErr) {
        for _, v := range apiErr.Violations {
            log.Printf("%s: %s [%s]", v.Path, v.Message, v.Code)
        }
        for _, w := range apiErr.Warnings {
            log.Printf("warning: %s", w)
        }
    }
    log.Fatal(err)
}
for _, w := range warnings {
    log.Printf("warning: %s", w)
}
```

## List
The listing operation returns a list of SSL certificates for the passed domainId.

### Example
```go
api.ListSSLCertificatesContext(ctx, domainId, nil)
if err != nil {
    log.Fatal(err)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListSSLCertificatesContext` function.

| Name | Description | Default |
|---|---|---|
| `page` | Specify the page of the result | 1 |
| `pageSize` | Specify the amount of results in the response | 50 |

## Read
The read operation returns a single SSL certificate by it's ID
```go
cert, err := api.GetSSLCertificateContext(ctx, domaindId, certId)
if err != nil {
    log.Fatal(err)
}
```

## Update
Updating a SSL certificate is very similar to creating a new one. You will need to provide the generated "id" and "modified" attributes to identify the version of object you are trying to update.  
Setting the certToRefresh, will update the existing certificate without changing its assigned IPs.  
Every time a certificate is refreshed with another non-matching certificate the operation is interrupted with an error. Setting certRefreshForced will ignore such errors and refresh the certificate
anyway. Please use it only, if you are sure you can ignore an error when refreshing a certificate

### Example
```go
ssl := &myrasec.SSLCertificate{
    ID:   0000,
    Modified: &types.DateTime{
        Time: modified,
    },
    SubDomains: []string{
        "www.example.com", 
        "example.com",
    },
}

cert, err := api.UpdateSSLCertificateContext(ctx, ssl, "example.com")
if err != nil {
    log.Fatal(err)
}
```

## Delete
It is not necessary to delete a certificate. When a certificate has not been assigned to any subdomain for more than one week it will be automatically removed.
Due to the sake of consistency and completeness, the `DeleteSSLCertificateContext` function is implemented and simply removes the assigned subdomains from the certificate.

### Example
```go
ssl := &myrasec.SSLCertificate{
    ID:   0000,
    Modified: &types.DateTime{
        Time: modified,
    },
    SubDomains: []string{
        "www.example.com", 
        "example.com",
    },
}

cert, err := api.DeleteSSLCertificateContext(ctx, ssl, "example.com")
if err != nil {
    log.Fatal(err)
}

log.Println(cert.SubDomains)
```