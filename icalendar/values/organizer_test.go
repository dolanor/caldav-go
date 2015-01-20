package values

import (
	"fmt"
	. "github.com/taviti/check"
	"testing"
)

type OrganizerSuite struct{}

var _ = Suite(new(OrganizerSuite))

func TestOrganizer(t *testing.T) { TestingT(t) }

func (s *OrganizerSuite) TestEncodeName(c *C) {
	n := "Foo Bar"
	e := fmt.Sprintf("%s <foo@bar.com>", n)
	o := NewOrganizer(e)
	c.Assert(o.EncodeICalName(), Equals, fmt.Sprintf("ORGANIZER;CN=%s", n))
	e = "<foo@bar.com>"
	o = NewOrganizer(e)
	c.Assert(o.EncodeICalName(), Equals, "ORGANIZER")
}

func (s *OrganizerSuite) TestEncodeValue(c *C) {
	e := "foo@bar.com"
	o := NewOrganizer(e)
	c.Assert(o.EncodeICalValue(), Equals, fmt.Sprintf("MAILTO:%s", e))
}

func (s *OrganizerSuite) TestValidate(c *C) {
	e := "@foobar"
	o := NewOrganizer(e)
	c.Assert(o.ValidateICalValue(), ErrorMatches, "address is invalid")
}
