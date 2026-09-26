package myrasec

import (
	"testing"
)

func TestNormalizeIPFilterValue(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		filterType string
		expected   string
	}{
		{
			name:       "request limiter bare IPv4 gets /32",
			value:      "172.16.1.1",
			filterType: IPFilterTypeWhitelistRequestLimiter,
			expected:   "172.16.1.1/32",
		},
		{
			name:       "request limiter bare IPv4 loopback gets /32",
			value:      "127.0.0.1",
			filterType: IPFilterTypeWhitelistRequestLimiter,
			expected:   "127.0.0.1/32",
		},
		{
			name:       "request limiter bare IPv6 gets /128",
			value:      "dead::beef",
			filterType: IPFilterTypeWhitelistRequestLimiter,
			expected:   "dead::beef/128",
		},
		{
			name:       "request limiter IPv4-mapped IPv6 gets /128",
			value:      "::ffff:1.2.3.4",
			filterType: IPFilterTypeWhitelistRequestLimiter,
			expected:   "::ffff:1.2.3.4/128",
		},
		{
			name:       "request limiter already CIDR IPv4 is unchanged",
			value:      "10.0.0.0/24",
			filterType: IPFilterTypeWhitelistRequestLimiter,
			expected:   "10.0.0.0/24",
		},
		{
			name:       "request limiter already CIDR IPv6 is unchanged",
			value:      "dead::beef/128",
			filterType: IPFilterTypeWhitelistRequestLimiter,
			expected:   "dead::beef/128",
		},
		{
			name:       "request limiter non-IP value is unchanged",
			value:      "not-an-ip",
			filterType: IPFilterTypeWhitelistRequestLimiter,
			expected:   "not-an-ip",
		},
		{
			name:       "whitelist bare IPv4 is left untouched",
			value:      "1.2.3.4",
			filterType: IPFilterTypeWhitelist,
			expected:   "1.2.3.4",
		},
		{
			name:       "blacklist bare IPv4 is left untouched",
			value:      "1.2.3.4",
			filterType: IPFilterTypeBlacklist,
			expected:   "1.2.3.4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeIPFilterValue(tt.value, tt.filterType)
			if got != tt.expected {
				t.Errorf("normalizeIPFilterValue(%q, %q) = %q, expected %q", tt.value, tt.filterType, got, tt.expected)
			}
		})
	}
}
