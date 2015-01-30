package values

import (
	"encoding/xml"
	"errors"
	"github.com/taviti/caldav-go/icalendar/values"
	"time"
)

// a representation of a date and time for iCalendar
type DateTime struct {
	name string
	t    time.Time
}

// creates a new caldav datetime representation, must be in UTC
func NewDateTime(name string, t time.Time) (*DateTime, error) {
	if t.Location() != time.UTC {
		return nil, errors.New("CalDAV datetime must be in UTC")
	} else {
		return &DateTime{name: name, t: t.Truncate(time.Second)}, nil
	}
}

// encodes the datetime value for the iCalendar specification
func (d *DateTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	layout := values.UTCDateTimeFormatString
	value := d.t.Format(layout)
	attr := xml.Attr{Name: name, Value: value}
	return attr, nil
}
