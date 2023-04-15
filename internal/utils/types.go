package utils

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"
)

// NullString is a wrapper around sql.NullString
type NullString sql.NullString

// Scan implements the Scanner interface for NullString
func (ns NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		ns = NullString{s.String, false}
	} else {
		ns = NullString{s.String, true}
	}

	return nil
}

// MarshalJSON method is called by json.Marshal,
// whenever it is of type NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

type SqlTime []byte
func (s SqlTime) Time() (time.Time, error) {
    return time.Parse("15:04:05",string(s))
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