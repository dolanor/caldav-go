package values

import (
	"fmt"
	. "github.com/taviti/check"
	"net/mail"
	"testing"
)

type AddressSuite struct{}

var _ = Suite(new(AddressSuite))

func TestAddress(t *testing.T) { TestingT(t) }

func (s *AddressSuite) TestEncodeName(c *C) {
	addy := mail.Address{Name: "Foo Bar", Address: "foo@bar.com"}
	o := NewOrganizerAddress(addy)
	c.Assert(o.EncodeICalName(), Equals, fmt.Sprintf("%s;CN=%s", o.role, addy.Name))
	addy.Name = ""
	a := NewAttendeeAddress(addy)
	c.Assert(a.EncodeICalName(), Equals, a.role)
}

func (s *AddressSuite) TestEncodeValue(c *C) {
	addy := mail.Address{Address: "foo@bar.com"}
	o := NewAttendeeAddress(addy)
	c.Assert(o.EncodeICalValue(), Equals, fmt.Sprintf("MAILTO:%s", addy.Address))
}
