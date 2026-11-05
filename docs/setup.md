# Setup myrasec-go

To be able to access the Myra API using the myrasec-go lib, you need a API Key and API Secret or an API token. You can create this on your own after accessing your own user page in the user management. Once you have access, the library can manage them as well, see [API Key](./apikey.md) and [API Token](./apitoken.md).

## Setup example
```go
api, err := myrasec.New(os.Getenv("MYRA_API_KEY"), os.Getenv("MYRA_API_SECRET"))
if err != nil {
    log.Fatal(err)
}
```

Using the API token for authentication:
```go
api, err := myrasec.NewWithToken(os.Getenv("MYRA_API_TOKEN"))
if err != nil {
    log.Fatal(err)
}
```

## List 100 domains example
```go
package main

import (
	"context"
	"log"
	"os"

	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
)

func main() {
	api, err := myrasec.New(os.Getenv("MYRA_API_KEY"), os.Getenv("MYRA_API_SECRET"))
	if err != nil {
		log.Fatal(err)
	}
	
	ctx := context.Background()

	domains, err := api.ListDomainsContext(ctx, map[string]string{"pageSize": "100"})
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range domains {
		log.Println(d.Name)
	}
}
```

## Cancellation and timeouts
Every function that sends a request has a variant with the suffix `Context` that takes a `context.Context` as its first parameter, for example `ListDomainsContext(ctx, params)`. The context bounds the whole call: the wait for the client side rate limiter, the HTTP request and the sleep between retries. When the context is cancelled or its deadline passes, the function returns `context.Canceled` or `context.DeadlineExceeded` (check with `errors.Is`) and no further request is sent. A deadline that is shorter than the wait the client side rate limiter needs for the next request fails immediately with `context.DeadlineExceeded` as well, the message names the rate limit as the cause. Independent of the context, the HTTP client times out after 30 seconds per request; use `SetHTTPClient` to change that.

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

domains, err := api.ListDomainsContext(ctx, map[string]string{"pageSize": "100"})
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Fatal("the API did not answer within 10 seconds")
    }
    log.Fatal(err)
}
```

The functions without the suffix, for example `ListDomains(params)`, are equivalent to the `Context` variant called with `context.Background()`. They are deprecated and will be removed in the next major version. Functions added since then exist as `Context` variants only. The examples in this documentation use the `Context` variants and assume a `ctx` variable, for example `ctx := context.Background()`.

## Error handling
Every function returns a `*myrasec.APIError` when the API answers with a non-successful HTTP status. The error message contains the status and, when the API sent them, the violations and the error message of the response. Use `errors.As` to inspect the status code and the violations:

```go
domain, err := api.GetDomainContext(ctx, domainId)
if err != nil {
    var apiErr *myrasec.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.StatusCode {
        case http.StatusForbidden:
            log.Fatal("missing permission or feature")
        case http.StatusBadRequest:
            for _, v := range apiErr.Violations {
                log.Println(v.Path, v.Message)
            }
        }
    }
    log.Fatal(err)
}
```

| Field | Type | Description |
|---|---|---|
| `StatusCode` | int | The HTTP status code of the response. |
| `ErrorMessage` | string | The `errorMessage` attribute of the response, empty when the API sent none. |
| `Violations` | []*Violation | The `violationList` entries of the response with `Path` and `Message`, empty when the API sent none. Access denied responses (403) carry no body, so both are empty then. |
| `Body` | string | The raw body of a non-successful response. Useful when the body is not a JSON error envelope, for example HTML from a proxy. In that case the error message carries the status only. |

A successful status with an error envelope (`error: true`, used by some endpoints for validation failures) is returned as `*APIError` as well. Its `StatusCode` is the successful status and its message carries the violations and the error message without a status prefix, exactly the message the library returned before `APIError` existed. Errors that happen before a response arrives (network errors, the client side rate limit, an unsupported action) are returned as plain errors.

## Retries
The client retries a request when `SetMaxRetries` is set to a value above 1 (the default is 1, a single attempt) and the API answers with one of these statuses:

| Status | Retried for |
|---|---|
| 429 Too Many Requests | every method, a rate limiter answers before the request is processed |
| 503 Service Unavailable | GET, and every method when the response carries `Retry-After`, the explicit invitation of a server that did not handle the request. Without the header a 503 can come from a gateway that lost the API after passing the request on, or from the API after a partial change, so writes are not repeated. |
| 500 Internal Server Error, 502 Bad Gateway, 504 Gateway Timeout | GET only. The API or the gateway in front of it may have processed the request already: a repeated POST could create a duplicate, a repeated PUT or DELETE would fail on the changed or missing entry. |

The wait between the attempts is the `Retry-After` header of the response when the API sends one (capped at 60 seconds). Without the header the wait starts at `SetRetrySleep` seconds (default 0) and doubles with every further attempt (for example 1s, 2s, 4s) up to 60 seconds, or up to the configured sleep when that is longer. A 429 or 503 without `Retry-After` waits at least one second. The wait is interrupted when the context is cancelled. When the wait, together with the wait for the client side rate limit, would reach past the deadline of the context, the client does not wait and returns the error of the last response instead, so the status the API sent is not lost behind a deadline error.

## Transport security
The client refuses to send requests to a plain `http://` endpoint on a host other than the local machine, because the API key and signature or the API token would travel in clear text, and fails with `ErrInsecureEndpoint` (check with `errors.Is`). Endpoints on `localhost` or a loopback address are exempt, so local test servers work without any setting. To talk to a plain HTTP endpoint on another host, for example a proxy in a lab network, opt in explicitly:

```go
api.SetInsecureHTTP(true)
```

The API does not redirect, so a redirect means a misconfigured endpoint or a hostile network. The HTTP client refuses redirects to another host and redirects from HTTPS to plain HTTP (the Go HTTP client keeps the credentials on a same-host redirect regardless of the scheme) and follows at most 10 same-host redirects. A client passed to `SetHTTPClient` without its own `CheckRedirect` is copied and the copy gets this policy; a client with its own `CheckRedirect` is used unchanged.
