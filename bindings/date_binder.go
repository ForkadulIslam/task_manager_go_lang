package bindings

import (
	"encoding/json"
	"strings"
	"time"
)

// DateOnly is a custom type for handling date strings in YYYY-MM-DD format
type DateOnly struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It handles parsing of date strings, including empty strings as zero time.
func (do *DateOnly) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		do.Time = time.Time{} // Set to zero value for empty/null
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	do.Time = t
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// It handles marshaling of DateOnly to JSON, returning empty string for zero time.
func (do DateOnly) MarshalJSON() ([]byte, error) {
	if do.Time.IsZero() {
		return []byte(`""`), nil // Return empty string for zero time
	}
	return []byte(`"` + do.Time.Format("2006-01-02") + `"`), nil
}
