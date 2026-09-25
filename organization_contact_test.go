package myrasec

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

func TestListOrganizationContacts(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/organization/contacts",
			`{"error":false, "violationList":[], "warningList":[], "pageSize":50, "page":1, "count":2, "data":[
				{"objectType":"OrganizationContactVO", "id":1, "userId":10, "name":"Alice Admin", "email":"alice@example.com",
					"phone":"+49111", "secondaryPhone":"+49222", "preferredCommunicationLanguage":"EN", "comments":"primary",
					"types":["SALES","COMPLIANCE","ESCALATION"], "receiveSSLReminders":true, "receiveTrafficAlerts":true,
					"technicalPriority":1, "escalationPriority":2,
					"created":"2025-01-09T16:31:13+0100", "modified":"2025-07-28T15:39:12+0200"},
				{"objectType":"OrganizationContactVO", "id":2, "email":"bob@example.com", "types":["COMPLIANCE"],
					"receiveSSLReminders":false, "receiveTrafficAlerts":false}
			]}`,
			"listOrganizationContacts",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	contacts, err := api.ListOrganizationContactsContext(context.Background(), nil)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(contacts) != 2 {
		t.Fatalf("Expected to get [%d] contacts but got [%d]", 2, len(contacts))
	}

	first := contacts[0]
	if first.ID != 1 || first.UserID != 10 || first.Name != "Alice Admin" || first.Email != "alice@example.com" {
		t.Errorf("Expected the first contact identity fields to decode, got %+v", first)
	}

	if first.Phone != "+49111" || first.SecondaryPhone != "+49222" || first.PreferredCommunicationLanguage != "EN" || first.Comments != "primary" {
		t.Errorf("Expected the first contact profile fields to decode, got %+v", first)
	}

	if len(first.Types) != 3 || first.Types[0] != "SALES" || first.Types[1] != "COMPLIANCE" || first.Types[2] != "ESCALATION" {
		t.Errorf("Expected the multi-select types to decode, got %v", first.Types)
	}

	if !first.ReceiveSSLReminders {
		t.Error("Expected first.ReceiveSSLReminders to be true")
	}

	if first.ReceiveTrafficAlerts == nil || !*first.ReceiveTrafficAlerts {
		t.Errorf("Expected first.ReceiveTrafficAlerts to decode to true, got %v", first.ReceiveTrafficAlerts)
	}

	if first.TechnicalPriority != 1 || first.EscalationChain != 2 {
		t.Errorf("Expected priorities to decode (technicalPriority=1, escalationPriority=2), got %d/%d", first.TechnicalPriority, first.EscalationChain)
	}

	if first.Created == nil || first.Modified == nil {
		t.Error("Expected the first contact's Created and Modified to be set")
	}

	second := contacts[1]
	if second.ID != 2 || second.Email != "bob@example.com" || len(second.Types) != 1 || second.Types[0] != "COMPLIANCE" {
		t.Errorf("Expected the second contact to decode, got %+v", second)
	}

	if second.ReceiveTrafficAlerts == nil || *second.ReceiveTrafficAlerts {
		t.Errorf("Expected second.ReceiveTrafficAlerts to decode to false, got %v", second.ReceiveTrafficAlerts)
	}
}

func TestCreateOrganizationContacts(t *testing.T) {
	// The bulk-create route acknowledges with an empty data list: no created
	// objects and no targetObject. The client must treat that as success rather
	// than raising "empty Data in API response".
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /organization/contacts": {Status: http.StatusCreated, Body: `{"error":false, "violationList":[], "warningList":[], "data":[]}`},
	})

	alerts := true
	input := []OrganizationContact{
		{
			Email:                          "new@example.com",
			Name:                           "New Contact",
			SecondaryPhone:                 "+49333",
			PreferredCommunicationLanguage: "DE",
			Comments:                       "created",
			Types:                          []string{"SALES", "ESCALATION"},
			ReceiveSSLReminders:            true,
			ReceiveTrafficAlerts:           &alerts,
			TechnicalPriority:              3,
			EscalationChain:                4,
		},
	}

	created, err := api.CreateOrganizationContactsContext(context.Background(), input)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	// The API returns no objects, so the call echoes the input back; the contacts
	// carry no server-generated ids and callers must re-list to obtain them.
	if len(created) != 1 || created[0].Email != "new@example.com" || created[0].EscalationChain != 4 {
		t.Errorf("Expected the input contacts to be returned unchanged, got %+v", created)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodPost || sent.Path != "/organization/contacts" {
		t.Errorf("Expected POST /organization/contacts but got %s %s", sent.Method, sent.Path)
	}

	// The bulk-create endpoint expects a JSON array of contacts.
	var payload []map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON array payload but got [%s]", sent.Body)
	}

	if len(payload) != 1 {
		t.Fatalf("Expected the payload to carry a single contact but got [%d]", len(payload))
	}

	entry := payload[0]
	if entry["email"] != "new@example.com" || entry["secondaryPhone"] != "+49333" || entry["preferredCommunicationLanguage"] != "DE" || entry["comments"] != "created" {
		t.Errorf("Expected the contact fields to round-trip in the payload, got %v", entry)
	}

	if entry["escalationPriority"] != float64(4) || entry["technicalPriority"] != float64(3) {
		t.Errorf("Expected escalationPriority=4 and technicalPriority=3 in the payload, got %v", entry)
	}

	types, ok := entry["types"].([]any)
	if !ok || len(types) != 2 || types[0] != "SALES" || types[1] != "ESCALATION" {
		t.Errorf("Expected the multi-select field to serialize as 'types', got %v", entry["types"])
	}

	// The multi-select must serialize as "types", never under a "contact"+"Type" key.
	forbiddenKey := "contact" + "Type"
	if _, present := entry[forbiddenKey]; present {
		t.Errorf("Expected no %q key in the payload; the multi-select is 'types'", forbiddenKey)
	}
}

// The Types multi-select uses omitzero: a nil slice must be absent from the
// payload (leaving the stored types unchanged) while an explicit empty slice
// must serialize as [] so the server removes all assigned types.
func TestCreateOrganizationContactsTypesOmitzero(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /organization/contacts": {Status: http.StatusCreated, Body: `{"error":false, "violationList":[], "warningList":[], "data":[]}`},
	})

	// nil Types must be omitted from the payload.
	if _, err := api.CreateOrganizationContactsContext(context.Background(), []OrganizationContact{
		{Email: "nil@example.com"},
	}); err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	var nilPayload []map[string]any
	if err := json.Unmarshal(requests.last(t).Body, &nilPayload); err != nil {
		t.Fatalf("Expected a JSON array payload but got [%s]", requests.last(t).Body)
	}
	if _, present := nilPayload[0]["types"]; present {
		t.Errorf("Expected a nil Types to be omitted from the payload, got %v", nilPayload[0])
	}

	// An explicit empty slice must serialize as [] to clear the assigned types.
	if _, err := api.CreateOrganizationContactsContext(context.Background(), []OrganizationContact{
		{Email: "empty@example.com", Types: []string{}},
	}); err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	var emptyPayload []map[string]any
	if err := json.Unmarshal(requests.last(t).Body, &emptyPayload); err != nil {
		t.Fatalf("Expected a JSON array payload but got [%s]", requests.last(t).Body)
	}
	types, present := emptyPayload[0]["types"]
	if !present {
		t.Fatalf("Expected an explicit empty Types to be present as [] in the payload, got %v", emptyPayload[0])
	}
	if arr, ok := types.([]any); !ok || len(arr) != 0 {
		t.Errorf("Expected types to serialize as an empty array, got %v", types)
	}
}

func TestUpdateOrganizationContact(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /organization/contacts/7": {Status: http.StatusOK, Body: `{"error":false, "violationList":[], "warningList":[], "data":[
			{"objectType":"OrganizationContactVO", "id":7, "email":"new@example.com", "name":"Renamed", "types":["COMPLIANCE"],
				"modified":"2026-08-02T09:00:00+0200"}
		]}`},
	})

	// The contact PUT is a full replace, so a real caller fetches the contact and
	// sends the whole object back, including Modified for the optimistic-lock check.
	modified := types.DateTime{Time: time.Date(2026, 8, 1, 9, 0, 0, 0, time.UTC)}
	updated, err := api.UpdateOrganizationContactContext(context.Background(), &OrganizationContact{
		ID:       7,
		Email:    "new@example.com",
		Name:     "Renamed",
		Types:    []string{"COMPLIANCE"},
		Modified: &modified,
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if updated.Name != "Renamed" {
		t.Errorf("Expected the updated contact to be renamed, got [%s]", updated.Name)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodPut || sent.Path != "/organization/contacts/7" {
		t.Errorf("Expected PUT /organization/contacts/7 but got %s %s", sent.Method, sent.Path)
	}

	// Modified must travel in the payload so the server can verify the version.
	var payload map[string]any
	if err := json.Unmarshal(sent.Body, &payload); err != nil {
		t.Fatalf("Expected a JSON payload but got [%s]", sent.Body)
	}
	if _, present := payload["modified"]; !present {
		t.Errorf("Expected the update payload to carry 'modified', got %v", payload)
	}
}

func TestDeleteOrganizationContact(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"DELETE /organization/contacts/7": {Status: http.StatusOK, Body: `{"error":false, "violationList":[], "warningList":[], "data":[]}`},
	})

	deleted, err := api.DeleteOrganizationContactContext(context.Background(), &OrganizationContact{ID: 7})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if deleted.ID != 7 {
		t.Errorf("Expected the deleted contact to keep its ID, got %+v", deleted)
	}

	sent := requests.last(t)
	if sent.Method != http.MethodDelete || sent.Path != "/organization/contacts/7" {
		t.Errorf("Expected DELETE /organization/contacts/7 but got %s %s", sent.Method, sent.Path)
	}
}

func TestListOrganizationContactTypes(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/organization/contact-types",
			`{"error":false, "violationList":[], "warningList":[], "data":[
				{"id":1, "key":"SALES", "label":"Sales"},
				{"id":2, "key":"COMPLIANCE", "label":"Compliance"}
			]}`,
			"listOrganizationContactTypes",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	catalog, err := api.ListOrganizationContactTypesContext(context.Background())
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(catalog) != 2 {
		t.Fatalf("Expected to get [%d] contact types but got [%d]", 2, len(catalog))
	}

	if catalog[0].ID != 1 || catalog[0].Key != "SALES" || catalog[0].Label != "Sales" {
		t.Errorf("Expected the first contact type to decode id/key/label, got %+v", catalog[0])
	}

	if catalog[1].Key != "COMPLIANCE" || catalog[1].Label != "Compliance" {
		t.Errorf("Expected the second contact type to decode, got %+v", catalog[1])
	}
}
