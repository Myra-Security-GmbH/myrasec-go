package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/pkg/types"
)

//
// Domain ...
//
type Domain struct {
	ID          int             `json:"id,omitempty"`
	Created     *types.DateTime `json:"created,omitempty"`
	Modified    *types.DateTime `json:"modified,omitempty"`
	Name        string          `json:"name"`
	AutoUpdate  bool            `json:"autoUpdate"`
	AutoDNS     bool            `json:"autoDns"`
	Paused      bool            `json:"paused"`
	PausedUntil *types.DateTime `json:"pausedUntil,omitempty"`
}

//
// ListDomains returns a slice containing all visible domains
//
func (api *API) ListDomains(params map[string]string) ([]Domain, error) {
	if _, ok := methods["listDomains"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listDomains")
	}

	definition := methods["listDomains"]
	definition.Action = fmt.Sprintf(definition.Action, 1)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var domains []Domain
	for _, v := range *result.(*[]Domain) {
		domains = append(domains, v)
	}

	return domains, nil
}

//
// CreateDomain creates a new domain using the MYRA API
//
func (api *API) CreateDomain(domain *Domain) (*Domain, error) {
	if _, ok := methods["createDomain"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "createDomain")
	}

	result, err := api.call(methods["createDomain"], domain)
	if err != nil {
		return nil, err
	}
	return result.(*Domain), nil
}

//
// UpdateDomain updates the passed domain using the MYRA API
//
func (api *API) UpdateDomain(domain *Domain) (*Domain, error) {
	if _, ok := methods["updateDomain"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "updateDomain")
	}

	result, err := api.call(methods["updateDomain"], domain)
	if err != nil {
		return nil, err
	}
	return result.(*Domain), nil
}

//
// DeleteDomain deletes the passed domain using the MYRA API
//
func (api *API) DeleteDomain(domain *Domain) (*Domain, error) {
	if _, ok := methods["deleteDomain"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteDomain")
	}

	result, err := api.call(methods["deleteDomain"], domain)
	if err != nil {
		return nil, err
	}
	return result.(*Domain), nil
}
