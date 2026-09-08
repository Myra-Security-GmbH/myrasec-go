package myrasec

import (
	"net"
	"strings"
)

// IP filter types supported by the MYRA API.
const (
	IPFilterTypeBlacklist               = "BLACKLIST"
	IPFilterTypeWhitelist               = "WHITELIST"
	IPFilterTypeWhitelistRequestLimiter = "WHITELIST_REQUEST_LIMITER"
)

// normalizeIPFilterValue normalizes a bare single IP address to its host CIDR
// notation (/32 for IPv4, /128 for IPv6) before it is sent to the API.
//
// Only WHITELIST_REQUEST_LIMITER rejects a bare single IP without CIDR notation;
// WHITELIST and BLACKLIST accept it, so their values are returned untouched to
// avoid changing what the API echoes back to callers (which would otherwise cause
// perpetual Terraform plan drift for existing bare-IP resources).
//
// Values that already contain a CIDR suffix, or that are not valid IP addresses,
// are returned unchanged so the API can validate them.
func normalizeIPFilterValue(value, filterType string) string {
	if filterType != IPFilterTypeWhitelistRequestLimiter {
		return value
	}

	if strings.Contains(value, "/") {
		return value
	}

	if net.ParseIP(value) == nil {
		return value
	}

	// Decide the host mask from the textual form: an IPv6 literal always contains
	// a colon. Relying on the text (rather than To4()) keeps IPv4-mapped IPv6
	// addresses such as "::ffff:1.2.3.4" on /128 instead of /32.
	if strings.Contains(value, ":") {
		return value + "/128"
	}

	return value + "/32"
}
