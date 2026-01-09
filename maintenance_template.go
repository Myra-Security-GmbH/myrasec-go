package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getMaintenanceTemplateMethods returns Maintenance template related API calls
func getMaintenanceTemplateMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listMaintenanceTemplates": {
			Name:   "listMaintenanceTemplates",
			Action: "domain/%d/maintenance-templates",
			Method: http.MethodGet,
			Result: []MaintenanceTemplate{},
		},
		"createMaintenanceTemplate": {
			Name:   "createMaintenanceTemplate",
			Action: "domain/%d/maintenance-templates",
			Method: http.MethodPost,
			Result: MaintenanceTemplate{},
		},
		"updateMaintenanceTemplate": {
			Name:   "updateMaintenanceTemplate",
			Action: "domain/%d/maintenance-templates/%d",
			Method: http.MethodPut,
			Result: MaintenanceTemplate{},
		},
		"deleteMaintenanceTemplate": {
			Name:   "deleteMaintenanceTemplate",
			Action: "domain/%d/maintenance-templates/%d",
			Method: http.MethodDelete,
			Result: MaintenanceTemplate{},
		},
	}
}

// MaintenanceTemplate ...
type MaintenanceTemplate struct {
	ID       int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set, while inserting a new object. To update or delete a maintenance template it is necessary to add this attribute to your object."`
	Created  *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. It will be created by the server after creating a new Maintenance object. This value is only informational so it is not necessary to add this an attribute to any API call."`
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletions. This value is always a date type with an ISO 8601 format."`
	Name     string          `json:"name" jsonschema:"A name to identify your maintenance template."`
	Content  string          `json:"content" jsonschema:"HTML content to show as maintenance page. Please note that it is not possible to include resources from the domain you have set to maintenance mode. If your maintenance page contains images use a different domain or use inline base64 encoded images."`
}

// ListMaintenanceTemplates returns a slice containing all maintenance templates for a domain
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

	return *result.(*[]MaintenanceTemplate), nil
}

// CreateMaintenanceTemplate creates a new maintenance template for the passed domain (id) using the MYRA API
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

// UpdateMaintenanceTemplate updates the passed maintenance template using the MYRA API
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

// DeleteMaintenanceTemplate deletes the passed maintenance template using the MYRA API
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
