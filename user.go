package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

func getUserMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"me": {
			Name:               "me",
			Action:             "user/me",
			Method:             http.MethodGet,
			Result:             User{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listUsers": {
			Name:   "listUsers",
			Action: "users",
			Method: http.MethodGet,
			Result: []User{},
		},
		"updateUser": {
			Name:   "updateUser",
			Action: "users/%d",
			Method: http.MethodPut,
			Result: User{},
		},
	}
}

// User represents a registered account holder in the system.
// It contains the authentication identity (Login) and metadata.
type User struct {
	// ID is the unique identifier for the user.
	// This value is server-generated and read-only.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the user. Server-generated and read-only. Ignored during creation."`

	// Created indicates when the user was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Login is the unique username for the account.
	// This must be a valid email address.
	Login string `json:"login,omitempty" jsonschema:"The user's login name. Must be a valid email address (format: user@example.com)."`

	// Email is the user's contact email address.
	Email string `json:"email,omitempty" jsonschema:"The user's contact email address."`

	// Firstname is the user's given name.
	Firstname string `json:"firstname,omitempty" jsonschema:"The user's given name."`

	// Lastname is the user's family name.
	Lastname string `json:"lastname,omitempty" jsonschema:"The user's family name."`

	// PrimaryPhone is the user's primary phone number.
	PrimaryPhone string `json:"primaryPhone,omitempty" jsonschema:"The user's primary phone number."`

	// SecondaryPhone is the user's secondary phone number.
	SecondaryPhone string `json:"secondaryPhone,omitempty" jsonschema:"The user's secondary phone number."`

	// PreferredCommunicationLanguage is the language used to communicate with the user.
	PreferredCommunicationLanguage string `json:"preferredCommunicationLanguage,omitempty" jsonschema:"The user's preferred communication language."`

	// OrganizationID is the unique identifier of the organization the user belongs to.
	OrganizationID int `json:"organizationId,omitempty" jsonschema:"The unique identifier of the organization the user belongs to."`

	// OrganizationName is the display name of the user's organization.
	OrganizationName string `json:"organizationName,omitempty" jsonschema:"The display name of the user's organization."`

	// Active indicates whether the user account is currently enabled.
	// The update route is a full replace and writes this flag unconditionally, so
	// it is sent without omitempty: a false value must reach the API to deactivate
	// a user (otherwise the flag would be dropped from the payload).
	Active bool `json:"active" jsonschema:"Indicates whether the user account is currently enabled."`

	// Locked indicates whether the user account is locked, e.g. after repeated failed login attempts.
	// Sent without omitempty for the same reason as Active: the full-replace update
	// must be able to carry a false value.
	Locked bool `json:"locked" jsonschema:"Indicates whether the user account is locked (e.g. after failed login attempts)."`

	// Deleted indicates whether the user has been soft-deleted.
	// Sent without omitempty for the same reason as Active: the full-replace update
	// must be able to carry a false value.
	Deleted bool `json:"deleted" jsonschema:"Indicates whether the user has been soft-deleted."`

	// Agent indicates whether the user has agent privileges.
	// The API sends this flag as a string ("" or "1") instead of a JSON boolean,
	// which types.Bool decodes; see its documentation.
	Agent types.Bool `json:"agent,omitempty" jsonschema:"Indicates whether the user has agent (support staff) privileges."`

	// TfaEnabled indicates whether two-factor authentication is currently active for this user.
	TfaEnabled bool `json:"tfaEnabled,omitempty" jsonschema:"Indicates whether two-factor authentication is currently active for this user."`

	// TfaRequired indicates whether two-factor authentication is required for this user.
	TfaRequired bool `json:"tfaRequired,omitempty" jsonschema:"Indicates whether two-factor authentication is enforced for this user."`

	// IsIndirectCustomer indicates whether the user belongs to an indirect-customer organization.
	IsIndirectCustomer bool `json:"isIndirectCustomer,omitempty" jsonschema:"Indicates whether the user belongs to an indirect-customer organization."`

	// Admin indicates whether the user is an administrator within their organization.
	Admin bool `json:"admin,omitempty" jsonschema:"Indicates whether the user has administrator role in their organization."`

	// RootAdmin indicates whether the user is a root administrator with platform-wide access.
	RootAdmin bool `json:"rootAdmin,omitempty" jsonschema:"Indicates whether the user is a root administrator with platform-wide access."`

	// Roles lists the user's role assignments across all groups they are a member of.
	Roles []UserRole `json:"roles,omitempty" jsonschema:"The user's role assignments across all groups they are a member of."`

	// RootGroupRoles lists the user's role assignments restricted to root (top-level) groups.
	RootGroupRoles []UserRole `json:"rootGroupRoles,omitempty" jsonschema:"The user's role assignments restricted to root (top-level) groups."`
}

// UserRole represents a user's role assignment within a specific group.
type UserRole struct {
	// ID is the unique identifier for the role assignment.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the role assignment. Server-generated and read-only."`

	// Created indicates when the role assignment was created.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified records the last update time of the role assignment.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Server-managed, read-only."`

	// GroupID is the identifier of the group this role applies to.
	GroupID int `json:"groupId,omitempty" jsonschema:"The identifier of the group this role applies to."`

	// GroupName is the display name of the group this role applies to.
	GroupName string `json:"groupName,omitempty" jsonschema:"The display name of the group this role applies to."`

	// Role is the role identifier (e.g. ADMINISTRATOR or USER).
	Role string `json:"role,omitempty" jsonschema:"The role identifier. Allowed values: 'ADMINISTRATOR', 'USER'."`
}

// MeContext returns the active user information
func (api *API) MeContext(ctx context.Context) (*User, error) {
	if _, ok := api.methods["me"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "me")
	}

	definition := api.methods["me"]
	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*User)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// Me is equivalent to MeContext with context.Background().
//
// Deprecated: use MeContext.
func (api *API) Me() (*User, error) {
	return api.MeContext(context.Background())
}

// ListUsersContext returns the users visible to the authenticated account.
// An account holding the ADMINISTRATOR role in at least one group (root or
// sub group) sees every user of its organization, any other account sees only
// itself. Pass query parameters such as "page", "pageSize" or "search" via the
// params map.
//
// The list carries the base user fields only: Admin, RootAdmin, Roles and
// RootGroupRoles are not populated by this endpoint.
func (api *API) ListUsersContext(ctx context.Context, params map[string]string) ([]User, error) {
	if _, ok := api.methods["listUsers"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listUsers")
	}

	definition := api.methods["listUsers"]

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]User)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// UpdateUserContext updates the passed user using the MYRA API.
// The route is a full replace, so it is a read-modify-write: fetch the user via
// ListUsersContext, mutate the fields you want to change, and send the whole
// object back (including Modified for the optimistic-lock check). Omitted fields
// are cleared server-side.
func (api *API) UpdateUserContext(ctx context.Context, user *User) (*User, error) {
	if _, ok := api.methods["updateUser"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateUser")
	}

	definition := api.methods["updateUser"]
	definition.Action = fmt.Sprintf(definition.Action, user.ID)

	result, err := api.call(ctx, definition, user)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*User)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}
