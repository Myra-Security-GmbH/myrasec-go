package myrasec

import (
	"testing"
)

func TestListTagRateLimits(t *testing.T) {

	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tag/1/ratelimits",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"id": 1, "network": "127.0.0.1/32", "value": 100, "burst": 50, "timeframe": 60, "tagId":2}, 
				{"id": 2, "network": "dead::beef/128", "value": 100, "burst": 50, "timeframe": 60, "tagId":2}, 
				{"id": 3, "network": "192.168.178.0/24", "value": 100, "burst": 50, "timeframe": 60, "tagId":2}
			]}`,
			methods["listTagRateLimits"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	limits, err := api.ListTagRateLimits(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(limits) != 3 {
		t.Errorf("Expected to fetch [%d] TagRateLimits but got [%d]", 3, len(limits))
	}

	for _, l := range limits {
		if !intInSlice(l.ID, []int{1, 2, 3}) {
			t.Errorf("Unexpected rate limit ID [%d]", l.ID)
		}

		if l.TagId != 2 {
			t.Errorf("Expected to get rate limit with tagId [%d] but got [%d]", 2, l.TagId)
		}

		if l.Value != 100 {
			t.Errorf("Expected to get rate limit with Value [%d] but got [%d]", 100, l.Value)
		}

		if l.Burst != 50 {
			t.Errorf("Expected to get rate limit with Burst [%d] but got [%d]", 50, l.Burst)
		}

		if l.Timeframe != 60 {
			t.Errorf("Expected to get rate limit with Timeframe [%d] but got [%d]", 60, l.Timeframe)
		}

		if l.ID == 1 && l.Network != "127.0.0.1/32" {
			t.Errorf("Expected to get rate limit with Network [%s] but got [%s]", "127.0.0.1/32", l.Network)
		}

		if l.ID == 2 && l.Network != "dead::beef/128" {
			t.Errorf("Expected to get rate limit with Network [%s] but got [%s]", "dead::beef/128", l.Network)
		}
		if l.ID == 3 && l.Network != "192.168.178.0/24" {
			t.Errorf("Expected to get rate limit with Network [%s] but got [%s]", "192.168.178.0/24", l.Network)
		}
	}

}
