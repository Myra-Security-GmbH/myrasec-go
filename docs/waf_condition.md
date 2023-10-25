# WAF condition

## List
The listing operation returns a list of WAF conditions.

### Example
```go
rules, err := api.ListWAFConditions()
if err != nil {
    log.Fatal(err)
}
```

## Condition types
| Condition display name | Name | Description | Phase (Request/in or Response/out) | Key-Value pair |
| - | - | - | - | - |
| Custom header | custom_header | The "Custom header" condition allows you to specify a freely definable key value pair for a header for matching. | in/out | Yes |
| Host header | host | You can match to a specific value within the host header. | in | Only value |
| User-Agent header | user_agent | Here you can match to the content of the User-Agent header. | in | Only value |
| Accept header | accept | The "Accept header" condition allows you to match to the requested Content-Type header (MIME type). | in | Only value |
| Accept-Encoding header | accept_encoding | The "Accept-Encoding header" condition allows you to match to the requested compression method. | in | Only value |
| Cookie | cookie | You can specify a freely definable key value pair for a cookie for matching. | in | Yes |
| Path | url | Here you can match a freely definable path. | in | Only value |
| Request method | method | The "Request method" condition allows you to match to the HTTP method used. | in | Only valuen|
| Query string argument | arg | The "Query string argument" condition allows you to specify a freely definable key value pair for matching. | in | Yes |
| Post argument | postarg | The "Post argument" condition allows you to match to a key value pair of a POST. | in | Yes |
| Query string | querystring | Here you can match to freely definable query strings. | in | Only value |
| Fingerprint | fingerprint | The "Fingerprint" condition allows you to define geo-blocking, AS blocking (Autonomous System), or Managed Myra fingerprint blocking. | in | Only value |
| Remote IP address | remote_addr | The "Remote IP address" condition allows you to match to the client IP address. This works with both IPv4 and IPv6 addresses. | in | Only value |
| Score | score | The "Score" condition allows you to match to points assigned by the scoring rules. | in | Only value |
| Set-Cookie header | set_cookie | The "Set-Cookie header" is used by web servers to send cookies to the browser by Nginx. | out | Only value |
| Content-Type header | content_type | The "Content-Type header" is used to filter the content allowd to access the website. | out | Only value |
