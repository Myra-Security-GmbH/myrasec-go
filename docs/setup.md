# Setup myrasec-go

To be able to access the Myra API using the myrasec-go lib, you need a API Key and API Secret or an API token. You can create this on your own after accessing your own user page in the user management.

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

The functions without the suffix, for example `ListDomains(params)`, are equivalent to the `Context` variant called with `context.Background()`. They are deprecated and will be removed in the next major version. The examples in this documentation use the `Context` variants and assume a `ctx` variable, for example `ctx := context.Background()`.

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

A successful status with an error envelope (`error: true`, used by some endpoints for validation failures) is returned as `*APIError` as well. Its `StatusCode` is the successful status and its message carries the violations and the error message without a status prefix, exactly the message the library returned before `APIError` existed. Errors that happen before a response arrives (network errors, the client side rate limit, an unsupported action) are returned as plain errors. The client retries HTTP 500 responses when `SetMaxRetries` is set to a value above 1, waiting `SetRetrySleep` seconds between the attempts.
