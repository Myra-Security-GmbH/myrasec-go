package types

import (
	"encoding/json"
	"strings"
	"time"
)

//
// DateTime ...
//
type DateTime struct {
	time.Time
}

//
// DateTimeNow ...
//
func DateTimeNow() *DateTime {
	ret := &DateTime{}
	ret.Time = time.Now()

	return ret
}

//
// MarshalJSON ...
//
func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.Format(time.RFC3339))
}

//
// UnmarshalJSON ...
//
func (dt *DateTime) UnmarshalJSON(b []byte) error {
	date := strings.Trim(string(b), "\"")

	t, err := time.Parse(time.RFC3339, date)
	dt.Time = t

	return err
}
