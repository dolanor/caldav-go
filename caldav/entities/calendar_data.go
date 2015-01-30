package entities

import (
	"encoding/xml"
	"github.com/taviti/caldav-go/caldav/values"
	"github.com/taviti/caldav-go/icalendar"
	"github.com/taviti/caldav-go/icalendar/components"
	"github.com/taviti/caldav-go/utils"
	"strings"
)

// a CalDAV calendar data object
type CalendarData struct {
	XMLName             xml.Name             `xml:"urn:ietf:params:xml:ns:caldav calendar-data"`
	Component           *Component           `xml:",omitempty"`
	RecurrenceSetLimit  *RecurrenceSetLimit  `xml:",omitempty"`
	ExpandRecurrenceSet *ExpandRecurrenceSet `xml:",omitempty"`
	Content             string               `xml:",chardata"`
}

func (c *CalendarData) CalendarComponent() (*components.Calendar, error) {
	cal := new(components.Calendar)
	if content := strings.TrimSpace(c.Content); content == "" {
		return nil, utils.NewError(c.CalendarComponent, "no calendar data to decode", c, nil)
	} else if err := icalendar.Unmarshal(content, cal); err != nil {
		return nil, utils.NewError(c.CalendarComponent, "decoding calendar data failed", c, err)
	} else {
		return cal, nil
	}
}

// an iCalendar specifier for returned calendar data
type Component struct {
	XMLName    xml.Name        `xml:"urn:ietf:params:xml:ns:caldav comp"`
	Properties []*PropertyName `xml:",omitempty"`
	Components []*Component    `xml:",omitempty"`
}

// used to restrict recurring event data to a particular time range
type RecurrenceSetLimit struct {
	XMLName   xml.Name         `xml:"urn:ietf:params:xml:ns:caldav limit-recurrence-set"`
	StartTime *values.DateTime `xml:"start,attr"`
	EndTime   *values.DateTime `xml:"end,attr"`
}

// used to expand recurring events into individual calendar event data
type ExpandRecurrenceSet struct {
	XMLName   xml.Name         `xml:"urn:ietf:params:xml:ns:caldav expand"`
	StartTime *values.DateTime `xml:"start,attr"`
	EndTime   *values.DateTime `xml:"end,attr"`
}
