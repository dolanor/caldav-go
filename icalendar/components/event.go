package components

import (
	"github.com/taviti/caldav-go/icalendar"
	"github.com/taviti/caldav-go/icalendar/values"
	"net/url"
)

type Event struct {

	// defines the persistent, globally unique identifier for the calendar component.
	UID string `ical:",required"`

	// indicates the date/time that the instance of the iCalendar object was created.
	DateStamp *values.DateTime `ical:"dtstamp,required"`

	// specifies when the calendar component begins.
	DateStart *values.DateTime `ical:"dtstart,required"`

	// specifies the date and time that a calendar component ends.
	DateEnd *values.DateTime `ical:"dtend,omitempty"`

	// specifies a positive duration of time.
	Duration *values.Duration `ical:",omitempty"`

	// defines the access classification for a calendar component.
	AccessClassification values.EventAccessClassification `ical:class,omitempty`

	// specifies the date and time that the calendar information was created by the calendar user agent in the
	// calendar store.
	// Note: This is analogous to the creation date and time for a file in the file system.
	Created *values.DateTime `ical:",omitempty"`

	// provides a more complete description of the calendar component, than that provided by the Summary property.
	Description string `ical:",omitempty"`

	// specifies information related to the global position for the activity specified by a calendar component.
	Geo *values.Geo `ical:",omitempty"`

	// specifies the date and time that the information associated with the calendar component was last revised in the
	// calendar store.
	// Note: This is analogous to the modification date and time for a file in the file system.
	LastModified *values.DateTime `ical:"last-modified,omitempty"`

	// defines the intended venue for the activity defined by a calendar component.
	Location string `ical:",omitempty"`

	// defines the organizer for a calendar component.
	Organizer *values.OrganizerAddress `ical:",omitempty"`

	// defines the relative priority for a calendar component.
	Priority int `ical:",omitempty"`

	// defines the revision sequence number of the calendar component within a sequence of revisions.
	Sequence int `ical:",omitempty"`

	// efines the overall status or confirmation for the calendar component.
	Status values.EventStatus `ical:",omitempty"`

	// defines a short summary or subject for the calendar component.
	Summary string `ical:",omitempty"`

	// defines whether an event is transparent or not to busy time searches.
	values.TimeTransparency `ical:"transp,omitempty"`

	// defines a Uniform Resource Locator (URL) associated with the iCalendar object.
	Url *url.URL `ical:",omitempty"`

	// used in conjunction with the "UID" and "SEQUENCE" property to identify a specific instance of a recurring
	// event calendar component. The property value is the effective value of the DateStart property of the
	// recurrence instance.
	RecurrenceId *values.DateTime `ical:"recurrence_id,omitempty"`

	// defines a rule or repeating pattern for recurring events, to-dos, or time zone definitions.
	*values.RecurrenceRule `ical:",omitempty"`

	// property provides the capability to associate a document object with a calendar component.
	Attachment *url.URL `ical:"attach,omitempty"`

	// defines an "Attendee" within a calendar component.
	Attendees []*values.AttendeeAddress `,omitempty`
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
