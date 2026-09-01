package myrasec

import (
	"crypto/sha256"
	"fmt"
	"slices"
)

// IntInSlice checks if the haystack []int slice contains the passed needle int
func intInSlice(needle int, haystack []int) bool {
	return slices.Contains(haystack, needle)
}

// BoolPtr returns a pointer to the passed bool value.
// This is a convenience helper for setting *bool fields in structs like Settings.
func BoolPtr(v bool) *bool {
	return &v
}

// BuildSHA256 builds the SHA256 for the passed string
func BuildSHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
