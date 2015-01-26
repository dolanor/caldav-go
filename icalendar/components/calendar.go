package components

import (
	"github.com/taviti/caldav-go/icalendar/values"
	"github.com/taviti/caldav-go/utils"
	"time"
)

type Calendar struct {

	// specifies the identifier corresponding to the highest version number or the minimum and maximum
	// range of the iCalendar specification that is required in order to interpret the iCalendar object.
	Version string `ical:",2.0"`

	// specifies the identifier for the product that created the iCalendar object
	ProductId string `ical:"prodid,-//taviti/caldav-go//NONSGML v1.0.0//EN"`

	// specifies the text value that uniquely identifies the "VTIMEZONE" calendar component.
	TimeZoneId string `ical:"tzid,omitempty"`

	// defines the iCalendar object method associated with the calendar object.
	Method values.Method `ical:",omitempty"`

	// defines the calendar scale used for the calendar information specified in the iCalendar object.
	values.CalScale `ical:",omitempty"`

	*TimeZone `ical:",omitempty"`
	*Event    `ical:",omitempty"`
}

func (c *Calendar) UseTimeZone(location *time.Location) *TimeZone {
	c.TimeZone = NewDynamicTimeZone(location)
	c.TimeZoneId = c.TimeZone.Id
	return c.TimeZone
}

func (c *Calendar) UsingTimeZone() bool {
	return len(c.TimeZoneId) > 0
}

func (c *Calendar) UsingGlobalTimeZone() bool {
	return c.UsingTimeZone() && c.TimeZoneId[0] == '/'
}

func (c *Calendar) ValidateICalValue() error {

	e := c.Event

	if e == nil {
		return nil
	}

	if err := e.ValidateICalValue(); err != nil {
		return utils.NewError(c.ValidateICalValue, "event failed validation", c, err)
	}

	if e.DateStart == nil && c.Method == "" {
		return utils.NewError(c.ValidateICalValue, "no value for method and no start date defined on event", c, nil)
	}

	if c.UsingTimeZone() && !c.UsingGlobalTimeZone() {
		if c.TimeZone == nil || c.TimeZone.Id != c.TimeZoneId {
			return utils.NewError(c.ValidateICalValue, "calendar timezone ID does not match timezone component", c, nil)
		}
	}

	return nil

}

func NewCalendar(event *Event) *Calendar {
	cal := new(Calendar)
	cal.Event = event
	return cal
}
