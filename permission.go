package myrasec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

const (
	// PermissionActionRead grants read access on an object.
	PermissionActionRead = "READ"

	// PermissionActionCreate grants the right to create new objects of the type.
	PermissionActionCreate = "CREATE"

	// PermissionActionEdit grants the right to modify existing objects.
	PermissionActionEdit = "EDIT"

	// PermissionActionSwitch grants the right to switch context (e.g. impersonate).
	PermissionActionSwitch = "SWITCH"

	// PermissionActionAdmin grants administrative rights on an object.
	PermissionActionAdmin = "ADMIN"

	// PermissionActionAgentMode grants access to agent (support) mode.
	PermissionActionAgentMode = "AGENT_MODE"

	// PermissionActionUserMode grants access to user mode.
	PermissionActionUserMode = "USER_MODE"

	// PermissionActionCacheClear grants the right to invalidate cached content.
	PermissionActionCacheClear = "CACHE_CLEAR"

	// PermissionActionReview grants the right to review pending changes.
	PermissionActionReview = "REVIEW"

	// PermissionActionPublish grants the right to publish reviewed changes.
	PermissionActionPublish = "PUBLISH"

	// PermissionActionLoadbalanceRead grants read access on load balancing data.
	PermissionActionLoadbalanceRead = "LOADBALANCE_READ"

	// PermissionActionLoadbalanceEdit grants modify access on load balancing data.
	PermissionActionLoadbalanceEdit = "LOADBALANCE_EDIT"

	// PermissionTypeUser identifies a permission attached to an individual user.
	PermissionTypeUser = "USER"

	// PermissionTypeGroup identifies a permission attached to a group.
	PermissionTypeGroup = "GROUP"
)

// getPermissionMethods returns permission related API calls.
func getPermissionMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listMyPermissions": {
			Name:   "listMyPermissions",
			Action: "user/permissions",
			Method: http.MethodGet,
			Result: []ObjectPermission{},
		},
		"checkMyPermission": {
			Name:               "checkMyPermission",
			Action:             "user/permissions/check",
			Method:             http.MethodPost,
			Result:             PermissionCheckResult{},
			ResponseDecodeFunc: decodePermissionCheckResponse,
		},
		"listUserPermissions": {
			Name:   "listUserPermissions",
			Action: "user/%d/permissions",
			Method: http.MethodGet,
			Result: []ObjectPermission{},
		},
		"addUserPermission": {
			Name:   "addUserPermission",
			Action: "user/%d/permissions",
			Method: http.MethodPost,
			Result: ObjectPermission{},
		},
		"removeUserPermission": {
			Name:   "removeUserPermission",
			Action: "user/%d/permissions/%d",
			Method: http.MethodDelete,
			Result: ObjectPermission{},
		},
		"listUserGroupPermissions": {
			Name:   "listUserGroupPermissions",
			Action: "user/group/%d/permissions",
			Method: http.MethodGet,
			Result: []ObjectPermission{},
		},
		"addUserGroupPermission": {
			Name:   "addUserGroupPermission",
			Action: "user/group/%d/permissions",
			Method: http.MethodPost,
			Result: ObjectPermission{},
		},
		"removeUserGroupPermission": {
			Name:   "removeUserGroupPermission",
			Action: "user/group/%d/permissions/%d",
			Method: http.MethodDelete,
			Result: ObjectPermission{},
		},
	}
}

// ObjectPermission represents a single object-level access rule granted to a
// user or a group. A permission combines an action (READ, EDIT, ...) with a
// target object type and, optionally, a concrete object instance.
type ObjectPermission struct {
	// ID is the unique identifier for the permission.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the permission. Server-generated; required for delete operations, ignored on create."`

	// Created indicates when the permission was created.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified records the last update time of the permission.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Server-managed, read-only."`

	// Action is the action being granted.
	// Allowed values: READ, CREATE, EDIT, SWITCH, ADMIN, AGENT_MODE, USER_MODE,
	// CACHE_CLEAR, REVIEW, PUBLISH, LOADBALANCE_READ, LOADBALANCE_EDIT.
	Action string `json:"action,omitempty" jsonschema:"The action being granted. Allowed values: 'READ', 'CREATE', 'EDIT', 'SWITCH', 'ADMIN', 'AGENT_MODE', 'USER_MODE', 'CACHE_CLEAR', 'REVIEW', 'PUBLISH', 'LOADBALANCE_READ', 'LOADBALANCE_EDIT'."`

	// ObjectType is the name of the entity type the permission applies to
	// (e.g. "Domain", "DnsRecord").
	ObjectType string `json:"objectType,omitempty" jsonschema:"The entity type the permission applies to (e.g. 'Domain', 'DnsRecord')."`

	// ObjectPermissionType identifies the recipient kind. Allowed values: USER, GROUP.
	ObjectPermissionType string `json:"objectPermissionType,omitempty" jsonschema:"The recipient kind of the permission. Allowed values: 'USER', 'GROUP'."`

	// ObjectInstance is the identifier of a specific object instance the
	// permission is scoped to. Zero or omitted means the permission applies
	// to all instances of ObjectType within the visible scope.
	ObjectInstance int `json:"objectInstance,omitempty" jsonschema:"The identifier of a specific object instance the permission is scoped to. Omit to apply to all instances of the object type within the visible scope."`

	// Name is an optional human-readable label for the permission.
	Name string `json:"name,omitempty" jsonschema:"An optional human-readable label for the permission."`

	// Scopes lists the identifiers of additional permissions that act as scope
	// restrictions for this permission.
	Scopes []int `json:"scopes,omitempty" jsonschema:"Identifiers of additional permissions that act as scope restrictions for this permission."`

	// Parents lists permissions this one inherits from.
	// This field is read-only and populated by the API.
	Parents []ObjectPermission `json:"parents,omitempty" jsonschema:"Permissions this one inherits from. Read-only, server-populated."`

	// Value is an optional auxiliary value carried with the permission.
	Value string `json:"value,omitempty" jsonschema:"An optional auxiliary value carried with the permission."`
}

// PermissionCheckResult represents the response of a permission check call.
type PermissionCheckResult struct {
	// IsAuthorized is true if the requesting user is allowed to perform the
	// requested action on the requested target.
	IsAuthorized bool `json:"isAuthorized" jsonschema:"True if the requesting user is allowed to perform the requested action on the requested target."`
}

// ListMyPermissions returns the permissions of the currently authenticated user.
// Supported query parameters include "language", "actions" and "objects".
func (api *API) ListMyPermissions(params map[string]string) ([]ObjectPermission, error) {
	if _, ok := api.methods["listMyPermissions"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listMyPermissions")
	}

	definition := api.methods["listMyPermissions"]

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]ObjectPermission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CheckMyPermission asks the API whether the authenticated user is allowed to
// perform the action described by the passed permission template on the
// referenced object type / instance.
func (api *API) CheckMyPermission(permission *ObjectPermission) (*PermissionCheckResult, error) {
	if _, ok := api.methods["checkMyPermission"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "checkMyPermission")
	}

	definition := api.methods["checkMyPermission"]

	result, err := api.call(definition, permission)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*PermissionCheckResult)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// ListUserPermissions returns all permissions assigned to the user with the
// given userID.
func (api *API) ListUserPermissions(userID int, params map[string]string) ([]ObjectPermission, error) {
	if _, ok := api.methods["listUserPermissions"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listUserPermissions")
	}

	definition := api.methods["listUserPermissions"]
	definition.Action = fmt.Sprintf(definition.Action, userID)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]ObjectPermission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// AddUserPermission grants the given permission to the user with the given userID.
func (api *API) AddUserPermission(permission *ObjectPermission, userID int) (*ObjectPermission, error) {
	if _, ok := api.methods["addUserPermission"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "addUserPermission")
	}

	definition := api.methods["addUserPermission"]
	definition.Action = fmt.Sprintf(definition.Action, userID)

	result, err := api.call(definition, permission)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*ObjectPermission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// RemoveUserPermission revokes the passed permission from the user identified
// by userID.
func (api *API) RemoveUserPermission(permission *ObjectPermission, userID int) (*ObjectPermission, error) {
	if _, ok := api.methods["removeUserPermission"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "removeUserPermission")
	}

	definition := api.methods["removeUserPermission"]
	definition.Action = fmt.Sprintf(definition.Action, userID, permission.ID)

	_, err := api.call(definition, permission)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

// ListUserGroupPermissions returns all permissions attached to the user group
// with the given groupID.
func (api *API) ListUserGroupPermissions(groupID int, params map[string]string) ([]ObjectPermission, error) {
	if _, ok := api.methods["listUserGroupPermissions"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listUserGroupPermissions")
	}

	definition := api.methods["listUserGroupPermissions"]
	definition.Action = fmt.Sprintf(definition.Action, groupID)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]ObjectPermission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// AddUserGroupPermission grants the given permission to the user group with the
// given groupID.
func (api *API) AddUserGroupPermission(permission *ObjectPermission, groupID int) (*ObjectPermission, error) {
	if _, ok := api.methods["addUserGroupPermission"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "addUserGroupPermission")
	}

	definition := api.methods["addUserGroupPermission"]
	definition.Action = fmt.Sprintf(definition.Action, groupID)

	result, err := api.call(definition, permission)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*ObjectPermission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// RemoveUserGroupPermission revokes the passed permission from the user group
// identified by groupID.
func (api *API) RemoveUserGroupPermission(permission *ObjectPermission, groupID int) (*ObjectPermission, error) {
	if _, ok := api.methods["removeUserGroupPermission"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "removeUserGroupPermission")
	}

	definition := api.methods["removeUserGroupPermission"]
	definition.Action = fmt.Sprintf(definition.Action, groupID, permission.ID)

	_, err := api.call(definition, permission)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

// decodePermissionCheckResponse decodes the response of POST /user/permissions/check.
// Unlike most endpoints, the API returns "data" as a JSON object ({"isAuthorized": bool})
// rather than an array, so the standard decoders cannot be used.
func decodePermissionCheckResponse(resp *http.Response, definition APIMethod) (any, error) {
	var raw struct {
		Error         bool            `json:"error"`
		ErrorMessage  string          `json:"errorMessage"`
		ViolationList []*Violation    `json:"violationList,omitempty"`
		Data          json.RawMessage `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	if raw.Error {
		return nil, newAPIErrorFromEnvelope(resp.StatusCode, raw.ErrorMessage, raw.ViolationList)
	}

	if definition.Result == nil {
		return raw.Data, nil
	}

	if len(raw.Data) == 0 {
		return nil, fmt.Errorf("empty data in API response")
	}

	retValue := reflect.New(reflect.TypeOf(definition.Result))
	res := retValue.Interface()
	if err := json.NewDecoder(bytes.NewReader(raw.Data)).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}
