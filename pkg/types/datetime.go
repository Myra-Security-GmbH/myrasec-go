package types

import (
	"encoding/json"
	"strings"
	"time"
)

// DateTime ...
type DateTime struct {
	time.Time
}

// DateTimeNow ...
func DateTimeNow() *DateTime {
	ret := &DateTime{}
	ret.Time = time.Now()

	return ret
}

// MarshalJSON ...
func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.Format("2006-01-02T15:04:05Z0700"))
}

// UnmarshalJSON ...
func (dt *DateTime) UnmarshalJSON(b []byte) error {
	date := strings.Trim(string(b), "\"")

	t, err := time.Parse(
		"2006-01-02T15:04:05Z0700",
		date,
	)

	dt.Time = t

	return err
}

// ToUnixDate ...
func (dt *DateTime) ToUnixDate() string {
	format := "_2. Jan 2006 "

	if dt.Year() == time.Now().Year() {
		format = "_2. Jan 15:04"
	}

	return dt.Format(format)
}

// ParseDate transforms the passed date string (time.RFC3339) to a DateTime struct
func ParseDate(date string) (*DateTime, error) {
	if len(date) <= 0 {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}

	return &DateTime{
		Time: parsed,
	}, nil
}
