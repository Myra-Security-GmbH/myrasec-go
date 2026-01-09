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

// MaintenanceTemplate represents a reusable HTML layout for maintenance pages.
// It defines the visual content displayed to users when a domain is switched to maintenance mode.
type MaintenanceTemplate struct {
	// ID is the unique identifier for the maintenance template.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the template. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the template was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Name is a descriptive label for the template.
	Name string `json:"name" jsonschema:"A descriptive name to identify this maintenance template."`

	// Content contains the HTML code for the maintenance page.
	// Note: Avoid linking to resources on the maintenance domain itself.
	Content string `json:"content" jsonschema:"The raw HTML content. Important: Do not link to resources (images/CSS) hosted on the domain being put into maintenance. Use external domains or inline Base64 encoding for assets."`
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
