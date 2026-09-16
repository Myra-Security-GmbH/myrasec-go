package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getOrganizationNoteMethods returns OrganizationNote related API calls
func getOrganizationNoteMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getOrganizationNote": {
			Name:               "getOrganizationNote",
			Action:             "organization/notes",
			Method:             http.MethodGet,
			Result:             OrganizationNote{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"updateOrganizationNote": {
			Name:   "updateOrganizationNote",
			Action: "organization/notes",
			Method: http.MethodPut,
			Result: OrganizationNote{},
		},
	}
}

// OrganizationNote represents the free-form note of an organization.
// There is a single note per organization; both create and edit go through
// the same upsert (PUT) route, and reading an organization without a saved
// note returns an empty note shape.
type OrganizationNote struct {
	// ID is the unique identifier for the note.
	// This value is server-generated and read-only.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the note. Server-generated and read-only."`

	// Created indicates when the note was first saved.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified records the last update time in ISO 8601 format.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Server-managed, read-only."`

	// Notes is the free-form note text of the organization.
	Notes string `json:"notes,omitempty" jsonschema:"The free-form note text of the organization."`
}

// GetOrganizationNoteContext returns the note of the authenticated organization
func (api *API) GetOrganizationNoteContext(ctx context.Context) (*OrganizationNote, error) {
	if _, ok := api.methods["getOrganizationNote"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getOrganizationNote")
	}

	definition := api.methods["getOrganizationNote"]

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*OrganizationNote)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// GetOrganizationNote is equivalent to GetOrganizationNoteContext with context.Background().
//
// Deprecated: use GetOrganizationNoteContext.
func (api *API) GetOrganizationNote() (*OrganizationNote, error) {
	return api.GetOrganizationNoteContext(context.Background())
}

// UpdateOrganizationNoteContext saves (upserts) the note of the authenticated organization
func (api *API) UpdateOrganizationNoteContext(ctx context.Context, note *OrganizationNote) (*OrganizationNote, error) {
	if _, ok := api.methods["updateOrganizationNote"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateOrganizationNote")
	}

	definition := api.methods["updateOrganizationNote"]

	result, err := api.call(ctx, definition, note)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*OrganizationNote)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateOrganizationNote is equivalent to UpdateOrganizationNoteContext with context.Background().
//
// Deprecated: use UpdateOrganizationNoteContext.
func (api *API) UpdateOrganizationNote(note *OrganizationNote) (*OrganizationNote, error) {
	return api.UpdateOrganizationNoteContext(context.Background(), note)
}
