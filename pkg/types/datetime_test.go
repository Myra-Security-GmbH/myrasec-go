package types

import (
	"strings"
	"testing"
	"time"
)

func TestDateTimeNow(t *testing.T) {
	dtn := DateTimeNow()

	now := time.Now()

	if dtn.Format("2006-01-02") != now.Format("2006-01-02") {
		t.Errorf("Expected to get a date like [%s] but got [%s]", now.Format("2006-01-02"), dtn.Format("2006-01-02"))
	}
}

func TestMarshalJSON(t *testing.T) {
	dtn := DateTimeNow()
	res, err := dtn.MarshalJSON()
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	now := time.Now()

	if !strings.Contains(string(res), now.Format("2006-01-02")) {
		t.Errorf("Expected to find [%s] in the date [%s]", now.Format("2006-01-02"), string(res))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	date := "2022-06-08T16:29:16+0200"

	dt := &DateTime{}
	err := dt.UnmarshalJSON([]byte(date))
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if dt.Format("2006-01-02") != "2022-06-08" {
		t.Errorf("Expected to get a date like [%s] but got [%s]", "2022-06-08", dt.Format("2006-01-02"))
	}

}

func TestParseDate(t *testing.T) {
	now := time.Now()
	dt, err := ParseDate(now.Format(time.RFC3339))
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if dt.Format("2006-01-02") != now.Format("2006-01-02") {
		t.Errorf("Expected to get a date like [%s] but got [%s]", now.Format("2006-01-02"), dt.Format("2006-01-02"))
	}
}
