package values

import (
	"fmt"
	"time"
)

const DateTimeFormatString = "20060102T150405"

// a representation of a date and time for iCalendar
type DateTime struct {
	t time.Time
}

// creates a new icalendar datetime representation
func NewDateTime(t time.Time) *DateTime {
	return &DateTime{t: t}
}

// encodes the datetime value for the iCalendar specification
func (d *DateTime) EncodeICalValue() string {
	val := d.t.Format(DateTimeFormatString)
	loc := d.t.Location()
	if loc == time.UTC {
		val = fmt.Sprintf("%sZ", val)
	}
	return val
}

// encodes the datetime value for the iCalendar specification
func (d *DateTime) String() string {
	return d.EncodeICalValue()
}
