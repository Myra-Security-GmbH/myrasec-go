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

// VHost ...
type VHost struct {
	ID         int    `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Value      string `json:"value,omitempty"`
	DomainName string `json:"domainName,omitempty"`
	Access     bool   `json:"access"`
	Paused     bool   `json:"paused"`
}

// ListAllSubdomains ...
func (api *API) ListAllSubdomains(params map[string]string) ([]VHost, error) {
	if _, ok := methods["listAllSubdomains"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listAllSubdomains")
	}

	definition := methods["listAllSubdomains"]
	definition.Action = fmt.Sprintf(definition.Action)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var vhosts []VHost
	vhosts = append(vhosts, *result.(*[]VHost)...)

	return vhosts, nil
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
