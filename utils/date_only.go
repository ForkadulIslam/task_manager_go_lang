package utils

import (
	"strings"
	"time"
)

// DateOnly is a custom type for handling date strings in YYYY-MM-DD format
type DateOnly struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It handles parsing of date strings.
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