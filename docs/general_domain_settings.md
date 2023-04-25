## General Domain Settings Access
When accessing General Domain Settings, which includes:
- settings
- redirects
- cache settings
- IP blacklist/whitelist
- IP rate limit
- WAF rules

there is specific format that needs to be submitted to the API.

Instead of specifying a subdomain name, ```"ALL-1234"``` needs to be passed, where 1234 is domain ID.

### Example - listing redirects for General Domain:
```go
redirects, err := api.ListRedirects(1234, "ALL-1234", nil)
if err != nil {
    log.Fatal(err)
}
```