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
	ID          int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set, while inserting a new object. To update or delete a Maintenance it is necessary to add this attribute to your object."`
	Created     *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. It will be created by the server after creating a new Maintenance object. This value is only informational so it is not necessary to add this an attribute to any API call."`
	Modified    *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletions. This value is always a date type with an ISO 8601 format."`
	Start       *types.DateTime `json:"start,omitempty" jsonschema:"Start is a date type attribute with an ISO 8601 format. This attribute shows the start date for a maintenance. This date have to be lower than end or null to start now."`
	End         *types.DateTime `json:"end,omitempty" jsonschema:"nd is a date type attribute with an ISO 8601 format and shows the end date for a maintenance. This date have to be higher than start or null to end now."`
	Active      bool            `json:"active" jsonschema:"This information shows if this maintenance page is currently active. You cannot set this attribute directly instead you have to set start and end attribute to activate maintenance."`
	Content     string          `json:"content" jsonschema:"HTML content to show as maintenance page. Please note that it is not possible to include resources from the domain you have set to maintenance mode. If your maintenance page contains images use a different domain or use inline base64 encoded images."`
	ContentFrom string          `json:"contentFrom,omitempty" jsonschema:"This property can be used instead of the property content to reference an existing maintenance page’s content. Instead of sending the actual content, specify a valid FQDN here. This will copy the content from the referenced maintenance page to the newly created."`
	FQDN        string          `json:"fqdn" jsonschema:"Shows a FQDN (fully qualified domain name) for a maintenance. This attribute shows the domain to handle maintenance for."`
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

	return *result.(*[]Maintenance), nil
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
