package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Bool is a boolean flag that decodes both the JSON boolean and the string
// encoding the API uses for some flags.
//
// The API serializes a few boolean attributes as strings, for example the
// user's agent flag, which arrives as "" (false) or "1" (true) instead of a
// JSON boolean. Bool accepts a JSON boolean, the strings "", "0", "false",
// "no" and "off" (false) and "1", "true", "yes" and "on" (true, case
// insensitive), the numbers 0 and 1, and null (false). Any other value is a
// decoding error, so an unexpected change of the API's encoding surfaces
// instead of being read as false. Bool marshals as a plain JSON boolean and
// behaves like bool in every other respect; convert with bool(v) where a
// plain bool is required.
type Bool bool

// UnmarshalJSON decodes a JSON boolean, a string or a number into the Bool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case nil:
		*b = false
	case bool:
		*b = Bool(v)
	case float64:
		switch v {
		case 0:
			*b = false
		case 1:
			*b = true
		default:
			return fmt.Errorf("cannot decode number %s as a boolean flag", data)
		}
	case string:
		switch strings.ToLower(v) {
		case "", "0", "false", "no", "off":
			*b = false
		case "1", "true", "yes", "on":
			*b = true
		default:
			return fmt.Errorf("cannot decode string %q as a boolean flag", v)
		}
	default:
		return fmt.Errorf("cannot decode %s as a boolean flag", data)
	}

	return nil
}
