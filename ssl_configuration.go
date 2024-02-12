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
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Ciphers   string `json:"ciphers"`
	Protocols string `json:"protocols"`
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

	var records []SslConfiguration
	records = append(records, *result.(*[]SslConfiguration)...)

	return records, nil
}
