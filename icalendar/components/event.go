package components

import (
	"github.com/taviti/caldav-go/icalendar"
	"github.com/taviti/caldav-go/icalendar/values"
)

type Event struct {
	UID         string           `ical:,required`
	DateStamp   *values.DateTime `ical:dtstamp,required`
	DateStart   *values.DateTime `ical:dtstart,required`
	DateEnd     *values.DateTime `ical:dtend,omitempty`
	Duration    *values.Duration `ical:,omitempty`
	Class       values.Class     `ical:,omitempty`
	Created     *values.DateTime `ical:,omitempty`
	Description string           `ical:,omitempty`

	//optional_single_property :geo, Icalendar::Values::Float
	//optional_single_property :last_modified, Icalendar::Values::DateTime
	//optional_single_property :location
	//optional_single_property :organizer, Icalendar::Values::CalAddress
	//optional_single_property :priority, Icalendar::Values::Integer
	//optional_single_property :sequence, Icalendar::Values::Integer
	//optional_single_property :status
	//optional_single_property :summary
	//optional_single_property :transp
	//optional_single_property :url, Icalendar::Values::Uri
	//optional_single_property :recurrence_id, Icalendar::Values::DateTime
	//
	//optional_property :rrule, Icalendar::Values::Recur, true
	//optional_property :attach, Icalendar::Values::Uri
	//optional_property :attendee, Icalendar::Values::CalAddress
	//optional_property :categories
	//optional_property :comment
	//optional_property :contact
	//optional_property :exdate, Icalendar::Values::DateTime
	//optional_property :request_status
	//optional_property :related_to
	//optional_property :resources
	//optional_property :rdate, Icalendar::Values::DateTime
	//
	//component :alarm, false
}

func (e *Event) ValidateICalValue() error {

	if e.DateEnd != nil && e.Duration != nil {
		return icalendar.NewError(e.ValidateICalValue, "DateEnd and Duration are mutually exclusive fields", e, nil)
	}

	return nil

}
