package utils

import (
	"database/sql/driver"
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
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	d.Time = t
	return nil
}

// Value implements the driver.Valuer interface.
// This method is called by GORM when writing to the database.
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// Scan implements the sql.Scanner interface.
// This method is called by GORM when reading from the database.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan type %T into Date", value)
	}
	d.Time = t
	return nil
}