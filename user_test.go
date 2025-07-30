package myrasec

import "testing"

func TestMe(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/me",
			`{"error":false, "violationList":[], "warningList":[], "data":[
				{"objectType":"UserExtendedVO", "id": 12345, "login":"test@example.com", "modified":"2025-07-28T15:39:12+0200", "created":"2025-01-09T16:31:13+0100"}
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
}
