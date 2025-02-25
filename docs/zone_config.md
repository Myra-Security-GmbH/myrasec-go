# Bind Zone Config

## Bind Zone Config string
To get the bind zone config as string as it's required in the bind config file
```go
zoneConfig, err := api.GetZoneConfigRaw(domainId, map[string]string)
if err != nil {
    log.Fatal(err)
}
log.Println(zoneConfig)
```
Response will be something like this
```
$ttl 38400
example.com. IN SOA a.myradns.net. support.myrasecurity.com. (
        2740493835
        16384
        2048
        1048576
        2560 )
example.com. IN NS a.myradns.net
example.com. IN NS a.myradns.net
www.example.com. 300 IN A 192.168.1.2
www.example.com. 300 IN AAAA FE80::1
```

## Bind Zone Config as json
To get the data printed to the bind zone config as json you can use this method
```go
zoneConfig, err := api.GetZoneConfigJson(domainId, map[string]string)
if err != nil {
    log.Fatal(err)
}
log.Println(zoneConfig)
```

Response will be something like this
```json
{
    "error": false,
    "violationList": [],
    "warningList": [],
    "data": [
        {
            "objectType": "DomainBindVO",
            "domainName": "example.com",
            "ttl": 38400,
            "primaryDns": "a.myradns.net",
            "hostmaster": "support.myrasecurity.com",
            "timeToRefresh": 16384,
            "timeToRetry": 2048,
            "timeToExpire": 1048576,
            "minimumTTL": 2560,
            "version": 2740493713,
            "records": [
                {
                    "domain": "example.com",
                    "type": "NS",
                    "value": "a.myradns.net",
                    "rtype": "dns"
                },
                {
                    "domain": "example.com",
                    "type": "NS",
                    "value": "b.myradns.net",
                    "rtype": "dns"
                },
                {
                    "domain": "www.example.com",
                    "ttl": 300,
                    "type": "A",
                    "value": "192.168.1.2",
                    "cnameAlt": "www-example-com.ax4z.com.",
                    "rtype": "record",
                    "active": true
                },
                {
                    "domain": "www.example.com",
                    "ttl": 300,
                    "type": "AAAA",
                    "value": "FE80::1",
                    "cnameAlt": "www-example-com.ax4z.com.",
                    "rtype": "record",
                    "active": true
                }
            ]
        }
    ]
}
```
