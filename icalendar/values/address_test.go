package values

import (
	"fmt"
	. "github.com/taviti/check"
	"testing"
)

type AddressSuite struct{}

var _ = Suite(new(AddressSuite))

func TestAddress(t *testing.T) { TestingT(t) }

func (s *AddressSuite) TestEncodeName(c *C) {
	n := "Foo Bar"
	e := "foo@bar.com"
	uri := fmt.Sprintf("%s <%s>", n, e)
	o := NewOrganizerAddress(uri)
	c.Assert(o.EncodeICalName(), Equals, fmt.Sprintf("%s;CN=%s", o.Role(), n))
	uri = fmt.Sprintf("<%s>", e)
	a := NewAttendeeAddress(uri)
	c.Assert(a.EncodeICalName(), Equals, a.Role())
}

func (s *AddressSuite) TestEncodeValue(c *C) {
	e := "foo@bar.com"
	o := NewAttendeeAddress(e)
	c.Assert(o.EncodeICalValue(), Equals, fmt.Sprintf("MAILTO:%s", e))
}

func (s *AddressSuite) TestValidate(c *C) {
	e := "@foobar"
	o := NewOrganizerAddress(e)
	c.Assert(o.ValidateICalValue(), ErrorMatches, "address is invalid")
}
