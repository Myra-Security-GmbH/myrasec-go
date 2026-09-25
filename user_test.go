package myrasec

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

func TestMe(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/user/me",
			`{"error":false, "violationList":[], "warningList":[], "data":[
				{"objectType":"UserExtendedVO", "id": 12345, "login":"test@example.com", "modified":"2025-07-28T15:39:12+0200", "created":"2025-01-09T16:31:13+0100",
					"admin": true, "active": true, "tfaEnabled": true,
					"primaryPhone":"+49111", "secondaryPhone":"+49222", "preferredCommunicationLanguage":"EN",
					"roles":[
						{"id": 10, "groupId": 1, "groupName": "root", "role": "ADMINISTRATOR"},
						{"id": 11, "groupId": 2, "groupName": "team", "role": "USER"}
					],
					"rootGroupRoles":[
						{"id": 10, "groupId": 1, "groupName": "root", "role": "ADMINISTRATOR"}
					]
				}
			]}`,
			"me",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	user, err := api.Me()
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if user.ID != 12345 {
		t.Errorf("Expected to get User with ID [%d] but got [%d]", 12345, user.ID)
	}

	if user.Login != "test@example.com" {
		t.Errorf("Expected to get User with ID [%s] but got [%s]", "test@example.com", user.Login)
	}

	if !user.Admin {
		t.Error("Expected user.Admin to be true")
	}

	if !user.TfaEnabled {
		t.Error("Expected user.TfaEnabled to be true")
	}

	if user.PrimaryPhone != "+49111" || user.SecondaryPhone != "+49222" || user.PreferredCommunicationLanguage != "EN" {
		t.Errorf("Expected the profile fields to decode, got primaryPhone [%s] secondaryPhone [%s] preferredCommunicationLanguage [%s]", user.PrimaryPhone, user.SecondaryPhone, user.PreferredCommunicationLanguage)
	}

	if len(user.Roles) != 2 {
		t.Errorf("Expected user to have [%d] roles but got [%d]", 2, len(user.Roles))
	}

	if user.Roles[0].GroupName != "root" || user.Roles[0].Role != GroupRoleAdministrator {
		t.Errorf("Expected first role to be ADMINISTRATOR on 'root' but got [%s/%s]", user.Roles[0].GroupName, user.Roles[0].Role)
	}

	if len(user.RootGroupRoles) != 1 {
		t.Errorf("Expected user to have [%d] root group role but got [%d]", 1, len(user.RootGroupRoles))
	}

	if user.RootGroupRoles[0].GroupID != 1 {
		t.Errorf("Expected root group role GroupID to be [%d] but got [%d]", 1, user.RootGroupRoles[0].GroupID)
	}
}

// TestMeWithAgentAsString covers the case where the API serializes the user's agent flag
// as a string ("" or "1") instead of a JSON boolean, which made every Me() call fail
// while decoding the response. The fixture of TestMe never carried the field, and
// preCacheRequest drops decoding errors, so the failure only showed against the
// live API.
func TestMeWithAgentAsString(t *testing.T) {
	tests := []struct {
		name     string
		agent    string
		expected bool
	}{
		{name: "empty string is not an agent", agent: `""`, expected: false},
		{name: "1 is an agent", agent: `"1"`, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, err := setupPreCachedAPI(preCacheRequest(
				"https://apiv2.myracloud.com/user/me",
				`{"error":false, "violationList":[], "warningList":[], "data":[
					{"objectType":"UserExtendedVO", "id": 12345, "login":"test@example.com", "firstname":"Test", "lastname":"User",
						"organizationId": 50, "organizationName": "Example Corp",
						"active": true, "locked": false, "deleted": false, "agent": `+tt.agent+`,
						"tfaEnabled": true, "tfaRequired": false, "admin": true}
				]}`,
				"me",
			))
			if err != nil {
				t.Error("Unexpected error.")
			}

			user, err := api.Me()
			if err != nil {
				t.Fatalf("Expected not to get an error but got [%s]", err.Error())
			}

			if bool(user.Agent) != tt.expected {
				t.Errorf("Expected user.Agent to be [%v] but got [%v]", tt.expected, user.Agent)
			}

			if user.Login != "test@example.com" || user.Firstname != "Test" || user.OrganizationName != "Example Corp" {
				t.Errorf("Expected the rest of the user to decode, got login [%s] firstname [%s] organization [%s]", user.Login, user.Firstname, user.OrganizationName)
			}

			if !user.Admin || !user.TfaEnabled {
				t.Error("Expected user.Admin and user.TfaEnabled to be true")
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /users": {Status: http.StatusOK, Body: `{"error":false,"violationList":[],"warningList":[],"data":[
			{"objectType":"UserVO","id":12345,"login":"admin@example.com","email":"admin@example.com","firstname":"Ada","lastname":"Admin",
				"organizationId":50,"organizationName":"Example Corp","active":true,"locked":false,"deleted":false,"agent":"",
				"tfaEnabled":true,"tfaRequired":false,"isIndirectCustomer":false,
				"created":"2025-01-09T16:31:13+0100","modified":"2025-07-28T15:39:12+0200"},
			{"objectType":"UserVO","id":12346,"login":"support@example.com","email":"support@example.com","firstname":"Sam","lastname":"Support",
				"organizationId":50,"organizationName":"Example Corp","active":false,"locked":true,"deleted":false,"agent":"1",
				"tfaEnabled":false,"tfaRequired":true,"isIndirectCustomer":true}
		],"page":2,"count":7,"pageSize":5}`},
	})

	users, err := api.ListUsersContext(context.Background(), map[string]string{
		ParamPage:     "2",
		ParamPageSize: "5",
		ParamSearch:   "example.com",
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	sent := requests.last(t)
	if sent.Method != http.MethodGet || sent.Path != "/users" {
		t.Errorf("Expected GET /users but got %s %s", sent.Method, sent.Path)
	}

	for key, expected := range map[string]string{ParamPage: "2", ParamPageSize: "5", ParamSearch: "example.com"} {
		if got := sent.Query.Get(key); got != expected {
			t.Errorf("Expected query parameter [%s] to be [%s] but got [%s]", key, expected, got)
		}
	}

	if len(users) != 2 {
		t.Fatalf("Expected to get [%d] users but got [%d]", 2, len(users))
	}

	if users[0].ID != 12345 || users[0].Login != "admin@example.com" {
		t.Errorf("Expected first user to be [12345/admin@example.com] but got [%d/%s]", users[0].ID, users[0].Login)
	}

	if users[0].Firstname != "Ada" || users[0].Lastname != "Admin" || users[0].OrganizationID != 50 || users[0].OrganizationName != "Example Corp" {
		t.Errorf("Expected the first user's profile fields to decode, got %+v", users[0])
	}

	if !users[0].Active || users[0].Locked || !users[0].TfaEnabled || bool(users[0].Agent) {
		t.Errorf("Expected the first user's flags to decode, got %+v", users[0])
	}

	if users[0].Created == nil || users[0].Modified == nil {
		t.Error("Expected the first user's Created and Modified to be set")
	}

	if users[1].Active || !users[1].Locked || !users[1].TfaRequired || !users[1].IsIndirectCustomer || !bool(users[1].Agent) {
		t.Errorf("Expected the second user's flags to decode, got %+v", users[1])
	}

	if users[0].Admin || users[0].RootAdmin || len(users[0].Roles) != 0 || len(users[0].RootGroupRoles) != 0 {
		t.Errorf("Expected the list not to carry role information, got %+v", users[0])
	}
}

func TestListUsersWithoutParams(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"GET /users": {Status: http.StatusOK, Body: `{"error":false,"violationList":[],"warningList":[],"data":[
			{"objectType":"UserVO","id":1,"login":"me@example.com"}
		],"page":1,"count":1,"pageSize":50}`},
	})

	users, err := api.ListUsersContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(users) != 1 || users[0].ID != 1 {
		t.Errorf("Expected a single user with ID [1] but got %+v", users)
	}

	if sent := requests.last(t); len(sent.Query) != 0 {
		t.Errorf("Expected no query parameters but got %v", sent.Query)
	}
}

func TestUpdateUser(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /users/12345": {Status: http.StatusOK, Body: `{"error":false,"violationList":[],"warningList":[],"data":[
			{"objectType":"UserVO","id":12345,"login":"test@example.com","firstname":"Test","lastname":"User",
				"primaryPhone":"+49111","secondaryPhone":"+49222","preferredCommunicationLanguage":"DE","active":false}
		]}`},
	})

	// The update route is a full replace and writes active/locked/deleted
	// unconditionally, so deactivating a user means sending active:false. Modified
	// travels back for the optimistic-lock check.
	modified := types.DateTime{Time: time.Date(2025, 7, 28, 15, 39, 12, 0, time.UTC)}
	updated, err := api.UpdateUserContext(context.Background(), &User{
		ID:                             12345,
		Login:                          "test@example.com",
		Firstname:                      "Test",
		Lastname:                       "User",
		PrimaryPhone:                   "+49111",
		SecondaryPhone:                 "+49222",
		PreferredCommunicationLanguage: "DE",
		Active:                         false,
		Modified:                       &modified,
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if updated.PreferredCommunicationLanguage != "DE" {
		t.Errorf("Expected the updated user's language to decode, got [%s]", updated.PreferredCommunicationLanguage)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodPut || sent.Path != "/users/12345" {
		t.Errorf("Expected PUT /users/12345 but got %s %s", sent.Method, sent.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	if payload["primaryPhone"] != "+49111" || payload["secondaryPhone"] != "+49222" || payload["preferredCommunicationLanguage"] != "DE" {
		t.Errorf("Expected the profile fields in the update payload, got %v", payload)
	}

	// active:false must reach the API; if 'active' were dropped (omitempty) a user
	// could never be deactivated through a full-replace update.
	active, present := payload["active"]
	if !present || active != false {
		t.Errorf("Expected the payload to carry active=false so the user can be deactivated, got %v", payload["active"])
	}

	// Modified must travel so the server can verify the version.
	if _, present := payload["modified"]; !present {
		t.Errorf("Expected the update payload to carry 'modified', got %v", payload)
	}
}

func TestListUsersForbidden(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"GET /users": {Status: http.StatusForbidden, Body: `{"error":true,"violationList":[{"propertyPath":"","message":"Access denied"}],"warningList":[],"data":[]}`},
	})

	users, err := api.ListUsersContext(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected an error for a forbidden response")
	}

	if users != nil {
		t.Errorf("Expected no users on error but got %+v", users)
	}
}
