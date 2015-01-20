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
type Organizer struct {
	address string
}

// creates a new icalendar event organizer representation
// Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
func NewOrganizer(address string) *Organizer {
	return &Organizer{address: address}
}

// creates an RFC 5322 compliant address representation for the Organizer
func (o *Organizer) Address() (*mail.Address, error) {
	return mail.ParseAddress(o.address)
}

// validates the geo value against the iCalendar specification
func (o *Organizer) ValidateICalValue() error {

	if _, err := o.Address(); err != nil {
		return icalendar.NewError(o.ValidateICalValue, "organizer address is invalid", o, err)
	}

	return nil

}

// encodes the geo value for the iCalendar specification
func (o *Organizer) EncodeICalValue() string {
	a, _ := o.Address()
	return fmt.Sprintf("MAILTO:%s", a.Address)

}

// encodes the geo value for the iCalendar specification
func (o *Organizer) EncodeICalName() string {
	if a, _ := o.Address(); a.Name != "" {
		return fmt.Sprintf("ORGANIZER;CN=%s", a.Name)
	} else {
		return "ORGANIZER"
	}
}
