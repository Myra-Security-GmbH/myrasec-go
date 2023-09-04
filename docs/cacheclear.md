# Cache Clear

To enqueue a clear operation you need to send a CacheClear object. The `FQDN` attribute in the object must be a subdomain of the domain mentioned in the call.

```go
type CacheClear struct {
    FQDN      string `json:"fqdn"`
    Resource  string `json:"resource"`
    Recursive bool   `json:"recursive"`
}
```
| Field | Type | Description|
|---|---|---|
| `FQDN` | string | A specific SubdomainName. |
| `Resource` | string | A specific path to be cache-cleared. Wildcards can be used. (The path is relative to the root of the website).|
| `Recursive` | bool |  A boolean value to define if the cache purge should be done recursively. |



## Clear the cache of a subdomain
To enqueue a cache clear operation for all `*.jpg` resources in `www.example.com` you can do the following call:

### Example
```go
domainId := 123
cc := &myrasec.ClearCache{
    FQDN: "www.example.com",
    Resource: "/*.jpg",
    Recursive: true
}
cc, err := api.ClearCache(cc, domainId)
if err != nil {
    panic(err)
}

log.Println(cc)
```


## Resource Pattern matching
Internally we use the [fnmatch](https://man7.org/linux/man-pages/man3/fnmatch.3.html) (flags=FNM_PATHNAME) function to find the matching resources that should be deleted. To allow you to do recurseve deletion in folder, the flag `recursive` was added.

| Pattern | Recursive | Resource | Result |
|---|---|---|---|
| `/*.js` | No | `/main.js` | **match** |
| `/*.js` | No | `/folder/main.js` | **no match** |
| `/*.js` | No | `/testmain.js` | **no match** |
| `.js` | Yes | `/assets/script.js` | **match** |
| `.js` | Yes | `/assets/jquery/jquery.js` | **match** |
| `.js` | Yes | `/main.js` | **match** |
| `.js` | Yes | `/main.css` | **no match** |
| `.js` | Yes | `/assets/js/source.map` | **no match** |
| `/assets/*.js` | No | `/assets/script.js` | **match** |
| `/assets/*.js` | No | `/asset/script.js` | **no match** |
| `/assets/*.js` | No | `/main.js` | **no match** |
| `/assets/*.js` | No | `/folder/js/script.js` | **no match** |
| `/assets/*.js` | Yes | `/assets/script.js` | **match** |
| `/assets/*.js` | Yes | `/assets/jquery/jquery.js` | **match** |
| `/assets/*.js` | Yes | `/main.js` | **no match** |
| `/assets/*.js` | Yes | `/js/angular.js` | **no match** |
| `/*.*` | Yes | `/main.js` | **match** |
| `/*` | Yes | `/main.js` | **match** |
