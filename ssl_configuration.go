package myrasec

import (
	"fmt"
	"net/http"
)

func getSslConfigurationMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listSSLConfigurations": {
			Name:   "listSSLConfigurations",
			Action: "ssl-configurations",
			Method: http.MethodGet,
			Result: []SslConfiguration{},
		},
	}
}

type SslConfiguration struct {
	ID        int    `json:"id" jsonschema:"Identifier of that configuration."`
	Name      string `json:"name" jsonschema:"Unique string identifier of this configuration."`
	Ciphers   string `json:"ciphers" jsonschema:"List of ciphers for this configuration."`
	Protocols string `json:"protocols" jsonschema:"List of protocols for this configuration."`
}

func (api *API) ListSslConfigurations() ([]SslConfiguration, error) {
	if _, ok := methods["listSSLConfigurations"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSSLConfigurations")
	}

	definition := methods["listSSLConfigurations"]

	result, err := api.call(definition)
	if err != nil {
		return nil, err
	}

	return *result.(*[]SslConfiguration), nil
}
