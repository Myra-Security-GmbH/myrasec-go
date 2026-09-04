package myrasec

import (
	"net"
	"strings"
)

// normalizeIPFilterValue normalizes a bare single IP address to its host CIDR
// notation (/32 for IPv4, /128 for IPv6) before it is sent to the API. This keeps
// the stored value consistent across all filter types, in particular
// WHITELIST_REQUEST_LIMITER, which rejects a single IP without CIDR notation.
// Values that already contain a CIDR suffix or that are not valid IP addresses are
// returned unchanged so the API can validate them.
func normalizeIPFilterValue(value string) string {
	if strings.Contains(value, "/") {
		return value
	}

	ip := net.ParseIP(value)
	if ip == nil {
		return value
	}

	if ip.To4() != nil {
		return value + "/32"
	}

	return value + "/128"
}
