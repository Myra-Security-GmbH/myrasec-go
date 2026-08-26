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

// Maintenance represents a scheduled maintenance window for a specific domain (FQDN).
// It controls when the maintenance page is displayed to visitors.
type Maintenance struct {
	// ID is the unique identifier for the maintenance entry.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the maintenance entry. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the maintenance entry was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Start specifies when the maintenance window begins.
	// If nil, the maintenance starts immediately.
	Start *types.DateTime `json:"start,omitempty" jsonschema:"The start timestamp (ISO 8601). If null, maintenance starts immediately. Must be earlier than 'End'."`

	// End specifies when the maintenance window finishes.
	// Must be later than Start.
	End *types.DateTime `json:"end,omitempty" jsonschema:"The end timestamp (ISO 8601). Must be later than 'Start'."`

	// Active indicates if the maintenance mode is currently live.
	// This is a computed read-only value based on Start and End.
	Active bool `json:"active" jsonschema:"Read-only status flag indicating if maintenance is currently active. Do not set this directly; adjust 'start' and 'end' to control activation."`

	// Content contains the HTML code for the maintenance page.
	// Note: Avoid linking to resources on the maintenance domain itself.
	Content string `json:"content" jsonschema:"The raw HTML content. Important: Do not link to resources (images/CSS) hosted on the domain being put into maintenance. Use external domains or inline Base64 encoding for assets."`

	// ContentFrom allows copying content from another existing maintenance page.
	// If specified, the system copies the HTML from the referenced FQDN.
	ContentFrom string `json:"contentFrom,omitempty" jsonschema:"Optional: A valid FQDN to copy maintenance content from. If set, the content is copied from the referenced domain's existing maintenance page."`

	// FQDN is the fully qualified domain name this maintenance applies to.
	FQDN string `json:"fqdn" jsonschema:"The Fully Qualified Domain Name (FQDN) to apply maintenance mode to (e.g., 'www.example.com')."`
}

// ListMaintenances returns a slice containing all maintenance pages for a subdomain
func (api *API) ListMaintenances(domainId int, subDomainName string, params map[string]string) ([]Maintenance, error) {
	if _, ok := api.methods["listMaintenances"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listMaintenances")
	}

	definition := api.methods["listMaintenances"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]Maintenance)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateMaintenance creates a new maintenance page for the passed subdomain (name) using the MYRA API
func (api *API) CreateMaintenance(maintenance *Maintenance, domainId int, subDomainName string) (*Maintenance, error) {
	if _, ok := api.methods["createMaintenance"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createMaintenance")
	}

	definition := api.methods["createMaintenance"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, maintenance)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Maintenance)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateMaintenance updates the passed maintenance page using the MYRA API
func (api *API) UpdateMaintenance(maintenance *Maintenance, domainId int, subDomainName string) (*Maintenance, error) {
	if _, ok := api.methods["updateMaintenance"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateMaintenance")
	}

	definition := api.methods["updateMaintenance"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, maintenance.ID)

	result, err := api.call(definition, maintenance)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Maintenance)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteMaintenance deletes the passed maintenance page using the MYRA API
func (api *API) DeleteMaintenance(maintenance *Maintenance, domainId int, subDomainName string) (*Maintenance, error) {
	if _, ok := api.methods["deleteMaintenance"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteMaintenance")
	}

	definition := api.methods["deleteMaintenance"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, maintenance.ID)

	_, err := api.call(definition, maintenance)
	if err != nil {
		return nil, err
	}
	return maintenance, nil
}
