package myrasec

import (
	"testing"
)

func TestGetIPFilter(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/ip-filters/www.example.com/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "value": "127.0.0.1/32", "type": "WHITELIST", "expireDate": null, "enabled": true}
			]}`,
			methods["getIPFilter"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	filter, err := api.GetIPFilter(1, "www.example.com", 1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if filter.ID != 1 {
		t.Errorf("Expected to get IP filter with ID [%d] but got [%d]", 1, filter.ID)
	}

	if filter.Value != "127.0.0.1/32" {
		t.Errorf("Expected to get IP filter with Value [%s] but got [%s]", "127.0.0.1", filter.Value)
	}

	if filter.Type != "WHITELIST" {
		t.Errorf("Expected to get IP filter with Type [%s] but got [%s]", "WHITELIST", filter.Type)
	}

	if filter.ExpireDate != nil {
		t.Errorf("Expected to get IP filter without ExpireDate but got [%s]", filter.ExpireDate)
	}

	if !filter.Enabled {
		t.Errorf("Expected to get an enabled IP filter")
	}
}

func TestListIPFilters(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/ip-filters/www.example.com",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"id": 1, "value": "127.0.0.1/32", "type": "WHITELIST_REQUEST_LIMITER", "expireDate": null, "enabled": true}, 
				{"id": 2, "value": "dead::beef/128", "type": "BLACKLIST", "expireDate": "2022-06-10T22:00:00+0200", "enabled": true}, 
				{"id": 3, "value": "192.168.178.0/24", "type": "WHITELIST", "expireDate": null, "enabled": false}
			]}`,
			methods["listIPFilters"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	filters, err := api.ListIPFilters(1, "www.example.com", nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(filters) != 3 {
		t.Errorf("Expected to fetch [%d] IP filters but got [%d]", 3, len(filters))
	}

	for _, f := range filters {

		if !intInSlice(f.ID, []int{1, 2, 3}) {
			t.Errorf("Unexpected IP filter ID [%d]", f.ID)
		}

		if f.ID == 1 {
			if f.Type != "WHITELIST_REQUEST_LIMITER" {
				t.Errorf("Expected to get IP filter with Type [%s] but got [%s]", "WHITELIST_REQUEST_LIMITER", f.Type)
			}

			if f.Value != "127.0.0.1/32" {
				t.Errorf("Expected to get IP filter with Value [%s] but got [%s]", "192.168.178.0/24", f.Value)
			}

			if f.ExpireDate != nil {
				t.Errorf("Expected to get IP filter without ExpireDate but got [%s]", f.ExpireDate)
			}

			if !f.Enabled {
				t.Errorf("Expected to get an enabled IP filter")
			}
		}

		if f.ID == 2 {
			if f.Type != "BLACKLIST" {
				t.Errorf("Expected to get IP filter with Type [%s] but got [%s]", "BLACKLIST", f.Type)
			}

			if f.Value != "dead::beef/128" {
				t.Errorf("Expected to get IP filter with Value [%s] but got [%s]", "dead::beef/128", f.Value)
			}

			if f.ExpireDate == nil {
				t.Errorf("Expected to get IP filter with ExpireDate")
			}
			if f.ExpireDate.Format("2006-01-02") != "2022-06-10" {
				t.Errorf("Expected to get IP filter without ExpireDate [%s] but got [%s]", "2022-06-10", f.ExpireDate.Format("2006-01-02"))
			}

			if !f.Enabled {
				t.Errorf("Expected to get an enabled IP filter")
			}
		}

		if f.ID == 3 {
			if f.Type != "WHITELIST" {
				t.Errorf("Expected to get IP filter with Type [%s] but got [%s]", "WHITELIST", f.Type)
			}

			if f.Value != "192.168.178.0/24" {
				t.Errorf("Expected to get IP filter with Value [%s] but got [%s]", "192.168.178.0/24", f.Value)
			}

			if f.ExpireDate != nil {
				t.Errorf("Expected to get IP filter without ExpireDate but got [%s]", f.ExpireDate)
			}

			if f.Enabled {
				t.Errorf("Expected to get a disabled IP filter")
			}
		}
	}

}
