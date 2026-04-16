package myrasec

import "testing"

func TestMe(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/me",
			`{"error":false, "violationList":[], "warningList":[], "data":[
				{"objectType":"UserExtendedVO", "id": 12345, "login":"test@example.com", "modified":"2025-07-28T15:39:12+0200", "created":"2025-01-09T16:31:13+0100",
					"admin": true, "active": true, "tfaEnabled": true,
					"roles":[
						{"id": 10, "groupId": 1, "groupName": "root", "role": "ADMINISTRATOR"},
						{"id": 11, "groupId": 2, "groupName": "team", "role": "USER"}
					],
					"rootGroupRoles":[
						{"id": 10, "groupId": 1, "groupName": "root", "role": "ADMINISTRATOR"}
					]
				}
			]}`,
			methods["me"],
		),
	})
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
