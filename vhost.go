package myrasec

import "fmt"

//
// VHost ...
//
type VHost struct {
	ID         int    `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Value      string `json:"value,omitempty"`
	DomainName string `json:"domainName,omitempty"`
	Access     bool   `json:"access,omitempty"`
	Paused     bool   `json:"paused,omitempty"`
}

//
// ListAllSubdomains ...
//
func (api *API) ListAllSubdomains(params map[string]string) ([]VHost, error) {
	if _, ok := methods["listAllSubdomains"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listAllSubdomains")
	}

	definition := methods["listAllSubdomains"]
	definition.Action = fmt.Sprintf(definition.Action)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var vhosts []VHost
	for _, v := range *result.(*[]VHost) {
		vhosts = append(vhosts, v)
	}

	return vhosts, nil
}

//
// ListAllSubdomainsForDomain ...
//
func (api *API) ListAllSubdomainsForDomain(domainId int, params map[string]string) ([]VHost, error) {
	if _, ok := methods["listSubdomainsForDomain"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listSubdomainsForDomain")
	}

	definition := methods["listSubdomainsForDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var vhosts []VHost
	for _, v := range *result.(*[]VHost) {
		vhosts = append(vhosts, v)
	}

	return vhosts, nil
}
