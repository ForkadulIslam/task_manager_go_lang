package utils

import (
	"fmt"
	"strings"
	"time"
)

// RequiredDate is a custom type for handling required date strings in YYYY-MM-DD format
type RequiredDate struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It handles parsing of date strings.
func (d *RequiredDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return fmt.Errorf("date field is required")
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	d.Time = t
	return nil
}