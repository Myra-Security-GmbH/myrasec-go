package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

const (
	// UserGroupTypeUser identifies a regular user group.
	UserGroupTypeUser = "USER"

	// UserGroupTypeAgent identifies an agent (support staff) group.
	UserGroupTypeAgent = "AGENT"

	// GroupRoleAdministrator grants administrative privileges within a group.
	GroupRoleAdministrator = "ADMINISTRATOR"

	// GroupRoleUser grants standard membership within a group.
	GroupRoleUser = "USER"
)

// getUserGroupMethods returns user group related API calls.
func getUserGroupMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listUserGroups": {
			Name:   "listUserGroups",
			Action: "user/groups",
			Method: http.MethodGet,
			Result: []UserGroup{},
		},
		"getUserGroup": {
			Name:               "getUserGroup",
			Action:             "user/groups/%d",
			Method:             http.MethodGet,
			Result:             UserGroup{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"createUserGroup": {
			Name:   "createUserGroup",
			Action: "user/groups",
			Method: http.MethodPost,
			Result: UserGroup{},
		},
		"updateUserGroup": {
			Name:   "updateUserGroup",
			Action: "user/groups/%d",
			Method: http.MethodPut,
			Result: UserGroup{},
		},
		"deleteUserGroup": {
			Name:   "deleteUserGroup",
			Action: "user/groups/%d",
			Method: http.MethodDelete,
			Result: UserGroup{},
		},
		"listUsersFromGroup": {
			Name:   "listUsersFromGroup",
			Action: "user/group/%d/users",
			Method: http.MethodGet,
			Result: []User{},
		},
		"addUserToGroup": {
			Name:   "addUserToGroup",
			Action: "user/group/%d/users",
			Method: http.MethodPost,
			Result: GroupRole{},
		},
		"removeUserFromGroup": {
			Name:   "removeUserFromGroup",
			Action: "user/group/%d/users/%d",
			Method: http.MethodDelete,
			Result: User{},
		},
	}
}

// UserGroup represents a collection of users sharing common access rights and roles.
// Groups can be nested hierarchically: a child group inherits scope from its parent
// but exposes only its own members and resources.
type UserGroup struct {
	// ID is the unique identifier for the group.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the group. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the group was created.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Name is the human-readable display name of the group.
	Name string `json:"name,omitempty" jsonschema:"The display name of the group."`

	// Parent is the identifier of the parent group, or zero for a root group.
	Parent int `json:"parent,omitempty" jsonschema:"The identifier of the parent group. Zero or omitted for a root group."`

	// Children lists the immediate child groups nested under this group.
	// This field is read-only and populated by the API on list/read responses.
	Children []UserGroup `json:"children,omitempty" jsonschema:"The immediate child groups nested under this group. Read-only, server-populated."`

	// Roles lists the role identifiers the requesting user holds on this group.
	// This field is read-only.
	Roles []string `json:"roles,omitempty" jsonschema:"The role identifiers the requesting user holds on this group. Read-only."`

	// MembersCount is the number of users currently assigned to this group.
	// This field is read-only.
	MembersCount int `json:"membersCount,omitempty" jsonschema:"The number of users currently assigned to this group. Read-only."`

	// Type distinguishes between a regular user group and an agent group.
	// Allowed values: USER, AGENT. Defaults to USER on creation if omitted.
	Type string `json:"type,omitempty" jsonschema:"Distinguishes between a regular user group and an agent group. Allowed values: 'USER', 'AGENT'. Defaults to 'USER' on creation."`
}

// GroupRole represents the assignment of a user to a group with a specific role.
// It is the payload used when adding a user to a group.
type GroupRole struct {
	// ID is the unique identifier for the role assignment.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the role assignment. Server-generated and read-only."`

	// Created indicates when the role assignment was created.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified records the last update time of the role assignment.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Server-managed, read-only."`

	// UserID is the identifier of the user being assigned. Required.
	UserID int `json:"userId" jsonschema:"The identifier of the user being assigned to the group. Required."`

	// Role is the role identifier granted to the user within the group.
	// Allowed values: ADMINISTRATOR, USER. Required.
	Role string `json:"role" jsonschema:"The role identifier granted to the user within the group. Allowed values: 'ADMINISTRATOR', 'USER'. Required."`
}

// ListUserGroupsContext returns all user groups visible to the authenticated user.
// Pass query parameters such as "page", "pageSize" or "role" via the params map.
func (api *API) ListUserGroupsContext(ctx context.Context, params map[string]string) ([]UserGroup, error) {
	if _, ok := api.methods["listUserGroups"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listUserGroups")
	}

	definition := api.methods["listUserGroups"]

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]UserGroup)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListUserGroups is equivalent to ListUserGroupsContext with context.Background().
//
// Deprecated: use ListUserGroupsContext.
func (api *API) ListUserGroups(params map[string]string) ([]UserGroup, error) {
	return api.ListUserGroupsContext(context.Background(), params)
}

// GetUserGroupContext returns the user group identified by the given id.
func (api *API) GetUserGroupContext(ctx context.Context, id int) (*UserGroup, error) {
	if _, ok := api.methods["getUserGroup"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getUserGroup")
	}

	definition := api.methods["getUserGroup"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*UserGroup)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// GetUserGroup is equivalent to GetUserGroupContext with context.Background().
//
// Deprecated: use GetUserGroupContext.
func (api *API) GetUserGroup(id int) (*UserGroup, error) {
	return api.GetUserGroupContext(context.Background(), id)
}

// CreateUserGroupContext creates a new user group using the MYRA API.
func (api *API) CreateUserGroupContext(ctx context.Context, group *UserGroup) (*UserGroup, error) {
	if _, ok := api.methods["createUserGroup"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createUserGroup")
	}

	definition := api.methods["createUserGroup"]

	result, err := api.call(ctx, definition, group)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*UserGroup)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateUserGroup is equivalent to CreateUserGroupContext with context.Background().
//
// Deprecated: use CreateUserGroupContext.
func (api *API) CreateUserGroup(group *UserGroup) (*UserGroup, error) {
	return api.CreateUserGroupContext(context.Background(), group)
}

// UpdateUserGroupContext updates the passed user group using the MYRA API.
func (api *API) UpdateUserGroupContext(ctx context.Context, group *UserGroup) (*UserGroup, error) {
	if _, ok := api.methods["updateUserGroup"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateUserGroup")
	}

	definition := api.methods["updateUserGroup"]
	definition.Action = fmt.Sprintf(definition.Action, group.ID)

	result, err := api.call(ctx, definition, group)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*UserGroup)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateUserGroup is equivalent to UpdateUserGroupContext with context.Background().
//
// Deprecated: use UpdateUserGroupContext.
func (api *API) UpdateUserGroup(group *UserGroup) (*UserGroup, error) {
	return api.UpdateUserGroupContext(context.Background(), group)
}

// DeleteUserGroupContext deletes the passed user group using the MYRA API.
func (api *API) DeleteUserGroupContext(ctx context.Context, group *UserGroup) (*UserGroup, error) {
	if _, ok := api.methods["deleteUserGroup"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteUserGroup")
	}

	definition := api.methods["deleteUserGroup"]
	definition.Action = fmt.Sprintf(definition.Action, group.ID)

	_, err := api.call(ctx, definition, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// DeleteUserGroup is equivalent to DeleteUserGroupContext with context.Background().
//
// Deprecated: use DeleteUserGroupContext.
func (api *API) DeleteUserGroup(group *UserGroup) (*UserGroup, error) {
	return api.DeleteUserGroupContext(context.Background(), group)
}

// ListUsersFromGroupContext returns the users that are members of the given group.
// Supported query parameters include "search", "page", "pageSize", "language",
// "includeRoles" and "sort".
func (api *API) ListUsersFromGroupContext(ctx context.Context, groupID int, params map[string]string) ([]User, error) {
	if _, ok := api.methods["listUsersFromGroup"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listUsersFromGroup")
	}

	definition := api.methods["listUsersFromGroup"]
	definition.Action = fmt.Sprintf(definition.Action, groupID)

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

// ListUsersFromGroup is equivalent to ListUsersFromGroupContext with context.Background().
//
// Deprecated: use ListUsersFromGroupContext.
func (api *API) ListUsersFromGroup(groupID int, params map[string]string) ([]User, error) {
	return api.ListUsersFromGroupContext(context.Background(), groupID, params)
}

// AddUserToGroupContext assigns a user to the given group with the role specified in role.
func (api *API) AddUserToGroupContext(ctx context.Context, role *GroupRole, groupID int) (*GroupRole, error) {
	if _, ok := api.methods["addUserToGroup"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "addUserToGroup")
	}

	definition := api.methods["addUserToGroup"]
	definition.Action = fmt.Sprintf(definition.Action, groupID)

	result, err := api.call(ctx, definition, role)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*GroupRole)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// AddUserToGroup is equivalent to AddUserToGroupContext with context.Background().
//
// Deprecated: use AddUserToGroupContext.
func (api *API) AddUserToGroup(role *GroupRole, groupID int) (*GroupRole, error) {
	return api.AddUserToGroupContext(context.Background(), role, groupID)
}

// RemoveUserFromGroupContext removes the passed user from the group identified by groupID.
func (api *API) RemoveUserFromGroupContext(ctx context.Context, user *User, groupID int) (*User, error) {
	if _, ok := api.methods["removeUserFromGroup"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "removeUserFromGroup")
	}

	definition := api.methods["removeUserFromGroup"]
	definition.Action = fmt.Sprintf(definition.Action, groupID, user.ID)

	_, err := api.call(ctx, definition, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// RemoveUserFromGroup is equivalent to RemoveUserFromGroupContext with context.Background().
//
// Deprecated: use RemoveUserFromGroupContext.
func (api *API) RemoveUserFromGroup(user *User, groupID int) (*User, error) {
	return api.RemoveUserFromGroupContext(context.Background(), user, groupID)
}
