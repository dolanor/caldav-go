package values

import (
	"fmt"
	. "github.com/taviti/check"
	"testing"
)

type DateTimeSuite struct{ time DateTime }

var _ = Suite(new(DateTimeSuite))

func Test(t *testing.T) { TestingT(t) }

// tests the current server for CalDAV support
func (s *DateTimeSuite) TestEncode(c *C) {
	enc, err := s.time.EncodeICalValue()
	c.Assert(err, Not(NotNil))
	encoded := s.time.t.Format(DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf("%sZ", encoded))
}
