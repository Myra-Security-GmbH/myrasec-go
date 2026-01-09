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

// SslConfiguration represents a custom TLS/SSL profile.
// It defines the specific cipher suites and protocols allowed for secure connections.
type SslConfiguration struct {
	// ID is the unique identifier for the SSL configuration.
	// This value is server-generated and read-only.
	ID int `json:"id" jsonschema:"The unique numeric identifier for this configuration. Server-generated and read-only. Ignored during creation."`

	// Name is a unique label for this configuration profile.
	Name string `json:"name" jsonschema:"A unique name or label to identify this SSL configuration (e.g., 'Myra-Global-TLS-Default')."`

	// Ciphers defines the list of allowed cipher suites.
	// Expects an OpenSSL cipher string format.
	Ciphers string `json:"ciphers" jsonschema:"The OpenSSL cipher suite specification string. Defines which encryption algorithms are allowed (e.g., 'ECDHE-RSA-AES256-GCM-SHA384:...')."`

	// Protocols defines the list of allowed TLS protocols.
	// Typically a space-separated string (e.g., "TLSv1.2 TLSv1.3").
	Protocols string `json:"protocols" jsonschema:"Space-separated list of enabled TLS protocols (e.g., 'TLSv1.2 TLSv1.3')."`
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
