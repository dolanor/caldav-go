package components

import (
	"github.com/taviti/caldav-go/icalendar"
	. "github.com/taviti/check"
	"testing"
	"time"
)

type TimezoneSuite struct{}

var _ = Suite(new(TimezoneSuite))

func TestTimezone(t *testing.T) { TestingT(t) }

// tests the current server for CalDAV support
func (s *TimezoneSuite) TestMarshal(c *C) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	c.Assert(err, IsNil)
	tz := NewDynamicTimeZone(loc)
	enc, err := icalendar.Marshal(tz)
	c.Assert(err, IsNil)
	c.Assert(enc, Equals, "BEGIN:VTIMEZONE\r\nTZID:America/Los_Angeles\r\nX-LIC-LOCATION:America/Los_Angeles\r\n"+
		"TZURL;VALUE=URI:http://tzurl.org/zoneinfo/America/Los_Angeles\r\nEND:VTIMEZONE")
}
