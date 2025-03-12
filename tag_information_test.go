package myrasec

import (
	"testing"
)

func TestListTagInformation(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tags/1/information",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"id": 1, "key": "key1", "value": "value1", "comment": "comment1"}, 
				{"id": 2, "key": "key2", "value": "value2", "comment": "comment2"}, 
				{"id": 3, "key": "key3", "value": "value3", "comment": "comment3"}
			]}`,
			methods["listTagInformation"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	information, err := api.ListTagInformation(1, map[string]string{})
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(information) != 3 {
		t.Errorf("Expected to get [%d] information but got [%d]", 3, len(information))
	}

	for k, v := range information {
		if v.ID != k+1 {
			t.Errorf("Expected to get information with ID [%d] but got [%d]", k+1, v.ID)
		}

		if v.ID == 1 {
			if v.Key != "key1" {
				t.Errorf("Expected to get information with Key [%s] but got [%s]", "key1", v.Key)
			}

			if v.Value != "value1" {
				t.Errorf("Expected to get information with Value [%s] but got [%s]", "value1", v.Value)
			}

			if v.Comment != "comment1" {
				t.Errorf("Expected to get information with Comment [%s] but got [%s]", "comment1", v.Comment)
			}
		}

		if v.ID == 2 {
			if v.Key != "key2" {
				t.Errorf("Expected to get information with Key [%s] but got [%s]", "key2", v.Key)
			}

			if v.Value != "value2" {
				t.Errorf("Expected to get information with Value [%s] but got [%s]", "value2", v.Value)
			}

			if v.Comment != "comment2" {
				t.Errorf("Expected to get information with Comment [%s] but got [%s]", "comment2", v.Comment)
			}
		}

		if v.ID == 3 {
			if v.Key != "key3" {
				t.Errorf("Expected to get information with Key [%s] but got [%s]", "key3", v.Key)
			}

			if v.Value != "value3" {
				t.Errorf("Expected to get information with Value [%s] but got [%s]", "value3", v.Value)
			}

			if v.Comment != "comment3" {
				t.Errorf("Expected to get information with Comment [%s] but got [%s]", "comment3", v.Comment)
			}
		}
	}
}
