package values

import (
	"fmt"
	"net/mail"
)

// Specifies the organizer of a group scheduled calendar entity. The property is specified within the "VFREEBUSY"
// calendar component to specify the calendar user requesting the free or busy time. When publishing a "VFREEBUSY"
// calendar component, the property is used to specify the calendar that the published busy time came from.
//
// The property has the property parameters CN, for specifying the common or display name associated with the
// "Organizer", DIR, for specifying a pointer to the directory information associated with the "Organizer",
// SENT-BY, for specifying another calendar user that is acting on behalf of the "Organizer". The non-standard
// parameters may also be specified on this property. If the LANGUAGE property parameter is specified, the identified
// language applies to the CN parameter value.
type Contact struct {
	address mail.Address
}

type AttendeeContact Contact
type OrganizerContact Contact

// creates a new icalendar attendee representation
func NewAttendeeContact(address mail.Address) *AttendeeContact {
	return &AttendeeContact{address: address}
}

// creates a new icalendar organizer representation
func NewOrganizerContact(address mail.Address) *OrganizerContact {
	return &OrganizerContact{address: address}
}

// encodes the contact value for the iCalendar specification
func (c *Contact) EncodeICalValue() string {
	return fmt.Sprintf("MAILTO:%s", c.address.Address)
}

// encodes the contact params for the iCalendar specification
func (c *Contact) EncodeICalParams() (params map[string]string) {
	if c.address.Name != "" {
		params = map[string]string{"CN": c.address.Name}
	}
	return
}

// encodes the contact value for the iCalendar specification
func (o *OrganizerContact) EncodeICalValue() string {
	return (*Contact)(o).EncodeICalValue()
}

// encodes the contact params for the iCalendar specification
func (o *OrganizerContact) EncodeICalParams() map[string]string {
	return (*Contact)(o).EncodeICalParams()
}

// encodes the organizer name for the iCalendar specification
func (o *OrganizerContact) EncodeICalName() string {
	return "ORGANIZER"
}

// encodes the contact value for the iCalendar specification
func (a *AttendeeContact) EncodeICalValue() string {
	return (*Contact)(a).EncodeICalValue()
}

// encodes the contact params for the iCalendar specification
func (a *AttendeeContact) EncodeICalParams() map[string]string {
	return (*Contact)(a).EncodeICalParams()
}

// encodes the organizer name for the iCalendar specification
func (a *AttendeeContact) EncodeICalName() string {
	return "ATTENDEE"
}
