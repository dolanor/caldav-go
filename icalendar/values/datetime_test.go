package values

import (
	"fmt"
	. "github.com/taviti/check"
	"testing"
)

type DateTimeSuite struct{ time DateTime }

var _ = Suite(new(DateTimeSuite))

func TestDateTime(t *testing.T) { TestingT(t) }

func (s *DateTimeSuite) TestEncode(c *C) {
	actual := s.time.EncodeICalValue()
	expected := s.time.t.Format(DateTimeFormatString)
	c.Assert(actual, Equals, fmt.Sprintf("%sZ", expected))
}
