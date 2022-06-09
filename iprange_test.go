package myrasec

import (
	"testing"
)

func TestListIPRanges(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/ip-ranges",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"id": 1, "network": "127.0.0.1/32", "validFrom": "2022-01-01T00:00:00+0200", "validTo": null}, 
				{"id": 2, "network": "192.168.0.0/16", "validFrom": "2022-01-01T00:00:00+0200", "validTo": null}, 
				{"id": 3, "network": "dead::beef/128", "validFrom": "2022-01-01T00:00:00+0200", "validTo": "2022-12-31T00:00:00+0200"}
			]}`,
			methods["listIPRanges"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	ranges, err := api.ListIPRanges(nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(ranges) != 3 {
		t.Errorf("Expected to fetch [%d] IP ranges but got [%d]", 3, len(ranges))
	}

	for _, r := range ranges {
		if !intInSlice(r.ID, []int{1, 2, 3}) {
			t.Errorf("Unexpected IP range ID [%d]", r.ID)
		}

		if r.ValidFrom == nil {
			t.Errorf("Expected to have a ValidFrom date but none is set")
		}

		if r.ValidFrom.Format("2006-01-02") != "2022-01-01" {
			t.Errorf("Expected to have ValidFrom date [%s] but got [%s]", "2022-01-01", r.ValidFrom.Format("2006-01-02"))
		}

		if r.ID == 1 {
			if r.Network != "127.0.0.1/32" {
				t.Errorf("Expected to get IP range with Network [%s] but got [%s]", "127.0.0.1/32", r.Network)
			}

			if r.ValidTo != nil {
				t.Errorf("Expected not to have a ValidTo date but got [%s]", r.ValidTo)
			}
		}

		if r.ID == 2 {
			if r.Network != "192.168.0.0/16" {
				t.Errorf("Expected to get IP range with Network [%s] but got [%s]", "192.168.0.0/16", r.Network)
			}

			if r.ValidTo != nil {
				t.Errorf("Expected not to have a ValidTo date but got [%s]", r.ValidTo)
			}
		}

		if r.ID == 3 {
			if r.Network != "dead::beef/128" {
				t.Errorf("Expected to get IP range with Network [%s] but got [%s]", "dead::beef/128", r.Network)
			}

			if r.ValidTo == nil {
				t.Errorf("Expected to have a ValidTo date but none is set")
			}

			if r.ValidTo.Format("2006-01-02") != "2022-12-31" {
				t.Errorf("Expected to have ValidTo date [%s] but got [%s]", "2022-12-31", r.ValidTo.Format("2006-01-02"))
			}
		}
	}
}
