# Subdomain settings

```go
type Settings struct {
	AccessLog                   bool     `json:"access_log"`
	AntibotPostFlood            bool     `json:"antibot_post_flood"`
	AntibotPostFloodThreshold   int      `json:"antibot_post_flood_threshold,omitempty"`
	AntibotProofOfWork          bool     `json:"antibot_proof_of_work"`
	AntibotProofOfWorkThreshold int      `json:"antibot_proof_of_work_threshold,omitempty"`
	BalancingMethod             string   `json:"balancing_method,omitempty"`
	BlockNotWhitelisted         bool     `json:"block_not_whitelisted"`
	BlockTorNetwork             bool     `json:"block_tor_network"`
	CacheEnabled                bool     `json:"cache_enabled"`
	CacheRevalidate             bool     `json:"cache_revalidate"`
	CDN                         bool     `json:"cdn"`
	ClientMaxBodySize           int      `json:"client_max_body_size,omitempty"`
	DiffieHellmanExchange       int      `json:"diffie_hellman_exchange,omitempty"`
	EnableOriginSNI             bool     `json:"enable_origin_sni"`
	ForwardedForReplacement     string   `json:"forwarded_for_replacement,omitempty"`
	HSTS                        bool     `json:"hsts"`
	HSTSIncludeSubdomains       bool     `json:"hsts_include_subdomains"`
	HSTSMaxAge                  int      `json:"hsts_max_age,omitempty"`
	HSTSPreload                 bool     `json:"hsts_preload"`
	HTTPOriginPort              int      `json:"http_origin_port,omitempty"`
	IgnoreNoCache               bool     `json:"ignore_nocache"`
	ImageOptimization           bool     `json:"image_optimization"`
	IPv6Active                  bool     `json:"ipv6_active"`
	LimitAllowedHTTPMethod      []string `json:"limit_allowed_http_method,omitempty"`
	LimitTLSVersion             []string `json:"limit_tls_version,omitempty"`
	LogFormat                   string   `json:"log_format,omitempty"`
	MonitoringAlertThreshold    int      `json:"monitoring_alert_threshold,omitempty"`
	MonitoringContactEMail      string   `json:"monitoring_contact_email,omitempty"`
	MonitoringSendAlert         bool     `json:"monitoring_send_alert"`
	MyraSSLHeader               bool     `json:"myra_ssl_header"`
	MyraSSLCertificate          []string `json:"myra_ssl_certificate"`
	MyraSSLCertificateKey       []string `json:"myra_ssl_certificate_key"`
	NextUpstream                []string `json:"next_upstream,omitempty"`
	OnlyHTTPS                   bool     `json:"only_https"`
	OriginConnectionHeader      string   `json:"origin_connection_header,omitempty"`
	ProxyCacheBypass            string   `json:"proxy_cache_bypass,omitempty"`
	ProxyCacheStale             []string `json:"proxy_cache_stale,omitempty"`
	ProxyConnectTimeout         int      `json:"proxy_connect_timeout,omitempty"`
	ProxyReadTimeout            int      `json:"proxy_read_timeout,omitempty"`
	RequestLimitBlock           string   `json:"request_limit_block,omitempty"`
	RequestLimitLevel           int      `json:"request_limit_level,omitempty"`
	RequestLimitReport          bool     `json:"request_limit_report"`
	RequestLimitReportEMail     string   `json:"request_limit_report_email,omitempty"`
	Rewrite                     bool     `json:"rewrite"`
	SourceProtocol              string   `json:"source_protocol,omitempty"`
	Spdy                        bool     `json:"spdy"`
	SSLOriginPort               int      `json:"ssl_origin_port,omitempty"`
	WAFEnable                   bool     `json:"waf_enable"`
	WAFLevelsEnable             []string `json:"waf_levels_enable,omitempty"`
	WAFPolicy                   string   `json:"waf_policy,omitempty"`
}
```

| Field | Type | Description|
|---|---|---|
| AccessLog | bool | Activate separated access log. A access log from each Myra node delivering your website will be saved. You can download the access log files via sftp from custlogs.myracloud.com. |
| AntibotPostFlood | bool | Detection of POST floods by using a JavaScript based puzzle. |
| AntibotPostFloodThreshold | int | This parameter determines the frequency how often the puzzle has to be solved. The higher the value the less likely the puzzle needs to be solved. |
| AntibotProofOfWork | bool | Detection of valid clients by using a JavaScript based puzzle. |
| AntibotProofOfWorkThreshold | int | This parameter determines the frequency how often the puzzle has to be solved. The higher the value the less likely the puzzle needs to be solved. |
| BalancingMethod | string | Specifies with which method requests are balanced between upstream servers. The default behavior is the round-robin balancing. The value ip_hash will cause Myra to forward the same client IP always to the same upstream server. The value least_conn will cause Myra to forward the request to the upstream server with least connections |
| BlockNotWhitelisted | bool | Block all IPs which are not whitelisted in the IP filter settings |
| BlockTorNetwork | bool | Block traffic from the TOR network. |
| CacheEnabled | bool | Turn caching on or off. If you enable the cache, you also have to define the objects to be cached in the cache settings. |
| CacheRevalidate | bool | If enabled, expired cache items will be requested with the additional HTTP header "If-Modified-Since" and "If-None-Match" |
| CDN | bool | Should this subdomain be used as Content Delivery Node (CDN). After enabling the CDN you will be able to create buckets and upload files using the Myra upload API. |
| ClientMaxBodySize | int | Sets the maximum allowed size of the client request body, specified in the “Content-Length” request header field. |
| DiffieHellmanExchange | int | Defines the size of the Diffie-Hellman key exchange parameters in bits. Please, note that Java 6 and 7 do not support Diffie-Hellman parameters larger than 1024 bits. If your server expects to receive connections from java 6 clients and wants to enable PFS, it must provide a DHE parameter of 1024 bits |
| EnableOriginSNI | bool | Enable or disable origin SNI. |
| ForwardedForReplacement | string | Set your own X-Forwarded-For header. |
| HSTS | bool | Enable HSTS protection for a domain. This will tell browsers to use secure https connections only when interacting with your domain. |
| HSTSIncludeSubdomains | bool | This will extend the HSTS protection for all subdomains |
| HSTSMaxAge | int | Specified how long the HSTS header is valid before the browser has to revalidate. |
| HSTSPreload | bool | Allow the domain to be added to the HSTS preload list used by all major browsers (https://hstspreload.appspot.com/). |
| HTTPOriginPort | int | Allows to set a port for communication with origin via HTTP. |
| IgnoreNoCache | bool | If activated, no-cache headers (Cache-Control: [private|no-store|no-cache]) will be ignored. |
| ImageOptimization | bool | Activate lossless optimization of JPEG and PNG images (recommended setting). |
| IPv6Active | bool | Allow connections via IPv6 to your systems. IPv4 connections will be forwarded in any case. |
| LimitAllowedHTTPMethod | []string | Not selected HTTP methods will be blocked. |
| LimitTLSVersion | []string | Only selected TLS versions will be used. |
| LogFormat | string | Use a different log format. |
| MonitoringAlertThreshold | int | Errors per minute that must occur until a report is sent. |
| MonitoringContactEMail | string | Email addresses, to which monitoring emails should be send. Multiple addresses are separated with a space. |
| MonitoringSendAlert | bool | Enables / disables the upstream error reporting. |
| MyraSSLHeader | bool | Activate the X-Myra-SSL Header, which indicates if a request was received via SSL. |
| MyraSSLCertificate | []string | An SSL Certificate (and chain) to be used to make requests on the origin. |
| MyraSSLCertificateKey | []string | The private key for the SSL Certificate |
| NextUpstream | []string | Specify in which case the current upstream should be marked as "down". The values can be arbitrary combined, expect the value "off". |
| OnlyHTTPS | bool | If activated, Myra will forward all requests to the origin using HTTPS regardless of the used protocol of the originating request. |
| OriginConnectionHeader | string | Sets the Connection header, which is transmitted to the origin with a request. |
| ProxyCacheBypass | string | Defines the name of the cookie which forces Myra to deliver the response not from cache. The values of the cookie must be not empty or equal to 0 to enable bypassing. |
| ProxyCacheStale | []string | Determines in which cases a stale cached response can be used when an error occurs during communication with your server. The values can be arbitrary combined, expect the value "off". |
| ProxyConnectTimeout | int | Defines a timeout in seconds for establishing a connection with the origin server. The timeout cannot be greater than 60 seconds. |
| ProxyReadTimeout | int | Defines a timeout in seconds for reading a response from the proxied server. The timeout is set only between two successive read operations, not for the transmission of the whole response. |
| RequestLimitBlock | string | If activated, the user has to solve a CAPTCHA after exceeding the configured request limit. |
| RequestLimitLevel | int | Define how many requests are allowed from an IP per minute. If this limit is reached, the IP will be blocked. If request_limit_block is enabled, the user can solve a CAPTCHA to unblock his IP address. |
| RequestLimitReport | bool | If activated, an email will be send containing blocked ip addresses that exceeded the configured request limit. |
| RequestLimitReportEMail | string |  |
| Rewrite | bool | Enable automated JavaScript optimization. All JavaScript is collected and executed at the end of the page. This significantly decreases the DOM content loaded time. If not all JavaScript files should be collected you can set the value to "regex" and specify the regex to use while matching filenames in the option "rewrite_regex". |
| SourceProtocol | string | Define which protocol should be used when passing a request to your servers. The value "same" will ensure that the same protocol is used as in the originating request to Myra. The "http" and "https" value will force Myra to always use the specified protocol when connecting. |
| Spdy | bool | Activate the high performance HTTP/2 protocol. Please note that you have to enable HTTPS for Myra to get HTTP/2 enabled. |
| SSLOriginPort | int | Allows to set a port for communication with origin via SSL. |
| WAFEnable | bool | Enables / disables the Web Application Firewall. |
| WAFLevelsEnable | []string | Level of applied WAF rules. |
| WAFPolicy | string | Default policy for the Web Application Firewall in case of rule error. |


## Read
The listing operation returns a settings object for the given domainId and subdomain name.

### Example
```go
settings, err := api.ListSettings(domainId, "www.example.com", nil)
if err != nil {
    log.Fatal(err)
}
```

The listing full operation returns an interface of the settings. To get a map (key, value) of the domain settings you have to convert the structure as shown here:
```go
settingsData, err := api.ListSettingsFull(domainId, "www.example.com", nil)
if err != nil {
	log.Fatal(err)
}
allSettings, _ := settingsData.(*map[string]interface{})
domainSettings := (*allSettings)["domain"]
settingsMap, ok := domainSettings.(map[string]interface{})
```

**Note:** To have a consistent API, the ListSettings function allows to pass a params map. But in fact, no params is used/interpreted on this API request.

## Create/Update/Delete
To create/update/delete (sub)domain settings you have to send a `map[string]interface{}` to update or delete attributes.

To create/update you add the specific attribute to the map with the required value.

To delete an attribute you have to add the attribute with `nil`.

Only attributes in the map will be touched in the api.
### Example
```go
settingsMap := make(map[string]interface{})
settingsMap['only_https'] = true // update/create
settingsMap['myra_ssl_header'] = nil // delete
s, err = api.UpdateSettingsPartial(settingsMap, domainId, "www.example.com")
if err != nil {
    log.Fatal(err)
}
```