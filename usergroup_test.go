package myrasec

import "testing"

func TestListUserGroups(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/groups",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"objectType":"GroupVO","id":1,"name":"root","type":"USER","membersCount":3,"created":"2025-01-09T16:31:13+0100","modified":"2025-04-02T10:15:49+0200","children":[
					{"objectType":"GroupVO","id":2,"name":"child","parent":1,"type":"USER","membersCount":1}
				]},
				{"objectType":"GroupVO","id":3,"name":"another","type":"USER","membersCount":0}
			],"page":1,"count":2,"pageSize":50}`,
			methods["listUserGroups"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	groups, err := api.ListUserGroups(nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(groups) != 2 {
		t.Errorf("Expected to get [%d] user groups but got [%d]", 2, len(groups))
	}

	if groups[0].ID != 1 {
		t.Errorf("Expected first group ID to be [%d] but got [%d]", 1, groups[0].ID)
	}

	if groups[0].Name != "root" {
		t.Errorf("Expected first group name to be [%s] but got [%s]", "root", groups[0].Name)
	}

	if groups[0].Type != UserGroupTypeUser {
		t.Errorf("Expected first group type to be [%s] but got [%s]", UserGroupTypeUser, groups[0].Type)
	}

	if groups[0].MembersCount != 3 {
		t.Errorf("Expected first group MembersCount to be [%d] but got [%d]", 3, groups[0].MembersCount)
	}

	if len(groups[0].Children) != 1 {
		t.Errorf("Expected first group to have [%d] child but got [%d]", 1, len(groups[0].Children))
	}

	if groups[0].Children[0].Parent != 1 {
		t.Errorf("Expected child parent to be [%d] but got [%d]", 1, groups[0].Children[0].Parent)
	}
}

func TestGetUserGroup(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/groups/42",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"objectType":"GroupVO","id":42,"name":"engineering","type":"USER","membersCount":7,"roles":["ADMINISTRATOR"]}
			]}`,
			methods["getUserGroup"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	group, err := api.GetUserGroup(42)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if group.ID != 42 {
		t.Errorf("Expected group ID to be [%d] but got [%d]", 42, group.ID)
	}

	if group.Name != "engineering" {
		t.Errorf("Expected group name to be [%s] but got [%s]", "engineering", group.Name)
	}

	if len(group.Roles) != 1 || group.Roles[0] != GroupRoleAdministrator {
		t.Errorf("Expected group roles to contain [%s] but got %v", GroupRoleAdministrator, group.Roles)
	}
}

func TestListUsersFromGroup(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/group/7/users",
			`{"error":false,"violationList":[],"warningList":[],"data":[
				{"objectType":"UserVO","id":100,"login":"alice@example.com","email":"alice@example.com","firstname":"Alice","lastname":"Example","organizationId":50,"organizationName":"Example Corp","active":true,"agent":false,"admin":true},
				{"objectType":"UserVO","id":101,"login":"bob@example.com","firstname":"Bob","active":true}
			],"page":1,"count":2,"pageSize":50}`,
			methods["listUsersFromGroup"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	users, err := api.ListUsersFromGroup(7, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(users) != 2 {
		t.Errorf("Expected to get [%d] users but got [%d]", 2, len(users))
	}

	if users[0].ID != 100 {
		t.Errorf("Expected first user ID to be [%d] but got [%d]", 100, users[0].ID)
	}

	if users[0].Firstname != "Alice" {
		t.Errorf("Expected first user firstname to be [%s] but got [%s]", "Alice", users[0].Firstname)
	}

	if users[0].OrganizationID != 50 {
		t.Errorf("Expected first user OrganizationID to be [%d] but got [%d]", 50, users[0].OrganizationID)
	}

	if !users[0].Admin {
		t.Errorf("Expected first user Admin flag to be true")
	}
}
