package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getOrganizationContactMethods returns OrganizationContact related API calls
func getOrganizationContactMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listOrganizationContacts": {
			Name:   "listOrganizationContacts",
			Action: "organization/contacts",
			Method: http.MethodGet,
			Result: []OrganizationContact{},
		},
		"createOrganizationContacts": {
			Name:               "createOrganizationContacts",
			Action:             "organization/contacts",
			Method:             http.MethodPost,
			Result:             OrganizationContact{},
			ResponseDecodeFunc: decodeOrganizationContactCreateResponse,
		},
		"updateOrganizationContact": {
			Name:   "updateOrganizationContact",
			Action: "organization/contacts/%d",
			Method: http.MethodPut,
			Result: OrganizationContact{},
		},
		"deleteOrganizationContact": {
			Name:   "deleteOrganizationContact",
			Action: "organization/contacts/%d",
			Method: http.MethodDelete,
			Result: OrganizationContact{},
		},
		"listOrganizationContactTypes": {
			Name:   "listOrganizationContactTypes",
			Action: "organization/contact-types",
			Method: http.MethodGet,
			Result: []OrganizationContactType{},
		},
	}
}

// decodeOrganizationContactCreateResponse decodes the response of the bulk
// create route. That endpoint answers with an envelope carrying an empty data
// list (no created objects, no targetObject), so the shared prepareResult would
// reject it with "empty Data in API response". This decoder only surfaces an
// error envelope (error:true / violationList) and otherwise reports success
// without a decoded body.
func decodeOrganizationContactCreateResponse(resp *http.Response, definition APIMethod) (any, error) {
	if _, err := decodeBaseResponse(resp); err != nil {
		return nil, err
	}
	return nil, nil
}

// OrganizationContact represents a contact person of an organization.
// The organization is resolved from the authenticated session, so contacts
// carry no domain or subdomain context.
type OrganizationContact struct {
	// ID is the unique identifier for the contact.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the contact. Server-generated; required for updates and deletes, ignored during creation."`

	// Created indicates when the contact was added.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. It is read-only on
	// create, but must be echoed back unchanged on update so the server can
	// detect concurrent modifications.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Read-only on create; must be echoed back on update (optimistic locking)."`

	// UserID links the contact to a user of the organization, if any.
	UserID int `json:"userId,omitempty" jsonschema:"The identifier of the linked user of the organization, if the contact is tied to a user account."`

	// Name is the display name of the contact.
	Name string `json:"name,omitempty" jsonschema:"The display name of the contact."`

	// Email is the contact's email address.
	Email string `json:"email" jsonschema:"The contact's email address."`

	// Phone is the contact's primary phone number.
	Phone string `json:"phone,omitempty" jsonschema:"The contact's primary phone number."`

	// SecondaryPhone is the contact's secondary phone number.
	SecondaryPhone string `json:"secondaryPhone,omitempty" jsonschema:"The contact's secondary phone number."`

	// PreferredCommunicationLanguage is the language used to communicate with the contact.
	// The API validates it against the supported values 'EN' and 'DE'.
	PreferredCommunicationLanguage string `json:"preferredCommunicationLanguage,omitempty" jsonschema:"The preferred communication language. Valid values: 'EN', 'DE'."`

	// Comments carries free-form notes about the contact.
	Comments string `json:"comments,omitempty" jsonschema:"Free-form notes about the contact."`

	// Types lists the contact-type keys assigned to this contact (multi-select).
	// Valid keys come from the contact-type catalog (ListOrganizationContactTypes).
	// The tag uses omitzero on purpose: a nil slice is omitted (leaving the
	// stored types unchanged), while an explicit empty slice serializes as [] to
	// remove all assigned types.
	Types []string `json:"types,omitzero" jsonschema:"The contact-type keys assigned to this contact (multi-select). Valid keys come from the contact-type catalog. Send an empty array to remove all types; omit to leave them unchanged."`

	// ReceiveSSLReminders controls whether the contact receives SSL expiry reminders.
	ReceiveSSLReminders bool `json:"receiveSSLReminders" jsonschema:"Indicates whether the contact receives SSL expiry reminders."`

	// ReceiveTrafficAlerts controls whether the contact receives traffic alerts.
	// It is a pointer on purpose: omitting it leaves the stored value unchanged.
	ReceiveTrafficAlerts *bool `json:"receiveTrafficAlerts,omitempty" jsonschema:"Indicates whether the contact receives traffic alerts. Omitting it leaves the stored value unchanged."`

	// TechnicalPriority is the technical contact priority (1-6).
	TechnicalPriority int `json:"technicalPriority,omitempty" jsonschema:"The technical contact priority. Valid range: 1-6."`

	// EscalationChain is the escalation priority (1-6). The REST key stays 'escalationPriority'.
	// The server only stores it when the contact also carries the 'ESCALATION'
	// type; submitting types without 'ESCALATION' resets the stored value to null.
	EscalationChain int `json:"escalationPriority,omitempty" jsonschema:"The escalation priority (escalation chain). Valid range: 1-6. Only stored when the contact also carries the 'ESCALATION' type."`
}

// OrganizationContactType represents an entry of the contact-type catalog.
type OrganizationContactType struct {
	// ID is the unique identifier for the contact type.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the contact type. Server-generated and read-only."`

	// Key is the machine-readable identifier used in OrganizationContact.Types.
	Key string `json:"key" jsonschema:"The machine-readable key used in a contact's 'types' multi-select."`

	// Label is the human-readable name of the contact type.
	Label string `json:"label" jsonschema:"The human-readable label of the contact type."`
}

// ListOrganizationContactsContext returns all contacts of the authenticated organization
func (api *API) ListOrganizationContactsContext(ctx context.Context, params map[string]string) ([]OrganizationContact, error) {
	if _, ok := api.methods["listOrganizationContacts"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listOrganizationContacts")
	}

	definition := api.methods["listOrganizationContacts"]

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]OrganizationContact)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateOrganizationContactsContext creates the passed contacts using the MYRA API.
// The endpoint is a bulk create, so the contacts are sent as a JSON array.
//
// The API acknowledges the create with an empty body: it returns no created
// objects, so the returned slice is the input echoed back and carries no
// server-generated ids. Call ListOrganizationContactsContext afterwards to
// obtain the persisted contacts with their ids.
func (api *API) CreateOrganizationContactsContext(ctx context.Context, contacts []OrganizationContact) ([]OrganizationContact, error) {
	if _, ok := api.methods["createOrganizationContacts"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createOrganizationContacts")
	}

	definition := api.methods["createOrganizationContacts"]

	if _, err := api.call(ctx, definition, contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

// UpdateOrganizationContactContext updates the passed contact using the MYRA API
func (api *API) UpdateOrganizationContactContext(ctx context.Context, contact *OrganizationContact) (*OrganizationContact, error) {
	if _, ok := api.methods["updateOrganizationContact"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateOrganizationContact")
	}

	definition := api.methods["updateOrganizationContact"]
	definition.Action = fmt.Sprintf(definition.Action, contact.ID)

	result, err := api.call(ctx, definition, contact)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*OrganizationContact)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteOrganizationContactContext deletes the passed contact using the MYRA API
func (api *API) DeleteOrganizationContactContext(ctx context.Context, contact *OrganizationContact) (*OrganizationContact, error) {
	if _, ok := api.methods["deleteOrganizationContact"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteOrganizationContact")
	}

	definition := api.methods["deleteOrganizationContact"]
	definition.Action = fmt.Sprintf(definition.Action, contact.ID)

	_, err := api.call(ctx, definition, contact)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

// ListOrganizationContactTypesContext returns the contact-type catalog of the organization
func (api *API) ListOrganizationContactTypesContext(ctx context.Context) ([]OrganizationContactType, error) {
	if _, ok := api.methods["listOrganizationContactTypes"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listOrganizationContactTypes")
	}

	definition := api.methods["listOrganizationContactTypes"]

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]OrganizationContactType)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}
