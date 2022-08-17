package myrasec

import (
	"testing"
)

func TestGetTag(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tags/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "name": "TagTest", "type": "TagType", "organization": 1, "assignments": [
					{"id": 1, "type": "DOMAIN", "title": "example.com"}
				]}
			]}`,
			methods["getTag"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	tag, err := api.GetTag(1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if tag.ID != 1 {
		t.Errorf("Expected to get Tag with ID [%d] but got [%d]", 1, tag.ID)
	}

	if tag.Name != "TagTest" {
		t.Errorf("Expected to get tag with name [%s] but got [%s]", "TestTag", tag.Name)
	}

	if tag.Type != "TagType" {
		t.Errorf("Expected to get tag with type [%s] but got [%s]", "TypeTag", tag.Type)
	}

	if tag.Organization != 1 {
		t.Errorf("Expected to get tag with organization id [%d] but got [%d]", 1, tag.Organization)
	}

	if len(tag.Assignments) != 1 {
		t.Errorf("Expected to get tag with number assignments [%d] but got [%d]", 1, len(tag.Assignments))
	}
}

func TestListTags(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/tags",
			`{"error": false, "pageSize": 10, "page": 1, "count": 3, "data": [
				{"id": 1, "name": "TagTest", "type": "TagType", "organization": "1", "assignments": [
					{"id": 1, "type": "DOMAIN", "title": "example.com"}
				]},
				{"id": 2, "name": "TagTest", "type": "TagType", "organization": "1", "assignments": [
					{"id": 1, "type": "DOMAIN", "title": "example.com"}
				]},
				{"id": 3, "name": "TagTest", "type": "TagType", "organization": "1", "assignments": [
					{"id": 1, "type": "DOMAIN", "title": "example.com"}
				]}
			]}`,
			methods["listTags"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	tags, err := api.ListTags(nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(tags) != 3 {
		t.Errorf("Expected to get [%d] tags but got [%d]", 3, len(tags))
	}
}
