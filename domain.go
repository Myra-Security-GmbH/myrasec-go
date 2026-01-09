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
	ID          int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a Domain it is necessary to add this attribute to your object."`
	Created     *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new Domain object. This value is informational so it is not necessary to add this attribute to any API call."`
	Modified    *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format."`
	Name        string          `json:"name" jsonschema:"Identifies the domain by its name. The value cannot be changed after creation. To change a typo you need to remove and recreate the domain."`
	AutoUpdate  bool            `json:"autoUpdate" jsonschema:"Shows if the current domain has autoUpdate activated. If autoUpdate is deactivated changes on your configuration are not deployed until you reactivate autoUpdate. This is primary used to change a lot of settings at once to prevent Myra to deploy a half done configuration. In some cases Myra support also deactivates this option to prevent Myra system from removing special configuration settings. Please note that turning autoUpdate off is not correlated to database transactions. This means that any changes are saved but not deployed."`
	AutoDNS     bool            `json:"autoDns" jsonschema:"If AutoDNS flag is set while creating a new domain Myra tries to get a list of subDomains for this domain. Depending on your DNS provider configuration this may fail or return a incomplete list. For best results Myra recomments to use the subDomain API to create DNS records."`
	Paused      bool            `json:"paused" jsonschema:"Shows if the domain is currently in pause mode."`
	PausedUntil *types.DateTime `json:"pausedUntil,omitempty" jsonschema:"Shows the date when Myra protection will be reactivated automatically."`
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

	return *result.(*[]Domain), nil
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
