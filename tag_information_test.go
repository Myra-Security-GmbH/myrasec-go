package myrasec

import (
	"testing"
)

func TestGetTagInformation(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tag/information/1",
			`{"error": false, "violationList": [], "warningList": [], "data": [{
            "objectType": "TagInformationVO",
            "name": "waf tag",
            "type": "INFORMATION",
            "id": 1,
            "modified": "2024-08-29T11:34:31+0200",
            "created": "2024-08-29T11:34:31+0200",
            "label": "waf tag",
            "informations": { "first": "some value", "second": "another value" }
        }
    ]
}`,
			methods["getTagInformation"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	tag, err := api.GetTagInformation(1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	info, exists := tag["informations"]
	if !exists {
		t.Error("Expected that attribute `information` exsist, but does not.")
	}

	informations := info.(map[string]interface{})

	first, exists := informations["first"]
	if !exists {
		t.Error("Expeted that information `first` exists")
	}
	if first != "some value" {
		t.Errorf("Expected first value to be `some value`, but got [%s]", first)
	}

	second, exists := informations["second"]
	if !exists {
		t.Error("Expeted that information `second` exists")
	}
	if second != "another value" {
		t.Errorf("Expected second value to be `another value`, but got [%s]", second)
	}
}
