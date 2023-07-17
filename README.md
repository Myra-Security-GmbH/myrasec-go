# myrasec-go

[![Go Report Card](https://goreportcard.com/badge/github.com/Myra-Security-GmbH/myrasec-go)](https://goreportcard.com/report/github.com/Myra-Security-GmbH/myrasec-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/Myra-Security-GmbH/myrasec-go.svg)](https://pkg.go.dev/github.com/Myra-Security-GmbH/myrasec-go)
[![tests](https://github.com/Myra-Security-GmbH/myrasec-go/actions/workflows/test.yml/badge.svg)](https://github.com/Myra-Security-GmbH/myrasec-go/actions/workflows/test.yml)

A Go library for interacting with Myra Security API.

> **Note**: This library is under active development. 
> Upcoming changes may break existing functionality.
> Consider this library as unstable.

## Example usage

```go
package main

import (
	"log"
	"os"

	myrasec "github.com/Myra-Security-GmbH/myrasec-go/"
)

func main() {
	api, err := myrasec.New(os.Getenv("MYRA_API_KEY"), os.Getenv("MYRA_API_SECRET"))
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
- [IP rate limit](./docs/ratelimit.md)
- [Redirect](./docs/redirect.md)
- [IP ranges](./docs/ip_range.md)
- [SSL certificates](./docs/ssl.md)
- [Maintenance](./docs/maintenance.md)
- [Maintenance templates](./docs/maintenance_template.md)
- [Error page](./docs/error_page.md)
- [Tag](./docs/tag.md)
- [Tag cache setting](./docs/tag_cachesetting.md)
- [Tag rate limit](./docs/tag_ratelimit.md)
- [Tag settings](./docs/tag_settings.md)
- [Tag WAF rule](./docs/tag_wafrule.md)
- [CDN bucket](./docs/cdn_bucket.md)
- [CDN file](./docs/cdn_file.md)
- [Statistics](./docs/statistics.md)
- [General Domain settings](./docs/general_domain_settings.md)
