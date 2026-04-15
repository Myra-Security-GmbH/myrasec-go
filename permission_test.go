package myrasec

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestListMyPermissions(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/permissions",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"objectType":"ObjectPermissionVO","id":1,"action":"READ","objectType":"Domain","objectPermissionType":"USER"},
				{"objectType":"ObjectPermissionVO","id":2,"action":"EDIT","objectType":"Domain","objectPermissionType":"GROUP","objectInstance":99}
			],"page":1,"count":2,"pageSize":50}`,
			methods["listMyPermissions"],
		),
	})
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
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/12/permissions",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"objectType":"ObjectPermissionVO","id":5,"action":"READ","objectType":"Domain","objectPermissionType":"USER"}
			],"page":1,"count":1,"pageSize":50}`,
			methods["listUserPermissions"],
		),
	})
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
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/group/3/permissions",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"objectType":"ObjectPermissionVO","id":11,"action":"ADMIN","objectType":"UserGroup","objectPermissionType":"GROUP","scopes":[1,2,3]}
			],"page":1,"count":1,"pageSize":50}`,
			methods["listUserGroupPermissions"],
		),
	})
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
	resp := http.Response{
		Status: strconv.Itoa(http.StatusOK),
		Body:   io.NopCloser(bytes.NewBufferString(`{"this will not work": func(){}}`)),
	}

	_, err := decodePermissionCheckResponse(&resp, methods["checkMyPermission"])
	if err == nil {
		t.Error("Expected to get an error but got nil")
	}
}
