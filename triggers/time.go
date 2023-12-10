package triggers

import (
	"strings"
	"time"
)

const (
	// iso8601 format.
	iso8601 = "2006-01-02T15:04:05"
	// iso8601TZ format with timezone.
	iso8601TZ = "2006-01-02T15:04:05-07:00"
)

// TimeISO8601 is a wrapper around time.Time that supports JSON marshalling/unmarshalling
// with format ISO8601 and ISO8601 with a provided timezone.
type TimeISO8601 struct {
	time.Time
	tz bool
}

// UnmarshalJSON unmarshals a JSON time string into a TimeISO8601.
// It supports with and without timezone.
func (t *TimeISO8601) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || len(s) == 0 {
		return nil
	}
	parsed, err := time.Parse(iso8601, s)
	if err != nil {
		parsed, err = time.Parse(iso8601TZ, s)
		if err != nil {
			return err
		}
		t.tz = true
	}
	t.Time = parsed
	return nil
}

// MarshalJSON marshals a TimeISO8601 into a JSON time string.
func (t TimeISO8601) MarshalJSON() ([]byte, error) {
	var formatted string
	if t.tz {
		formatted = t.Format(iso8601TZ)
	} else {
		formatted = t.Format(iso8601)
	}
	return []byte(`"` + formatted + `"`), nil
}
