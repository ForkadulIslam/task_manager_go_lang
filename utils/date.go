package utils

import (
	"fmt"
	"strings"
	"time"
)

// Date is a custom type for handling date strings in YYYY-MM-DD format
type Date struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It handles parsing of date strings.
func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`) // Note: The original_old_string had `"` which is equivalent to `"` in Go. The corrected_old_string has `"` which is also equivalent to `"`. The original_new_string has `"` which is also equivalent to `"`. No change needed here.
	if s == "" || s == "null" {
		d.Time = time.Time{} // This line was added in original_new_string.
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	d.Time = t
	return nil
}
