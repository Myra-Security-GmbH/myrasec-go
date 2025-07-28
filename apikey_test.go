package myrasec

import "testing"

func TestListAPIKeys(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/user/me",
			`{"error":false, "violationList":[], "warningList":[], "data":[
				{"objectType":"UserExtendedVO", "id": 12345, "login":"test@example.com", "modified":"2025-07-28T15:39:12+0200", "created":"2025-01-09T16:31:13+0100"}
			]}`,
			methods["me"],
		),
		preCacheRequest(
			"https://apiv2.myracloud.com/user/12345/api-keys",
			`{"error":false,"violationList":[],"warningList":[],"data":[{"objectType":"ApiKeyVO","key":"aexuqu7ohc3ihahchair9aijooch4Poh","name":"myrasec-go","id":12345,"modified":"2025-04-02T10:15:49+0200","created":"2025-04-02T10:15:49+0200","deleted":false},{"objectType":"ApiKeyVO","key":"vouZaesiexael4me4ohchopee7Oe7qua","name":"terraform","id":67890,"modified":"2025-04-08T14:16:45+0200","created":"2025-04-08T14:16:45+0200","deleted":false}],"page":1,"count":2,"pageSize":50}`,
			methods["listApiKeys"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	apikeys, err := api.ListApiKeys(nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(apikeys) != 2 {
		t.Errorf("Expected to get [%d] API keys but got [%d]", 2, len(apikeys))
	}
}
