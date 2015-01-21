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
type Address struct {
	role    string
	address mail.Address
}

type AttendeeAddress Address
type OrganizerAddress Address

// creates a new icalendar attendee representation
func NewAttendeeAddress(address mail.Address) *AttendeeAddress {
	return &AttendeeAddress{address: address, role: "ATTENDEE"}
}

// creates a new icalendar organizer representation
func NewOrganizerAddress(address mail.Address) *OrganizerAddress {
	return &OrganizerAddress{address: address, role: "ORGANIZER"}
}

// encodes the address value for the iCalendar specification
func (a *Address) EncodeICalValue() string {
	return fmt.Sprintf("MAILTO:%s", a.address.Address)
}

// encodes the organizer name for the iCalendar specification
func (a *Address) EncodeICalName() string {
	if a.address.Name != "" {
		return fmt.Sprintf("%s;CN=%s", a.role, a.address.Name)
	} else {
		return a.role
	}
}

// encodes the address value for the iCalendar specification
func (a *OrganizerAddress) EncodeICalValue() string {
	return (*Address)(a).EncodeICalValue()
}

// encodes the organizer name for the iCalendar specification
func (a *OrganizerAddress) EncodeICalName() string {
	return (*Address)(a).EncodeICalName()
}

// encodes the address value for the iCalendar specification
func (a *AttendeeAddress) EncodeICalValue() string {
	return (*Address)(a).EncodeICalValue()
}

// encodes the organizer name for the iCalendar specification
func (a *AttendeeAddress) EncodeICalName() string {
	return (*Address)(a).EncodeICalName()
}
