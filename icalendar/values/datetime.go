package values

import (
	"fmt"
	"time"
)

const DateTimeFormatString = "20060102T150405"

type DateTime struct {
	t time.Time
}

func NewDateTime(t time.Time) *DateTime {
	return &DateTime{t: t}
}

func (d *DateTime) EncodeICalValue() (string, error) {
	val := d.t.Format(DateTimeFormatString)
	loc := d.t.Location()
	if loc == time.UTC {
		val = fmt.Sprintf("%sZ", val)
	}
	return val, nil
}
