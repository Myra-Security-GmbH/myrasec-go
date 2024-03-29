package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getDomainMethods returns Domain related API calls
func getDomainMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getDomain": {
			Name:               "getDomain",
			Action:             "domains/%d",
			Method:             http.MethodGet,
			Result:             Domain{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listDomains": {
			Name:   "listDomains",
			Action: "domains",
			Method: http.MethodGet,
			Result: []Domain{},
		},
		"createDomain": {
			Name:   "createDomain",
			Action: "domains",
			Method: http.MethodPost,
			Result: Domain{},
		},
		"updateDomain": {
			Name:   "updateDomain",
			Action: "domains/%d",
			Method: http.MethodPut,
			Result: Domain{},
		},
		"deleteDomain": {
			Name:   "deleteDomain",
			Action: "domains/%d",
			Method: http.MethodDelete,
			Result: Domain{},
		},
	}
}

// Domain ...
type Domain struct {
	ID          int             `json:"id,omitempty"`
	Created     *types.DateTime `json:"created,omitempty"`
	Modified    *types.DateTime `json:"modified,omitempty"`
	Name        string          `json:"name"`
	AutoUpdate  bool            `json:"autoUpdate"`
	AutoDNS     bool            `json:"autoDns"`
	Paused      bool            `json:"paused"`
	PausedUntil *types.DateTime `json:"pausedUntil,omitempty"`
	Reversed    bool            `json:"reversed"`
}

// GetDomain returns a single domain with/for the given identifier
func (api *API) GetDomain(id int) (*Domain, error) {
	if _, ok := methods["getDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getDomain")
	}

	definition := methods["getDomain"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*Domain), nil
}

// ListDomains returns a slice containing all visible domains
func (api *API) ListDomains(params map[string]string) ([]Domain, error) {
	if _, ok := methods["listDomains"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listDomains")
	}

	definition := methods["listDomains"]

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var domains []Domain
	domains = append(domains, *result.(*[]Domain)...)

	return domains, nil
}

// CreateDomain creates a new domain using the MYRA API
func (api *API) CreateDomain(domain *Domain) (*Domain, error) {
	if _, ok := methods["createDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createDomain")
	}

	definition := methods["createDomain"]

	result, err := api.call(definition, domain)
	if err != nil {
		return nil, err
	}
	return result.(*Domain), nil
}

// UpdateDomain updates the passed domain using the MYRA API
func (api *API) UpdateDomain(domain *Domain) (*Domain, error) {
	if _, ok := methods["updateDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateDomain")
	}

	definition := methods["updateDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domain.ID)

	result, err := api.call(definition, domain)
	if err != nil {
		return nil, err
	}
	return result.(*Domain), nil
}

// DeleteDomain deletes the passed domain using the MYRA API
func (api *API) DeleteDomain(domain *Domain) (*Domain, error) {
	if _, ok := methods["deleteDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteDomain")
	}

	definition := methods["deleteDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domain.ID)

	_, err := api.call(definition, domain)
	if err != nil {
		return nil, err
	}
	return domain, nil
}
