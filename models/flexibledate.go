package models

import (
	"database/sql"
	"errors"
	"time"
)

type FlexibleDate struct {
	time.Time
}

func (self *FlexibleDate) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	s = s[1 : len(s)-1]

	result, err := FlexibleDateParse(s)
	if err != nil {
		return err
	}

	self.Time = result

	return
}

// FlexibleDateParse is a convinience function to parse a date that could be
// a string, a time.Time, or a sql.NullTime.
func FlexibleDateParse(date interface{}) (time.Time, error) {
	// If it's within a sql.NullTime wrapper, return the time from that.
	nulltime, ok := date.(sql.NullTime)
	if ok {
		return nulltime.Time, nil
	}

	// Parse as string
	datestring, ok := date.(string)
	if ok {
		t, err := time.Parse(time.RFC3339Nano, datestring)
		if err == nil {
			return t, nil
		}

		t, err = time.Parse("2006-01-02T15:04:05.999999999Z0700", datestring)
		if err == nil {
			return t, nil
		}

		t, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", datestring)
		if err == nil {
			return t, nil
		}
	}

	dateobject, ok := date.(time.Time)
	if ok {
		return dateobject, nil
	}

	return time.Time{}, errors.New("unable to parse date")
}
