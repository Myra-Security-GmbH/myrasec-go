package myrasec

import (
	"testing"
)

func TestNormalizeIPFilterValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "bare IPv4 gets /32",
			value:    "172.16.1.1",
			expected: "172.16.1.1/32",
		},
		{
			name:     "bare IPv4 loopback gets /32",
			value:    "127.0.0.1",
			expected: "127.0.0.1/32",
		},
		{
			name:     "bare IPv6 gets /128",
			value:    "dead::beef",
			expected: "dead::beef/128",
		},
		{
			name:     "already CIDR IPv4 is unchanged",
			value:    "10.0.0.0/24",
			expected: "10.0.0.0/24",
		},
		{
			name:     "already CIDR IPv6 is unchanged",
			value:    "dead::beef/128",
			expected: "dead::beef/128",
		},
		{
			name:     "non-IP value is unchanged",
			value:    "not-an-ip",
			expected: "not-an-ip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeIPFilterValue(tt.value)
			if got != tt.expected {
				t.Errorf("normalizeIPFilterValue(%q) = %q, expected %q", tt.value, got, tt.expected)
			}
		})
	}
}
