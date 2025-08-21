package utils

import (
	"strings"
	"time"
)

// NullTime represents a time.Time that can be null/empty in JSON
type NullTime struct {
	time.Time
	Valid bool // True if Time is not null
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It handles parsing of time strings, including empty strings as null.
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		nt.Valid = false
		return nil
	}

	// Try to parse as YYYY-MM-DD (date only)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		// If parsing fails, return the error
		return err
	}

	nt.Time = t
	nt.Valid = true
	return nil
}