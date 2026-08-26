package types

import (
	"encoding/json"
	"testing"
)

func TestBoolUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		expectErr bool
		expected  Bool
	}{
		{name: "json true", json: `true`, expected: true},
		{name: "json false", json: `false`, expected: false},
		{name: "empty string is false", json: `""`, expected: false},
		{name: "string 0", json: `"0"`, expected: false},
		{name: "string 1", json: `"1"`, expected: true},
		{name: "string false", json: `"false"`, expected: false},
		{name: "string true", json: `"true"`, expected: true},
		{name: "string True mixed case", json: `"True"`, expected: true},
		{name: "string no", json: `"no"`, expected: false},
		{name: "string yes", json: `"yes"`, expected: true},
		{name: "string off", json: `"off"`, expected: false},
		{name: "string ON", json: `"ON"`, expected: true},
		{name: "null is false", json: `null`, expected: false},
		{name: "number 0", json: `0`, expected: false},
		{name: "number 1", json: `1`, expected: true},
		{name: "unknown string", json: `"maybe"`, expectErr: true},
		{name: "string with whitespace", json: `" 1"`, expectErr: true},
		{name: "number 2", json: `2`, expectErr: true},
		{name: "array", json: `[]`, expectErr: true},
		{name: "object", json: `{}`, expectErr: true},
		{name: "invalid json", json: `tru`, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start from true so a decoder that leaves the value untouched is caught.
			b := Bool(true)
			err := json.Unmarshal([]byte(tt.json), &b)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("Expected an error for [%s] but got none (value %v)", tt.json, b)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected not to get an error for [%s] but got [%s]", tt.json, err.Error())
			}

			if b != tt.expected {
				t.Errorf("Expected [%s] to decode to [%v] but got [%v]", tt.json, tt.expected, b)
			}
		})
	}
}

// TestBoolInStruct covers the way the type is used on the API models: a flag
// with omitempty that arrives as a string, and is written back as a JSON boolean.
func TestBoolInStruct(t *testing.T) {
	type flags struct {
		Agent Bool `json:"agent,omitempty"`
		Admin Bool `json:"admin,omitempty"`
	}

	var decoded flags
	if err := json.Unmarshal([]byte(`{"agent":"1","admin":false}`), &decoded); err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if !decoded.Agent {
		t.Error("Expected Agent to be true")
	}

	if decoded.Admin {
		t.Error("Expected Admin to be false")
	}

	encoded, err := json.Marshal(decoded)
	if err != nil {
		t.Fatalf("Expected not to get an error but got [%s]", err.Error())
	}

	if string(encoded) != `{"agent":true}` {
		t.Errorf("Expected a plain JSON boolean with the false flag omitted but got [%s]", encoded)
	}
}
