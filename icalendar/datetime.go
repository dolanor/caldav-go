package icalendar

import "time"

const DateTimeFormatString = "20060102T150405Z"

type DateTime struct {
	Time time.Time
}

func NewDateTime(time time.Time) *DateTime {
	return &DateTime{Time: time}
}

func (d *DateTime) String() string {
	return d.Time.Format(DateTimeFormatString)
}
