package myrasec

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDomainVerifiedIsDecoded(t *testing.T) {
	api, _ := newTestAPI(t, map[string]testResponse{
		"GET /domains/1": {Status: http.StatusOK, Body: `{"error": false, "data": [{"id": 1, "name": "shop.example.com", "verified": false}]}`},
		"GET /domains/2": {Status: http.StatusOK, Body: `{"error": false, "data": [{"id": 2, "name": "example.com", "verified": true}]}`},
		"GET /domains/3": {Status: http.StatusOK, Body: `{"error": false, "data": [{"id": 3, "name": "example.org"}]}`},
	})

	cases := []struct {
		id       int
		verified *bool
		pending  bool
	}{
		{1, boolPtr(false), true},
		{2, boolPtr(true), false},
		// an API without domain verification sends no attribute: not pending
		{3, nil, false},
	}

	for _, c := range cases {
		domain, err := api.GetDomainContext(context.Background(), c.id)
		if err != nil {
			t.Fatalf("Expected not to get an error for domain [%d] but got [%s]", c.id, err.Error())
		}

		if (domain.Verified == nil) != (c.verified == nil) || (domain.Verified != nil && *domain.Verified != *c.verified) {
			t.Errorf("Domain [%d]: unexpected Verified [%v]", c.id, domain.Verified)
		}

		if domain.IsPendingVerification() != c.pending {
			t.Errorf("Domain [%d]: expected IsPendingVerification [%t]", c.id, c.pending)
		}
	}
}

func TestDomainVerifiedIsNotSentWhenUnset(t *testing.T) {
	body, err := json.Marshal(&Domain{Name: "example.com"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, ok := decoded["verified"]; ok {
		t.Errorf("Expected a create payload without verified, got %s", body)
	}
}

func boolPtr(v bool) *bool {
	return &v
}
