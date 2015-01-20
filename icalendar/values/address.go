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
type Address interface {
	Role() string
	URI() string
}

type AttendeeAddress struct {
	uri string
}

type OrganizerAddress struct {
	uri string
}

// creates an RFC 5322 compliant address representation for the owner
func ParseMailAddress(address Address) (*mail.Address, error) {
	return mail.ParseAddress(address.URI())
}

func ValidateMailAddress(address Address) error {
	if _, err := ParseMailAddress(address); err != nil {
		return icalendar.NewError(ValidateMailAddress, "mailing address is invalid", address, err)
	}
	return nil
}

func MailtoLink(address Address) string {
	m, _ := ParseMailAddress(address)
	return fmt.Sprintf("MAILTO:%s", m.Address)
}

func CanonicalName(address Address) string {
	if m, err := ParseMailAddress(address); err == nil && m.Name != "" {
		return fmt.Sprintf("%s;CN=%s", address.Role(), m.Name)
	} else {
		return address.Role()
	}
}

// creates a new icalendar attendee representation
// Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
func NewAttendeeAddress(uri string) *AttendeeAddress {
	return &AttendeeAddress{uri: uri}
}

// creates a new icalendar owner representation
// Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
func NewOrganizerAddress(uri string) *OrganizerAddress {
	return &OrganizerAddress{uri: uri}
}

// validates the address value against the iCalendar specification
func (a *AttendeeAddress) ValidateICalValue() error {
	return ValidateMailAddress(a)
}

// validates the address value against the iCalendar specification
func (a *OrganizerAddress) ValidateICalValue() error {
	return ValidateMailAddress(a)
}

// encodes the address value for the iCalendar specification
func (a *OrganizerAddress) EncodeICalValue() string {
	return MailtoLink(a)
}

// encodes the address value for the iCalendar specification
func (a *AttendeeAddress) EncodeICalValue() string {
	return MailtoLink(a)
}

// encodes the attendee name for the iCalendar specification
func (a *AttendeeAddress) EncodeICalName() string {
	return CanonicalName(a)
}

// encodes the organizer name for the iCalendar specification
func (a *OrganizerAddress) EncodeICalName() string {
	return CanonicalName(a)
}

// encodes the attendee role for the iCalendar specification
func (a *AttendeeAddress) Role() string {
	return "ATTENDEE"
}

// encodes the organizer role for the iCalendar specification
func (a *OrganizerAddress) Role() string {
	return "ORGANIZER"
}

// encodes the attendee address URI for the iCalendar specification
func (a *AttendeeAddress) URI() string {
	return a.uri
}

// encodes the organizer address URI for the iCalendar specification
func (a *OrganizerAddress) URI() string {
	return a.uri
}
