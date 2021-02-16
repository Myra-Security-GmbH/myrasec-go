# myrasec-go

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
