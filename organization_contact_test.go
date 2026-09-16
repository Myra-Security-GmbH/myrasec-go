package myrasec

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListOrganizationContacts(t *testing.T) {
	api, err := setupPreCachedAPI(
		preCacheRequest(
			"https://apiv2.myracloud.com/organization/contacts",
			`{"error":false, "violationList":[], "warningList":[], "pageSize":50, "page":1, "count":2, "data":[
				{"objectType":"OrganizationContactVO", "id":1, "userId":10, "name":"Alice Admin", "email":"alice@example.com",
					"phone":"+49111", "secondaryPhone":"+49222", "preferredCommunicationLanguage":"EN", "comments":"primary",
					"types":["technical","billing"], "receiveSSLReminders":true, "receiveTrafficAlerts":true,
					"technicalPriority":1, "escalationPriority":2,
					"created":"2025-01-09T16:31:13+0100", "modified":"2025-07-28T15:39:12+0200"},
				{"objectType":"OrganizationContactVO", "id":2, "email":"bob@example.com", "types":["security"],
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

	if len(first.Types) != 2 || first.Types[0] != "technical" || first.Types[1] != "billing" {
		t.Errorf("Expected the multi-select types to decode, got %v", first.Types)
	}

	if !first.ReceiveSSLReminders {
		t.Error("Expected first.ReceiveSSLReminders to be true")
	}

	if first.ReceiveTrafficAlerts == nil || !*first.ReceiveTrafficAlerts {
		t.Errorf("Expected first.ReceiveTrafficAlerts to decode to true, got %v", first.ReceiveTrafficAlerts)
	}

	if first.TechnicalPriority != 1 || first.EscalationChain != 2 {
		t.Errorf("Expected priorities to decode (technical=1, escalation=2), got %d/%d", first.TechnicalPriority, first.EscalationChain)
	}

	if first.Created == nil || first.Modified == nil {
		t.Error("Expected the first contact's Created and Modified to be set")
	}

	second := contacts[1]
	if second.ID != 2 || second.Email != "bob@example.com" || len(second.Types) != 1 || second.Types[0] != "security" {
		t.Errorf("Expected the second contact to decode, got %+v", second)
	}

	if second.ReceiveTrafficAlerts == nil || *second.ReceiveTrafficAlerts {
		t.Errorf("Expected second.ReceiveTrafficAlerts to decode to false, got %v", second.ReceiveTrafficAlerts)
	}
}

func TestCreateOrganizationContacts(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"POST /organization/contacts": {Status: http.StatusCreated, Body: `{"error":false, "violationList":[], "warningList":[], "targetObject":[
			{"objectType":"OrganizationContactVO", "id":7, "email":"new@example.com", "name":"New Contact",
				"secondaryPhone":"+49333", "preferredCommunicationLanguage":"DE", "comments":"created",
				"types":["technical","billing"], "receiveSSLReminders":true, "receiveTrafficAlerts":true,
				"technicalPriority":3, "escalationPriority":4,
				"created":"2026-08-01T09:00:00+0200", "modified":"2026-08-01T09:00:00+0200"}
		]}`},
	})

	alerts := true
	created, err := api.CreateOrganizationContactsContext(context.Background(), []OrganizationContact{
		{
			Email:                          "new@example.com",
			Name:                           "New Contact",
			SecondaryPhone:                 "+49333",
			PreferredCommunicationLanguage: "DE",
			Comments:                       "created",
			Types:                          []string{"technical", "billing"},
			ReceiveSSLReminders:            true,
			ReceiveTrafficAlerts:           &alerts,
			TechnicalPriority:              3,
			EscalationChain:                4,
		},
	})
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if created.ID != 7 || created.Email != "new@example.com" || created.EscalationChain != 4 {
		t.Errorf("Expected the created contact to be decoded from the response, got %+v", created)
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
	if !ok || len(types) != 2 || types[0] != "technical" || types[1] != "billing" {
		t.Errorf("Expected the multi-select field to serialize as 'types', got %v", entry["types"])
	}

	// The multi-select must serialize as "types", never under a "contact"+"Type" key.
	forbiddenKey := "contact" + "Type"
	if _, present := entry[forbiddenKey]; present {
		t.Errorf("Expected no %q key in the payload; the multi-select is 'types'", forbiddenKey)
	}
}

func TestUpdateOrganizationContact(t *testing.T) {
	api, requests := newTestAPI(t, map[string]testResponse{
		"PUT /organization/contacts/7": {Status: http.StatusOK, Body: `{"error":false, "violationList":[], "warningList":[], "data":[
			{"objectType":"OrganizationContactVO", "id":7, "email":"new@example.com", "name":"Renamed", "types":["billing"]}
		]}`},
	})

	updated, err := api.UpdateOrganizationContactContext(context.Background(), &OrganizationContact{
		ID:    7,
		Email: "new@example.com",
		Name:  "Renamed",
		Types: []string{"billing"},
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
				{"id":1, "key":"technical", "label":"Technical"},
				{"id":2, "key":"billing", "label":"Billing"}
			]}`,
			"listOrganizationContactTypes",
		),
	)
	if err != nil {
		t.Error("Unexpected error.")
	}

	types, err := api.ListOrganizationContactTypesContext(context.Background())
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(types) != 2 {
		t.Fatalf("Expected to get [%d] contact types but got [%d]", 2, len(types))
	}

	if types[0].ID != 1 || types[0].Key != "technical" || types[0].Label != "Technical" {
		t.Errorf("Expected the first contact type to decode id/key/label, got %+v", types[0])
	}

	if types[1].Key != "billing" || types[1].Label != "Billing" {
		t.Errorf("Expected the second contact type to decode, got %+v", types[1])
	}
}
