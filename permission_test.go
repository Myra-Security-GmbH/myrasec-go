package myrasec

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestListMyPermissions(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/user/permissions",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"id":1,"action":"READ","objectType":"Domain","objectPermissionType":"USER"},
				{"id":2,"action":"EDIT","objectType":"Domain","objectPermissionType":"GROUP","objectInstance":99}
			],"page":1,"count":2,"pageSize":50}`,
			"listMyPermissions",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	permissions, err := api.ListMyPermissions(nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(permissions) != 2 {
		t.Errorf("Expected to get [%d] permissions but got [%d]", 2, len(permissions))
	}

	if permissions[0].Action != PermissionActionRead {
		t.Errorf("Expected first permission action to be [%s] but got [%s]", PermissionActionRead, permissions[0].Action)
	}

	if permissions[1].ObjectInstance != 99 {
		t.Errorf("Expected second permission ObjectInstance to be [%d] but got [%d]", 99, permissions[1].ObjectInstance)
	}

	if permissions[1].ObjectPermissionType != PermissionTypeGroup {
		t.Errorf("Expected second permission ObjectPermissionType to be [%s] but got [%s]", PermissionTypeGroup, permissions[1].ObjectPermissionType)
	}
}

func TestListUserPermissions(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/user/12/permissions",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"id":5,"action":"READ","objectType":"Domain","objectPermissionType":"USER"}
			],"page":1,"count":1,"pageSize":50}`,
			"listUserPermissions",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	permissions, err := api.ListUserPermissions(12, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(permissions) != 1 {
		t.Errorf("Expected to get [%d] permission but got [%d]", 1, len(permissions))
	}

	if permissions[0].ID != 5 {
		t.Errorf("Expected permission ID to be [%d] but got [%d]", 5, permissions[0].ID)
	}
}

func TestListUserGroupPermissions(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/user/group/3/permissions",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"id":11,"action":"ADMIN","objectType":"UserGroup","objectPermissionType":"GROUP","scopes":[1,2,3]}
			],"page":1,"count":1,"pageSize":50}`,
			"listUserGroupPermissions",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	permissions, err := api.ListUserGroupPermissions(3, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(permissions) != 1 {
		t.Errorf("Expected to get [%d] permission but got [%d]", 1, len(permissions))
	}

	if permissions[0].Action != PermissionActionAdmin {
		t.Errorf("Expected permission action to be [%s] but got [%s]", PermissionActionAdmin, permissions[0].Action)
	}

	if len(permissions[0].Scopes) != 3 {
		t.Errorf("Expected permission to have [%d] scopes but got [%d]", 3, len(permissions[0].Scopes))
	}
}

func TestDecodePermissionCheckResponseAuthorized(t *testing.T) {
	methods := initializeMethods()

	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   io.NopCloser(bytes.NewBufferString(`{"error":false,"data":{"isAuthorized":true}}`)),
	}

	result, err := decodePermissionCheckResponse(&resp, methods["checkMyPermission"])
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	checkResult, ok := result.(*PermissionCheckResult)
	if !ok {
		t.Fatalf("Expected result to be *PermissionCheckResult but got %T", result)
	}

	if !checkResult.IsAuthorized {
		t.Error("Expected IsAuthorized to be true")
	}
}

func TestDecodePermissionCheckResponseUnauthorized(t *testing.T) {
	methods := initializeMethods()

	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   io.NopCloser(bytes.NewBufferString(`{"error":false,"data":{"isAuthorized":false}}`)),
	}

	result, err := decodePermissionCheckResponse(&resp, methods["checkMyPermission"])
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	checkResult, ok := result.(*PermissionCheckResult)
	if !ok {
		t.Fatalf("Expected result to be *PermissionCheckResult but got %T", result)
	}

	if checkResult.IsAuthorized {
		t.Error("Expected IsAuthorized to be false")
	}
}

func TestDecodePermissionCheckResponseError(t *testing.T) {
	methods := initializeMethods()

	resp := http.Response{
		Status: strconv.Itoa(http.StatusBadRequest),
		Body:   io.NopCloser(bytes.NewBufferString(`{"error":true,"violationList":[{"propertypath":"action","message":"invalid action"}]}`)),
	}

	_, err := decodePermissionCheckResponse(&resp, methods["checkMyPermission"])
	if err == nil {
		t.Error("Expected to get an error but got nil")
	}
}

func TestDecodePermissionCheckResponseInvalidBody(t *testing.T) {
	methods := initializeMethods()

	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   io.NopCloser(bytes.NewBufferString(`{"this will not work": func(){}}`)),
	}

	_, err := decodePermissionCheckResponse(&resp, methods["checkMyPermission"])
	if err == nil {
		t.Error("Expected to get an error but got nil")
	}
}

func TestAddUserPermission(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/user/12/permissions",
			`{"error":false,"violationList":[],"warningList":[],"targetObject":[
				{"id":77,"action":"READ","objectType":"Domain","objectPermissionType":"USER","objectInstance":12345}
			]}`,
			"addUserPermission",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	p, err := api.AddUserPermission(&ObjectPermission{
		Action:               PermissionActionRead,
		ObjectType:           "Domain",
		ObjectPermissionType: PermissionTypeUser,
		ObjectInstance:       12345,
	}, 12)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if p.ID != 77 {
		t.Errorf("Expected permission ID to be [%d] but got [%d]", 77, p.ID)
	}

	if p.ObjectInstance != 12345 {
		t.Errorf("Expected ObjectInstance to be [%d] but got [%d]", 12345, p.ObjectInstance)
	}
}

func TestCheckMyPermission(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/user/permissions/check",
			`{"error":false,"data":{"isAuthorized":true}}`,
			"checkMyPermission",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	result, err := api.CheckMyPermission(&ObjectPermission{
		Action:     PermissionActionEdit,
		ObjectType: "Domain",
	})
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if !result.IsAuthorized {
		t.Error("Expected IsAuthorized to be true")
	}
}

func TestAddUserGroupPermission(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/user/group/3/permissions",
			`{"error":false,"violationList":[],"warningList":[],"targetObject":[
				{"id":88,"action":"EDIT","objectType":"Domain","objectPermissionType":"GROUP"}
			]}`,
			"addUserGroupPermission",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	p, err := api.AddUserGroupPermission(&ObjectPermission{
		Action:               PermissionActionEdit,
		ObjectType:           "Domain",
		ObjectPermissionType: PermissionTypeGroup,
	}, 3)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if p.ID != 88 {
		t.Errorf("Expected permission ID to be [%d] but got [%d]", 88, p.ID)
	}

	if p.Action != PermissionActionEdit {
		t.Errorf("Expected Action to be [%s] but got [%s]", PermissionActionEdit, p.Action)
	}
}
