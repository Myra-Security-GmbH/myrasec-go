package myrasec

import (
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

	// OrganizationID is the unique identifier of the organization the user belongs to.
	OrganizationID int `json:"organizationId,omitempty" jsonschema:"The unique identifier of the organization the user belongs to."`

	// OrganizationName is the display name of the user's organization.
	OrganizationName string `json:"organizationName,omitempty" jsonschema:"The display name of the user's organization."`

	// Active indicates whether the user account is currently enabled.
	Active bool `json:"active,omitempty" jsonschema:"Indicates whether the user account is currently enabled."`

	// Locked indicates whether the user account is locked, e.g. after repeated failed login attempts.
	Locked bool `json:"locked,omitempty" jsonschema:"Indicates whether the user account is locked (e.g. after failed login attempts)."`

	// Deleted indicates whether the user has been soft-deleted.
	Deleted bool `json:"deleted,omitempty" jsonschema:"Indicates whether the user has been soft-deleted."`

	// Agent indicates whether the user has agent privileges.
	Agent bool `json:"agent,omitempty" jsonschema:"Indicates whether the user has agent (support staff) privileges."`

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

// Me returns the active user information
func (api *API) Me() (*User, error) {
	if _, ok := methods["me"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "me")
	}

	definition := methods["me"]
	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*User)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}
