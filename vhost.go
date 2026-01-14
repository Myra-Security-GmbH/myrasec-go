package myrasec

import (
	"fmt"
	"net/http"
)

// getVHostMethods returns VHost related API calls
func getVHostMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listAllSubdomains": {
			Name:   "listAllSubdomains",
			Action: "subdomains",
			Method: http.MethodGet,
			Result: []VHost{},
		},
		"listSubdomainsForDomain": {
			Name:   "listSubdomainsForDomain",
			Action: "domain/%d/subdomains",
			Method: http.MethodGet,
			Result: []VHost{},
		},
	}
}

// VHost represents a Virtual Host configuration within a domain.
// It maps specific hostnames (subdomains) to specific settings or backends.
type VHost struct {
	// ID is the unique identifier for the VHost.
	// Server-generated.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the Virtual Host. Server-generated; required for updates and deletes, but ignored during creation."`

	// Label is a descriptive name for the VHost.
	// Used for easier identification in the UI.
	Label string `json:"label,omitempty" jsonschema:"The specific hostname or FQDN associated with this VHost without trailing dot (e.g., 'shop.example.com')."`

	// Value specifies the actual hostname.
	// This is the FQDN handled by this VHost.
	Value string `json:"value,omitempty" jsonschema:"The specific hostname or FQDN associated with this VHost (e.g., 'shop.example.com.')."`

	// DomainName is the parent domain of the VHost.
	DomainName string `json:"domainName,omitempty" jsonschema:"The Fully Qualified Domain Name (FQDN) of the parent domain."`

	// Access indicates if the VHost is generally accessible.
	Access bool `json:"access" jsonschema:"Indicates if the VHost is configured to be accessible (active)."`

	// Paused indicates if the VHost is temporarily suspended.
	// If true, traffic might be blocked or not processed.
	Paused bool `json:"paused" jsonschema:"Indicates if the VHost is currently paused. If true, traffic processing is suspended."`
}

// ListAllSubdomains ...
func (api *API) ListAllSubdomains(params map[string]string) ([]VHost, error) {
	if _, ok := methods["listAllSubdomains"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listAllSubdomains")
	}

	definition := methods["listAllSubdomains"]

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]VHost), nil
}

// ListAllSubdomainsForDomain ...
func (api *API) ListAllSubdomainsForDomain(domainId int, params map[string]string) ([]VHost, error) {
	if _, ok := methods["listSubdomainsForDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSubdomainsForDomain")
	}

	definition := methods["listSubdomainsForDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var vhosts []VHost
	vhosts = append(vhosts, *result.(*[]VHost)...)

	return vhosts, nil
}
