package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// MaintenanceTemplate ...
//
type MaintenanceTemplate struct {
	ID       int             `json:"id,omitempty"`
	Created  *types.DateTime `json:"created,omitempty"`
	Modified *types.DateTime `json:"modified,omitempty"`
	Name     string          `json:"name"`
	Content  string          `json:"content"`
}

//
// ListMaintenanceTemplates returns a slice containing all maintenance templates for a domain
//
func (api *API) ListMaintenanceTemplates(domainId int, params map[string]string) ([]MaintenanceTemplate, error) {
	if _, ok := methods["listMaintenanceTemplates"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listMaintenanceTemplates")
	}

	definition := methods["listMaintenanceTemplates"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []MaintenanceTemplate
	records = append(records, *result.(*[]MaintenanceTemplate)...)

	return records, nil
}

//
// CreateMaintenanceTemplate creates a new maintenance template for the passed domain (id) using the MYRA API
//
func (api *API) CreateMaintenanceTemplate(template *MaintenanceTemplate, domainId int) (*MaintenanceTemplate, error) {
	if _, ok := methods["createMaintenanceTemplate"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createMaintenanceTemplate")
	}

	definition := methods["createMaintenanceTemplate"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, template)
	if err != nil {
		return nil, err
	}
	return result.(*MaintenanceTemplate), nil
}

//
// UpdateMaintenanceTemplate updates the passed maintenance template using the MYRA API
//
func (api *API) UpdateMaintenanceTemplate(template *MaintenanceTemplate, domainId int) (*MaintenanceTemplate, error) {
	if _, ok := methods["updateMaintenanceTemplate"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateMaintenanceTemplate")
	}

	definition := methods["updateMaintenanceTemplate"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, template.ID)

	result, err := api.call(definition, template)
	if err != nil {
		return nil, err
	}
	return result.(*MaintenanceTemplate), nil
}

//
// DeleteMaintenanceTemplate deletes the passed maintenance template using the MYRA API
//
func (api *API) DeleteMaintenanceTemplate(template *MaintenanceTemplate, domainId int) (*MaintenanceTemplate, error) {
	if _, ok := methods["deleteMaintenanceTemplate"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteMaintenanceTemplate")
	}

	definition := methods["deleteMaintenanceTemplate"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, template.ID)

	_, err := api.call(definition, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}
