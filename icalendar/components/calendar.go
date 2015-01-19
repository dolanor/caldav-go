package components

import (
	"github.com/taviti/caldav-go/icalendar"
	"github.com/taviti/caldav-go/icalendar/values"
)

type Calendar struct {
	Version string        `ical:",2.0"`
	ProdId  string        `ical:",-//taviti/caldav-go//NONSGML v1.0.0//EN"`
	Method  values.Method `ical:",omitempty"`
	*Event  `ical:",omitempty"`
}

func (c *Calendar) ValidateICalValue() error {

	e := c.Event

	if e == nil {
		return nil
	}

	if err := e.ValidateICalValue(); err != nil {
		return icalendar.NewError(c.ValidateICalValue, "event failed validation", c, err)
	}

	if e.DateStart == nil && c.Method == "" {
		return icalendar.NewError(c.ValidateICalValue, "no value for method and no start date defined on event", c, nil)
	}

	return nil

}
