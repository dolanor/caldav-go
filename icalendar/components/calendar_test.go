package components

import (
	. "github.com/go-check/check"
	"github.com/taviti/caldav-go/icalendar"
	"testing"
)

type CalendarSuite struct{ calendar Calendar }

var _ = Suite(new(CalendarSuite))

func Test(t *testing.T) { TestingT(t) }

// tests the current server for CalDAV support
func (s *CalendarSuite) TestMarshal(c *C) {
	enc, err := icalendar.Marshal(s.calendar)
	c.Assert(err, Not(NotNil))
	c.Assert(enc, Equals, "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//taviti/caldav-go//NONSGML v1.0.0//EN\r\nEND:VCALENDAR")
}
