# myrasec-go

[![Go Report Card](https://goreportcard.com/badge/github.com/Myra-Security-GmbH/myrasec-go)](https://goreportcard.com/report/github.com/Myra-Security-GmbH/myrasec-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/Myra-Security-GmbH/myrasec-go.svg)](https://pkg.go.dev/github.com/Myra-Security-GmbH/myrasec-go)

A Go library for interacting with Myra Security API.

> **Note**: This library is under active development. 
> Upcoming changes may break existing functionality.
> Consider this library as unstable.

## Features
### Domains
* [X] List existing Domains
* [X] Create a new Domain
* [X] Update a Domain
* [X] Delete a Domain
### DNS Records
* [X] List existing DNS Records
* [X] Create a new DNS Record
* [X] Update a DNS Record
* [X] Delete a DNS Record
### Redirects
* [X] List existing Redirects
* [X] Create a new Redirect
* [X] Update a Redirect
* [X] Delete a Redirect
### CacheSettings
* [X] List existing Cache Settings
* [X] Create a new Cache Setting
* [X] Update a Cache Setting
* [X] Delete a Cache Setting

### Settings
* [X] List Settings for a subdomain
* [X] Update Settings for a subdomain


### IP Filter
* [X] List existing IP Filters
* [X] Create a new IP Filter
* [X] Update an IP Filter
* [X] Delete an IP Filter

### Rate Limit
* [X] List existing Rate Limit Settings
* [X] Create a new Rate Limit Setting
* [X] Update a Rate Limit Setting
* [X] Delete a Rate Limit Setting

### WAF
* [X] List existing WAF rules
* [X] Fetch a singe WAF rule
* [X] Fetch existing WAF conditions
* [X] Fetch existing WAF actions
* [X] Create a new WAF rule
* [X] Update a WAF rule
* [X] Delete a WAF rule


## Usage
```go
package main

import (
	"log"
    "os"

	myrasec "github.com/Myra-Security-GmbH/myrasec-go"
)
func main() {
	api, err := myrasec.New(os.Getenv("MYRA_API_KEY"), os.Getenv("MYRA_API_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

    domains, err := api.ListDomains()
	if err != nil {
		return err
	}

	for _, d := range domains {
		fmt.Println(d.Name)
    }
}
```
