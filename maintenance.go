package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getMaintenanceMethods returns Maintenance related API calls
func getMaintenanceMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listMaintenances": {
			Name:   "listMaintenances",
			Action: "domain/%d/%s/maintenances",
			Method: http.MethodGet,
			Result: []Maintenance{},
		},
		"createMaintenance": {
			Name:   "createMaintenance",
			Action: "domain/%d/%s/maintenances",
			Method: http.MethodPost,
			Result: Maintenance{},
		},
		"updateMaintenance": {
			Name:   "updateMaintenance",
			Action: "domain/%d/%s/maintenances/%d",
			Method: http.MethodPut,
			Result: Maintenance{},
		},
		"deleteMaintenance": {
			Name:   "deleteMaintenance",
			Action: "domain/%d/%s/maintenances/%d",
			Method: http.MethodDelete,
			Result: Maintenance{},
		},
	}
}

// Maintenance ...
type Maintenance struct {
	ID          int             `json:"id,omitempty"`
	Created     *types.DateTime `json:"created,omitempty"`
	Modified    *types.DateTime `json:"modified,omitempty"`
	Start       *types.DateTime `json:"start,omitempty"`
	End         *types.DateTime `json:"end,omitempty"`
	Active      bool            `json:"active"`
	Content     string          `json:"content"`
	ContentFrom string          `json:"contentFrom,omitempty"`
	FQDN        string          `json:"fqdn"`
}

// ListMaintenances returns a slice containing all maintenance pages for a subdomain
func (api *API) ListMaintenances(domainId int, subDomainName string, params map[string]string) ([]Maintenance, error) {
	if _, ok := methods["listMaintenances"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listMaintenances")
	}

	definition := methods["listMaintenances"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []Maintenance
	records = append(records, *result.(*[]Maintenance)...)

	return records, nil
}

// CreateMaintenance creates a new maintenance page for the passed subdomain (name) using the MYRA API
func (api *API) CreateMaintenance(maintenance *Maintenance, domainId int, subDomainName string) (*Maintenance, error) {
	if _, ok := methods["createMaintenance"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createMaintenance")
	}

	definition := methods["createMaintenance"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, maintenance)
	if err != nil {
		return nil, err
	}
	return result.(*Maintenance), nil
}

// UpdateMaintenance updates the passed maintenance page using the MYRA API
func (api *API) UpdateMaintenance(maintenance *Maintenance, domainId int, subDomainName string) (*Maintenance, error) {
	if _, ok := methods["updateMaintenance"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateMaintenance")
	}

	definition := methods["updateMaintenance"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, maintenance.ID)

	result, err := api.call(definition, maintenance)
	if err != nil {
		return nil, err
	}
	return result.(*Maintenance), nil
}

// DeleteMaintenance deletes the passed maintenance page using the MYRA API
func (api *API) DeleteMaintenance(maintenance *Maintenance, domainId int, subDomainName string) (*Maintenance, error) {
	if _, ok := methods["deleteMaintenance"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteMaintenance")
	}

	definition := methods["deleteMaintenance"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, maintenance.ID)

	_, err := api.call(definition, maintenance)
	if err != nil {
		return nil, err
	}
	return maintenance, nil
}
