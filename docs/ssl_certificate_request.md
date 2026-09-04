# SSL certificate request (Myra Managed Certificate)

A managed certificate request asks Myra to obtain a certificate from Let's Encrypt, Sectigo or D-Trust, to renew it automatically and to assign it to the subdomains of the request. The issuance is asynchronous: the request is created immediately, the certificate arrives later and the `Status` of the request reports the progress. Poll the request with `GetSSLCertificateRequestContext` until the status is `CREATED` or `FAILED`.

## Requirements

- The organization needs the `Myra-Certificate` feature, which is enabled by Myra support. Without it every call answers with HTTP 403, see [Error handling](#error-handling).
- Sectigo and D-Trust additionally need stored [SSL provider credentials](./ssl_provider_credentials.md). Let's Encrypt needs none.
- Every subject alternative name must belong to a domain registered in the Myra system. A request containing an unknown domain is rejected with HTTP 400.
- The certificate is only assigned to subdomains the API user may access. On a multi-domain request a user without access to one of the domains receives HTTP 403 on create.

```go
type SSLCertificateRequest struct {
	ID                       int                               `json:"id,omitempty"`
	Created                  *types.DateTime                   `json:"created,omitempty"`
	Modified                 *types.DateTime                   `json:"modified,omitempty"`
	Algorithm                string                            `json:"algorithm"`
	Provider                 string                            `json:"provider"`
	Status                   string                            `json:"status,omitempty"`
	FailureReason            string                            `json:"failureReason,omitempty"`
	CustomerActionable       bool                              `json:"customerActionable,omitempty"`
	SubjectAlternativeNames  []SSLCertificateRequestSAN        `json:"subjectAlternativeNames"`
	Assignments              []SSLCertificateRequestAssignment `json:"assignments"`
	MultiDomain              bool                              `json:"multiDomain,omitempty"`
	SSLProviderCredentialsID int                               `json:"sslProviderCredentialsId,omitempty"`
	RenewalInterval          int                               `json:"renewalInterval,omitempty"`
	SignatureAlgorithm       string                            `json:"signatureAlgorithm,omitempty"`
}
```

| Field | Type | Description |
|---|---|---|
| `ID` | int | ID is a unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a request it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created is a date type attribute with an `ISO 8601` format. It will be created by the server after creating a new request. This value is only informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates. A deletion identifies the request by its `ID` only. This value is always a date type with an `ISO 8601` format. |
| `Algorithm` | string | The key algorithm of the requested certificate. Valid values: `RSA2048`, `RSA4096`, `RSA8192`, `ECDSA256`, `ECDSA384`. Let's Encrypt accepts `RSA2048`, `ECDSA256` and `ECDSA384` only. The algorithm is immutable after creation and has to be sent unchanged on updates, a changed value is rejected with HTTP 400. |
| `Provider` | string | The certificate provider. Valid values: `LETS_ENCRYPT`, `SECTIGO`, `DTRUST`. Use the constants `myrasec.SSLProviderLetsEncrypt`, `myrasec.SSLProviderSectigo` and `myrasec.SSLProviderDTrust`. |
| `Status` | string | Reports the progress of the issuance, see [Status](#status). Read-only, the transitions happen on the server. |
| `FailureReason` | string | The reason code when the status is `FAILED`, see [Status](#status). Read-only. |
| `CustomerActionable` | bool | True when the failure reason is a DNS or CNAME issue the customer can resolve, false when Myra or the certificate authority has to act. Read-only. |
| `SubjectAlternativeNames` | []SSLCertificateRequestSAN | The names the certificate has to cover. Wildcards (`*.example.com`) are allowed as the leftmost label. Names made redundant by a wildcard in the same request are dropped by the server. |
| `Assignments` | []SSLCertificateRequestAssignment | The subdomains the issued certificate is assigned to. Each subdomain has to exist as a DNS record in the Myra system. |
| `MultiDomain` | bool | True when the subject alternative names span more than one domain. Read-only. |
| `SSLProviderCredentialsID` | int | The `ID` of the [SSL provider credentials](./ssl_provider_credentials.md) used for the issuance. Required for `SECTIGO` and `DTRUST`, ignored for `LETS_ENCRYPT`. The server clears the value when the referenced credentials are deleted. |
| `RenewalInterval` | int | The number of days before the certificate expires at which it is renewed. Zero means the system default. Accepted for `SECTIGO` and `DTRUST` only. |
| `SignatureAlgorithm` | string | The signature algorithm of the requested certificate. Valid values: `SHA256`, `SHA384`, `SHA512`. Empty means the system default. Accepted for `SECTIGO` and `DTRUST` only. `SHA512` cannot be combined with `ECDSA256` or `ECDSA384`. |

```go
type SSLCertificateRequestSAN struct {
	ID       int             `json:"id,omitempty"`
	Created  *types.DateTime `json:"created,omitempty"`
	Modified *types.DateTime `json:"modified,omitempty"`
	Name     string          `json:"name"`
}

type SSLCertificateRequestAssignment struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	SubDomainName string          `json:"subDomainName"`
}
```

| Field | Type | Description |
|---|---|---|
| `Name` | string | The FQDN the certificate has to cover, for example `www.example.com` or `*.example.com`. |
| `SubDomainName` | string | The FQDN of the subdomain the issued certificate is assigned to. |

The `ID`, `Created` and `Modified` attributes of both nested objects are server-generated. On updates, send the `ID` of a subject alternative name back to keep its stored entry, an entry without `ID` replaces the stored entry of the same name. Assignments are matched by `SubDomainName`, their `ID` is ignored.

## Status

| Value | Constant | Meaning |
|---|---|---|
| `OPEN` | `SSLCertificateRequestStatusOpen` | The request has been accepted and the certificate is being issued. |
| `WAITING_FOR_CNAME` | `SSLCertificateRequestStatusWaitingForCNAME` | The issuance waits for the domain validation. Use [Check domains](#check-domains) to find out which CNAME record is expected. |
| `CREATED` | `SSLCertificateRequestStatusCreated` | The certificate has been issued and assigned to the subdomains of the request. The certificate itself is available through `ListSSLCertificatesContext` on the domain (with `Managed` set to true) or, for Sectigo and D-Trust, through `ListSSLProviderCertificatesContext`. |
| `FAILED` | `SSLCertificateRequestStatusFailed` | The issuance did not succeed. The request is not retried automatically. Fix the cause reported in `FailureReason`, then create a new request. Only two updates re-enter the issuance of a failed request: adding a name that no issued certificate covers or changing the `Provider`. |

| Failure reason | Customer actionable | Meaning |
|---|---|---|
| `CNAME_TIMEOUT` | yes | The expected CNAME record for the domain validation did not appear in time. |
| `VALIDATION_FAILED` | yes | The domain validation failed, usually because of the DNS setup of the domain. |
| `VALIDATION_ERROR` | no | The validation stage failed on the Myra or CA side, for example a rate limit or a lost connection. |
| `ORDER_FAILED` | no | The certificate authority rejected the order. |
| `CERT_LOAD_FAILED` | no | The issued certificate could not be loaded into the platform. |

## Create
To create a new managed certificate request send a `SSLCertificateRequest` without the attributes `ID`, `Created` and `Modified`. The response carries the persisted request in status `OPEN`. The certificate is issued asynchronously, see [End-to-end example](#end-to-end-example-lets-encrypt).

### Example (Let's Encrypt)
```go
request := &myrasec.SSLCertificateRequest{
    Provider:  myrasec.SSLProviderLetsEncrypt,
    Algorithm: myrasec.SSLCertificateRequestAlgorithmRSA2048,
    SubjectAlternativeNames: []myrasec.SSLCertificateRequestSAN{
        {Name: "www.example.com"},
    },
    Assignments: []myrasec.SSLCertificateRequestAssignment{
        {SubDomainName: "www.example.com"},
    },
}

created, err := api.CreateSSLCertificateRequestContext(ctx, request)
if err != nil {
    log.Fatal(err)
}

log.Println(created.ID, created.Status)
```

### Example (Sectigo or D-Trust)
Sectigo and D-Trust need the `ID` of stored [SSL provider credentials](./ssl_provider_credentials.md). They additionally accept `RSA4096`, `RSA8192`, a `RenewalInterval` and a `SignatureAlgorithm`.

```go
request := &myrasec.SSLCertificateRequest{
    Provider:                 myrasec.SSLProviderSectigo,
    Algorithm:                myrasec.SSLCertificateRequestAlgorithmRSA4096,
    SSLProviderCredentialsID: credentials.ID,
    RenewalInterval:          30,
    SignatureAlgorithm:       myrasec.SSLCertificateRequestSignatureAlgorithmSHA384,
    SubjectAlternativeNames: []myrasec.SSLCertificateRequestSAN{
        {Name: "www.example.com"},
        {Name: "api.example.com"},
    },
    Assignments: []myrasec.SSLCertificateRequestAssignment{
        {SubDomainName: "www.example.com"},
        {SubDomainName: "api.example.com"},
    },
}

created, err := api.CreateSSLCertificateRequestContext(ctx, request)
if err != nil {
    log.Fatal(err)
}
```

## List
The listing operation returns the managed certificate requests of the organization the API user may read.

**NOTE: Without a `status` parameter the API returns only requests in status `OPEN`, `WAITING_FOR_CNAME` and `FAILED`. Requests whose certificate has been issued (`CREATED`) are omitted. Pass `status` explicitly to see them.**

### Example
```go
requests, err := api.ListSSLCertificateRequestsContext(ctx, map[string]string{
    "status": "OPEN,WAITING_FOR_CNAME,CREATED,FAILED",
    "domain": "example.com",
})
if err != nil {
    log.Fatal(err)
}

for _, r := range requests {
    log.Println(r.ID, r.Provider, r.Status)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListSSLCertificateRequestsContext` function.

| Name | Description | Default |
|---|---|---|
| `status` | Comma separated list of statuses to include: `OPEN`, `WAITING_FOR_CNAME`, `CREATED`, `FAILED`. | `OPEN,WAITING_FOR_CNAME,FAILED` |
| `domain` | Restrict the result to requests of one domain. Accepts the domain ID or the domain name. | null |
| `search` | Restrict the result to requests with an assigned subdomain whose name contains the search string. Requests without assignments are omitted when `search` is set. | null |
| `page` | Specify the page of the result. | 1 |
| `pageSize` | Specify the amount of results in the response. | 50 |

## Read
The read operation returns a single managed certificate request by its `ID`. Use it to poll the `Status`.

### Example
```go
request, err := api.GetSSLCertificateRequestContext(ctx, requestId)
if err != nil {
    log.Fatal(err)
}

log.Println(request.Status, request.FailureReason)
```

## Update
Updating a request is very similar to creating a new one. You have to provide the `ID` and `Modified` attributes to identify the version of the object and the unchanged `Algorithm`, since it is immutable. Send the full lists of `SubjectAlternativeNames` and `Assignments`: entries missing from the lists are removed.

- Adding a name that no issued certificate of the request covers triggers the issuance of a new certificate. The current certificate keeps being served until the new one arrives.
- Removing a name depends on the provider. For Let's Encrypt the issued certificate keeps covering the removed name until its next renewal, for Sectigo and D-Trust a narrowed certificate is re-issued.
- On a `FAILED` request only adding a name that no issued certificate covers or changing the `Provider` re-enters the issuance. Removing a name changes nothing there.

### Example
```go
request, err := api.GetSSLCertificateRequestContext(ctx, requestId)
if err != nil {
    log.Fatal(err)
}

request.SubjectAlternativeNames = append(request.SubjectAlternativeNames, myrasec.SSLCertificateRequestSAN{Name: "shop.example.com"})
request.Assignments = append(request.Assignments, myrasec.SSLCertificateRequestAssignment{SubDomainName: "shop.example.com"})

updated, err := api.UpdateSSLCertificateRequestContext(ctx, request)
if err != nil {
    log.Fatal(err)
}
```

## Delete
Deleting a request removes the request and the certificates issued for it. The `ID` attribute identifies the request.

### Example
```go
_, err := api.DeleteSSLCertificateRequestContext(ctx, &myrasec.SSLCertificateRequest{ID: requestId})
if err != nil {
    log.Fatal(err)
}
```

## Update the SSL configuration
`UpdateSSLCertificateRequestConfigurationContext` applies a TLS profile to every certificate issued for the request. Valid profile names are returned by `ListSslConfigurationsContext`, see [SSL configuration](./ssl_configuration.md). Profiles restricted to ECDSA keys are rejected for RSA requests.

### Example
```go
request, err := api.UpdateSSLCertificateRequestConfigurationContext(ctx, requestId, "2023-mozilla-modern")
if err != nil {
    log.Fatal(err)
}
```

## Check domains
`CheckSSLCertificateRequestDomainsContext` runs a live name server check for up to 99 domains and reports per domain whether it is served by Myra name servers. A domain served elsewhere needs a CNAME record from `ChallengeName` to `ExpectedCName` before the domain validation can succeed. Run the check before creating a request or when a request stays in `WAITING_FOR_CNAME`.

```go
type SSLCertificateRequestDomainCheck struct {
	Exists        bool   `json:"exists"`
	IsMyraNS      bool   `json:"isMyraNS"`
	ChallengeName string `json:"challengeName,omitempty"`
	ExpectedCName string `json:"expectedCName,omitempty"`
}
```

| Field | Type | Description |
|---|---|---|
| `Exists` | bool | True when the domain has name server records. |
| `IsMyraNS` | bool | True when the domain is served by Myra name servers. No CNAME record is needed then. |
| `ChallengeName` | string | The record name that has to be created as CNAME for the domain validation. Empty when no record is needed. |
| `ExpectedCName` | string | The target the challenge record has to point to. Empty when no record is needed. |

### Example
```go
checks, err := api.CheckSSLCertificateRequestDomainsContext(ctx, []string{"www.example.com", "*.example.org"})
if err != nil {
    log.Fatal(err)
}

for name, check := range checks {
    if !check.IsMyraNS {
        log.Printf("%s: create CNAME %s -> %s", name, check.ChallengeName, check.ExpectedCName)
    }
}
```

Domains whose lookup failed are absent from the result.

## End-to-end example (Let's Encrypt)
```go
// Bound the whole flow, including the polling below, to two hours.
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
defer cancel()

request, err := api.CreateSSLCertificateRequestContext(ctx, &myrasec.SSLCertificateRequest{
    Provider:                myrasec.SSLProviderLetsEncrypt,
    Algorithm:               myrasec.SSLCertificateRequestAlgorithmECDSA256,
    SubjectAlternativeNames: []myrasec.SSLCertificateRequestSAN{{Name: "www.example.com"}},
    Assignments:             []myrasec.SSLCertificateRequestAssignment{{SubDomainName: "www.example.com"}},
})
if err != nil {
    log.Fatal(err)
}

// The issuance takes minutes and can take hours when a CNAME record has to be created.
// Poll with a generous interval, every request counts against the API rate limit.
for {
    request, err = api.GetSSLCertificateRequestContext(ctx, request.ID)
    if err != nil {
        log.Fatal(err)
    }

    switch request.Status {
    case myrasec.SSLCertificateRequestStatusCreated:
        certs, err := api.ListSSLCertificatesContext(ctx, domainId, nil)
        if err != nil {
            log.Fatal(err)
        }
        for _, c := range certs {
            if c.Managed {
                log.Println("issued:", c.Subject, c.ValidTo)
            }
        }
        return
    case myrasec.SSLCertificateRequestStatusFailed:
        log.Fatalf("issuance failed: %s (customer actionable: %t)", request.FailureReason, request.CustomerActionable)
    case myrasec.SSLCertificateRequestStatusWaitingForCNAME:
        checks, _ := api.CheckSSLCertificateRequestDomainsContext(ctx, []string{"www.example.com"})
        log.Printf("waiting for CNAME %s -> %s", checks["www.example.com"].ChallengeName, checks["www.example.com"].ExpectedCName)
    }

    select {
    case <-ctx.Done():
        log.Fatal(ctx.Err())
    case <-time.After(time.Minute):
    }
}
```

## End-to-end example (Sectigo or D-Trust)
```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
defer cancel()

credentials, err := api.CreateSSLProviderCredentialsContext(ctx, &myrasec.SSLProviderCredentials{
    Name:     "Sectigo OV",
    Provider: myrasec.SSLProviderSectigo,
    Endpoint: "https://acme.sectigo.com/v2/OV",
    EABKid:   os.Getenv("SECTIGO_EAB_KID"),
    EABHmac:  os.Getenv("SECTIGO_EAB_HMAC"),
})
if err != nil {
    log.Fatal(err)
}

request, err := api.CreateSSLCertificateRequestContext(ctx, &myrasec.SSLCertificateRequest{
    Provider:                 myrasec.SSLProviderSectigo,
    Algorithm:                myrasec.SSLCertificateRequestAlgorithmRSA2048,
    SSLProviderCredentialsID: credentials.ID,
    SubjectAlternativeNames:  []myrasec.SSLCertificateRequestSAN{{Name: "www.example.com"}},
    Assignments:              []myrasec.SSLCertificateRequestAssignment{{SubDomainName: "www.example.com"}},
})
if err != nil {
    log.Fatal(err)
}

// Poll GetSSLCertificateRequestContext as in the Let's Encrypt example. Once the status is
// CREATED, the certificates issued with the credentials are listed on the credentials.
certs, err := api.ListSSLProviderCertificatesContext(ctx, credentials.ID, nil)
if err != nil {
    log.Fatal(err)
}
for _, c := range certs {
    log.Println(c.RequestID, c.Subject, c.ValidTo)
}
```

## Rate limits
The issuance is subject to the following limits. They are checked before the certificate is ordered and apply to every provider, not only to Let's Encrypt.

| Limit | Description |
|---|---|
| 50 certificates per domain per week | Counted per registered domain, so `www.example.com` and `api.example.com` count against the limit of `example.com`. |
| 300 orders per account per 3 hours | Counted per certificate provider account. |
| 5 certificates per week for the same set of names | Repeatedly re-issuing an unchanged request reaches this limit first. |
| 99 subject alternative names per request | A request with 100 names or more is rejected before the certificate is ordered. The create call itself accepts up to 100 names. |

A request that reaches a limit does not wait for the window to free up. For a first issuance the status becomes `FAILED` with the failure reason `VALIDATION_ERROR`, create a new request once the window has passed. For the renewal of an already issued certificate the request stays `CREATED` and the next renewal run retries, the served certificate is not affected.

Independent of the issuance limits, the client limits API calls to 5 per second. Poll the request status with an interval of a minute or more.

## Error handling
Every method returns a `*myrasec.APIError` when the API answers with a non-successful HTTP status. Inspect it with `errors.As`:

| Status | Cause |
|---|---|
| 400 | Validation failed. `Violations` names the attribute and the reason, for example an unknown domain in `subjectAlternativeNames` or an algorithm the provider does not support. A stale `Modified` attribute on an update is reported the same way, with the violation message `api.error.already-modified`: read the request again and retry. |
| 403 | The organization lacks the `Myra-Certificate` feature or the user lacks the `SslCertRequest` permission on one of the domains. The API sends no body in this case, so the error carries the status code only. A request that does not exist answers 403 as well, with a violation `does-not-exist`. |

```go
request, err := api.GetSSLCertificateRequestContext(ctx, requestId)
if err != nil {
    var apiErr *myrasec.APIError
    if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden {
        log.Fatal("Myra Managed Certificates are not enabled for this organization or user")
    }
    log.Fatal(err)
}
```

The client retries HTTP 500 responses when `SetMaxRetries` is set, all other statuses are returned immediately. See [Setup](./setup.md#error-handling) for the general error handling of the client.
