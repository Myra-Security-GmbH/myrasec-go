package myrasec

import "testing"

func TestListBuckets(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://upload.myracloud.com/v2/bucket/list/example.com",
			`{"error": false, "result": [
				{"bucket": "b1", "linkedDomains": ["www.example.com"]},
				{"bucket": "b2", "linkedDomains": ["1.example.com", "2.example.com"]},
				{"bucket": "b3", "linkedDomains": []}
			]}`,
			methods["listBuckets"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	buckets, err := api.ListBuckets("example.com")
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(buckets) != 3 {
		t.Errorf("Expected to get [%d] buckets but got [%d]", 3, len(buckets))
	}

	for _, b := range buckets {
		switch b.Name {
		case "b1":
			if len(b.LinkedDomains) != 1 {
				t.Errorf("Expected to get [%d] linked domain got [%d]", 1, len(b.LinkedDomains))
			}
		case "b2":
			if len(b.LinkedDomains) != 2 {
				t.Errorf("Expected to get [%d] linked domain got [%d]", 2, len(b.LinkedDomains))
			}
		case "b3":
			if len(b.LinkedDomains) != 0 {
				t.Errorf("Expected to get [%d] linked domain got [%d]", 0, len(b.LinkedDomains))
			}
		default:
			t.Errorf("Unexpected bucket [%s]", b.Name)
		}
	}
}

func GetBucketStatus(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://upload.myracloud.com/v2/bucket/status/example.com/b1",
			`{
				"status": "Bucket available",
				"statusCode": 0
			}`,
			methods["getBucketStatus"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	status, err := api.GetBucketStatus("example.com", "b1")
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if status.Status != "Bucket available" {
		t.Errorf("Expected to get Status [%s] but got [%s]", "Bucket available", status.Status)
	}

	if status.StatusCode != 0 {
		t.Errorf("Expected to get StatusCode [%d] but got [%d]", 0, status.StatusCode)
	}
}

func GetBucketStatistics(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://upload.myracloud.com/v2/bucket/statistics/example.com/b1",
			`{"error": false, "result": {
				"files": 1, 
				"folders": 0, 
				"storageSize": 1048576, 
				"contentSize": 131072
			}`,
			methods["getBucketStatistics"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	statistics, err := api.GetBucketStatistics("example.com", "b1")
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if statistics.Files != 1 {
		t.Errorf("Expected to get Files [%d] but got [%d]", 1, statistics.Files)
	}

	if statistics.Folders != 0 {
		t.Errorf("Expected to get Folders [%d] but got [%d]", 0, statistics.Folders)
	}

	if statistics.StorageSize != 1048576 {
		t.Errorf("Expected to get StorageSize [%d] but got [%d]", 1048576, statistics.StorageSize)
	}

	if statistics.ContentSize != 131072 {
		t.Errorf("Expected to get ContentSize [%d] but got [%d]", 131072, statistics.ContentSize)
	}
}
