package utils

import (
	"fmt"
	"strings"
	"time"
)

// NullableDate is a custom type for handling nullable date strings in YYYY-MM-DD format
type NullableDate struct {
	time.Time
	Present bool
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It handles parsing of date strings.
func (d *NullableDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		d.Present = false
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	d.Time = t
	d.Present = true
	return nil
}