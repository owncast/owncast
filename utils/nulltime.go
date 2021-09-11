package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Value interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func (nt NullTime) UnmarshalJSON(data []byte) error {
	dateString := string(data)
	if dateString == "null" {
		return nil
	}

	dateStringWithoutQuotes := dateString[1 : len(dateString)-1]
	parsedDateTime, err := time.Parse(time.RFC3339, dateStringWithoutQuotes)
	if err != nil {
		return err
	}

	nt.Time = parsedDateTime // nolint
	return nil
}
