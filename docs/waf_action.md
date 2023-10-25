# WAF action

## List
The listing operation returns a list of WAF actions.

### Example
```go
rules, err := api.ListWAFActions()
if err != nil {
    log.Fatal(err)
}
```

## Action types
| Action display name | Name | Description | Direction/Phase (in/request or out/response) | CustomKey & Value pair | Valid Keys | Valid Values |
| --- | --- | --- | --- | --- | --- | --- |
| Allow | allow | The request is answered without processing by the request limiter and antibot mechanism. The WAF status "AW" (Allowed WAF) is set in the access log. | in | No | - | - |
| Block | block | The requset is blocked. The WAF status "BW" (Blocked WAF) is set in the access log. | in | No | - | - |
| Log | log | The request is highlighted in the access log. In the log, the request is marked as "LW" (Logged WAF). This allows you to test rules. | in | No | - | - |
| Modify header | modify_header | This action allows you to overwrite a header added to the request | in/out | Yes | Any | Any |
| Add header | add_header | This action allows you to add a header to the request | in/out | Yes | Any | Any |
| Remove header | remove_header | This action allows you to remove headers from the request | in/out | Only value | - | Any |
| CAPTCHA | verify_human | The request is answered with a CAPTCHA | in | No | - | - |
| Change upstream | change_upstream | The request is sent to another upstream, similar to a ProxyPass. The requested URL remains unchanged. The changed upsream must also be stored as a DNS record at Myra. | in/out | Only value | - | Any |
| Rate limit | origin_rate_limit | The "Rate limit" action allows you to set a request limit that is independent from the domain/subdomain. The first value determines the time period in seconds, the second value the number of allowed requests. When the set request limit is reached, all further requests are blocked until the end of the defined time period. In the next time period, requests are allowd again until the limit is exeeded once more. | in/out | Yes | 1,2,5,10,15,30,45,60,120, 180, 300, 600, 1200, 3600 | >= 1 |
| Score | score | The "Score" action allows you to assign points to parameters of a request. The scrore of a request is 0 by default and can be incremented, decremented, or multiplied. To do this, enter one of the operators `+`, `-`, or `*` in the Name field. In addition, you mut specify the operand in the Value field. You can set to 0 by multiplying by 0, which is also the default score. This is useful in combination with other rules, since this value can then be considered again in subsequent WAF rules. | in/out | Yes | `+`, `-`, `*` | >= 0 |
| URI subscription | uri_subst | This action allows you to replace parts of the requested URI. Only the first match in the URI is replaced. | in/out | Yes | Any | Any |
| Set HTTP status | set_http_status | The 'SetHTTP status` action allows you to specify the request to be answered with an HTTP status. | in/out | Yes | 301, 302, 404 | Any |