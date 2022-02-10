# Setup myrasec-go

To be able to access the Myra API using the myrasec-go lib, you need a API Key and API Secret. 
You can manage your API keys and other profile settings in the Myra user management on your profile page. Please make sure to note your API Secret as it is visible only once in the moment when you create new API Key.

## Setup example
```go
api, err := myrasec.New(os.Getenv("MYRA_API_KEY"), os.Getenv("MYRA_API_SECRET"))
if err != nil {
    log.Fatal(err)
}
```

## List 100 domains example
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
	
	domains, err := api.ListDomains(map[string]string{"pageSize": "100"})
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range domains {
		log.Println(d.Name)
	}
}
```
