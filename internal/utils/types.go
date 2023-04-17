package utils

import (
	"database/sql/driver"
	"time"
)

type SqlTime []byte

func (s SqlTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(s))
}

// A wrapper for time.Time
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
