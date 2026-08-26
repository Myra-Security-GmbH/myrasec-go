package myrasec

import (
	"testing"
)

func TestListWaitingRooms(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/waiting-rooms?domainId=1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"objectType":"WaitingRoomVO","name":"test name","maxConcurrent":300,"sessionTimeout":200,"waitRefresh":100,"paths":["test-path"],"subDomainName":"www.example.com","vhostId":1,"id":1,"modified":"2025-01-14T11:15:09+0100","created":"2025-01-14T11:15:09+0100"}
			]}`,
			"listWaitingRoomsForDomain",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	waitingrooms, err := api.ListWaitingRoomsForDomain(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(waitingrooms) != 1 {
		t.Errorf("Expected to get [%d] waiting rooms but got [%d]", 1, len(waitingrooms))
	}

	for _, m := range waitingrooms {
		if m.ID != 1 {
			t.Errorf("Expected to get WaitingRoom with ID [%d] but got [%d]", 1, m.ID)
		}
	}
}
