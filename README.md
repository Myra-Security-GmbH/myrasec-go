# myrasec-go

[![Go Report Card](https://goreportcard.com/badge/github.com/Myra-Security-GmbH/myrasec-go/v2)](https://goreportcard.com/report/github.com/Myra-Security-GmbH/myrasec-go/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/Myra-Security-GmbH/myrasec-go/v2.svg)](https://pkg.go.dev/github.com/Myra-Security-GmbH/myrasec-go/v2)
[![tests](https://github.com/Myra-Security-GmbH/myrasec-go/actions/workflows/test.yml/badge.svg)](https://github.com/Myra-Security-GmbH/myrasec-go/actions/workflows/test.yml)

A Go library for interacting with Myra Security API.

## Example usage

```go
package main

import (
	"log"
	"os"

	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
)

func main() {
	api, err := myrasec.New(os.Getenv("MYRA_API_KEY"), os.Getenv("MYRA_API_SECRET"))
	// Alternatively you can use the API token for authentication, too
	// api, err := myrasec.NewWithToken(os.Getenv("MYRA_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	domains, err := api.ListDomains(map[string]string{"pageSize": "100"})
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range domains {
		log.Println(d.Name)
	}
}
```

## Documentation
- [Setup](./docs/setup.md)
- [Domain](./docs/domain.md)
- [DNS record](./docs/dns_record.md)
- [Cache setting](./docs/cache_setting.md)
- [Subdomain settings](./docs/subdomain_settings.md)
- [IP filter](./docs/ip_filter.md)
- [Redirect](./docs/redirect.md)
- [IP ranges](./docs/ip_range.md)
- [SSL certificates](./docs/ssl.md)
- [Maintenance](./docs/maintenance.md)
- [Maintenance templates](./docs/maintenance_template.md)
- [Error page](./docs/error_page.md)
- [Tag](./docs/tag.md)
- [Tag cache setting](./docs/tag_cachesetting.md)
- [Tag settings](./docs/tag_settings.md)
- [Tag WAF rule](./docs/tag_wafrule.md)
- [CDN bucket](./docs/cdn_bucket.md)
- [CDN file](./docs/cdn_file.md)
- [Statistics](./docs/statistics.md)
- [CacheClear](./docs/cacheclear.md)
- [General Domain settings](./docs/general_domain_settings.md)
- [WAF](./docs/waf.md)
- [WAF Actions](./docs/waf_action.md)
- [WAF Conditions](./docs/waf_condition.md)
- [Waitingroom](./docs/waitingroom.md)
- [Bind Zone Config](./docs/zone_config.md)
