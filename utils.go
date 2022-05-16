package myrasec

import (
	"crypto/sha256"
	"fmt"
)

//
// IntInSlice checks if the haystack []int slice contains the passed needle int
//
func intInSlice(needle int, haystack []int) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}
	}
	return false
}

//
// BuildSHA256 builds the SHA256 for the passed string
//
func BuildSHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
