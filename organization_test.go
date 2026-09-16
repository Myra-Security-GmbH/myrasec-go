package myrasec

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetOrganizationNote(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/organization/notes",
			`{"error":false, "violationList":[], "warningList":[], "data":[
				{"objectType":"OrganizationNoteVO", "id":3, "notes":"remember the maintenance window",
					"created":"2025-01-09T16:31:13+0100", "modified":"2025-07-28T15:39:12+0200"}
			]}`,
			"getOrganizationNote",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	note, err := api.GetOrganizationNoteContext(context.Background())
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if note.ID != 3 || note.Notes != "remember the maintenance window" {
		t.Errorf("Expected the note to decode id and notes, got %+v", note)
	}

	if note.Created == nil || note.Modified == nil {
		t.Error("Expected the note's Created and Modified to be set")
	}
}

func TestUpdateOrganizationNote(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /organization/notes": {Status: http.StatusOK, Body: `{"error":false, "violationList":[], "warningList":[], "data":[
			{"objectType":"OrganizationNoteVO", "id":3, "notes":"updated note"}
		]}`},
	})

	updated, err := api.UpdateOrganizationNoteContext(context.Background(), &OrganizationNote{
		ID:    3,
		Notes: "updated note",
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if updated.Notes != "updated note" {
		t.Errorf("Expected the updated note to carry the new text, got [%s]", updated.Notes)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodPut || sent.Path != "/organization/notes" {
		t.Errorf("Expected PUT /organization/notes but got %s %s", sent.Method, sent.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	if payload["notes"] != "updated note" {
		t.Errorf("Expected the notes value in the payload, got %v", payload)
	}
}

// Clearing the note sends an empty string; the field must be present so the API
// stores the cleared value rather than leaving the previous one untouched.
func TestUpdateOrganizationNoteClearing(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /organization/notes": {Status: http.StatusOK, Body: `{"error":false, "violationList":[], "warningList":[], "data":[
			{"objectType":"OrganizationNoteVO", "id":3, "notes":""}
		]}`},
	})

	updated, err := api.UpdateOrganizationNoteContext(context.Background(), &OrganizationNote{ID: 3, Notes: ""})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if updated.Notes != "" {
		t.Errorf("Expected the cleared note to decode to an empty string, got [%s]", updated.Notes)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodPut || sent.Path != "/organization/notes" {
		t.Errorf("Expected PUT /organization/notes but got %s %s", sent.Method, sent.Path)
	}

	// The cleared value must be carried in the payload; if "notes" were dropped
	// (e.g. via omitempty) the upsert would leave the stored note untouched.
	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}

	notes, present := payload["notes"]
	if !present || notes != "" {
		t.Errorf("Expected the payload to carry notes=\"\" so the note is cleared, got %v", payload)
	}
}
