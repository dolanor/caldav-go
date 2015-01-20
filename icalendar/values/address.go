package values

import (
	"fmt"
	"github.com/taviti/caldav-go/icalendar"
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
	role, uri string
}

type AttendeeAddress Address
type OrganizerAddress Address
type RelationAddress Address

// creates a new icalendar attendee representation
// Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
func NewAttendeeAddress(uri string) *AttendeeAddress {
	return &AttendeeAddress{uri: uri, role: "ATTENDEE"}
}

// creates a new icalendar attendee representation
// Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
func NewRelationAddress(uri string) *RelationAddress {
	return &RelationAddress{uri: uri, role: "RELATED-TO"}
}

// creates a new icalendar organizer representation
// Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
func NewOrganizerAddress(uri string) *OrganizerAddress {
	return &OrganizerAddress{uri: uri, role: "ORGANIZER"}
}

// creates an RFC 5322 compliant address representation for the owner
func (a *Address) MailAddress() (*mail.Address, error) {
	return mail.ParseAddress(a.uri)
}

// validates the address value against the iCalendar specification
func (a *Address) ValidateICalValue() error {

	if _, err := a.MailAddress(); err != nil {
		return icalendar.NewError(a.ValidateICalValue, "mailing address is invalid", a, err)
	}

	return nil

}

// encodes the address value for the iCalendar specification
func (a *Address) EncodeICalValue() string {
	m, _ := a.MailAddress()
	return fmt.Sprintf("MAILTO:%s", m.Address)
}

// encodes the organizer name for the iCalendar specification
func (a *Address) EncodeICalName() string {
	if m, err := a.MailAddress(); err == nil && m.Name != "" {
		return fmt.Sprintf("%s;CN=%s", a.role, m.Name)
	} else {
		return a.role
	}
}
