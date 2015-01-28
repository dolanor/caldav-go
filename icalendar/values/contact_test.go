package values

import (
	"github.com/taviti/caldav-go/icalendar"
	. "github.com/taviti/check"
	"net/mail"
	"testing"
)

type ContactSuite struct{}

var _ = Suite(new(ContactSuite))

func TestContact(t *testing.T) { TestingT(t) }

func (s *ContactSuite) TestMarshalWithName(c *C) {
	addy := mail.Address{Name: "Foo Bar", Address: "foo@bar.com"}
	o := NewOrganizerContact(addy)
	enc, err := icalendar.Marshal(o)
	c.Assert(err, IsNil)
	c.Assert(enc, Equals, "ORGANIZER;CN=Foo Bar:MAILTO:foo@bar.com")
}

func (s *ContactSuite) TestMarshalWithoutName(c *C) {
	addy := mail.Address{Address: "foo@bar.com"}
	o := NewAttendeeContact(addy)
	enc, err := icalendar.Marshal(o)
	c.Assert(err, IsNil)
	c.Assert(enc, Equals, "ATTENDEE:MAILTO:foo@bar.com")
}

func (s *ContactSuite) TestItentity(c *C) {

	addy := mail.Address{Name: "Foo Bar", Address: "foo@bar.com"}

	before := NewOrganizerContact(addy)
	encoded, err := icalendar.Marshal(before)
	c.Assert(err, IsNil)

	after := new(OrganizerContact)
	err = icalendar.Unmarshal(encoded, after)
	c.Assert(err, IsNil)

	c.Assert(after, DeepEquals, before)

}
