package values

import (
	"github.com/dolanor/caldav-go/icalendar"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

type DurationSuite struct{}

var _ = Suite(new(DurationSuite))

func TestDuration(t *testing.T) { TestingT(t) }

func (s *DurationSuite) TestIsPast(c *C) {
	d := NewDuration(-time.Second)
	c.Assert(d.IsPast(), Equals, true)
}

func (s *DurationSuite) TestDecompose(c *C) {
	d := NewDuration(5*7*24*time.Hour + 4*24*time.Hour + 3*time.Hour + 2*time.Minute + time.Second)
	weeks, days, hours, minutes, seconds := d.Decompose()
	c.Assert(weeks, Equals, int64(5))
	c.Assert(days, Equals, int64(4))
	c.Assert(hours, Equals, int64(3))
	c.Assert(minutes, Equals, int64(2))
	c.Assert(seconds, Equals, int64(1))
}

func (s *DurationSuite) TestEncode(c *C) {
	d := NewDuration(-(5*7*24*time.Hour + 4*24*time.Hour + 3*time.Hour + 2*time.Minute + time.Second))
	encoded, err := d.EncodeICalValue()
	c.Assert(err, IsNil)
	c.Assert(encoded, Equals, "-P5W4DT3H2M1S")
	d = NewDuration(7*24*time.Hour + 2*24*time.Hour)
	encoded, err = d.EncodeICalValue()
	c.Assert(err, IsNil)
	c.Assert(encoded, Equals, "P1W2D")
	d = NewDuration(time.Hour + time.Minute + time.Second)
	encoded, err = d.EncodeICalValue()
	c.Assert(err, IsNil)
	c.Assert(encoded, Equals, "PT1H1M1S")
}

func (s *DurationSuite) TestIdentity(c *C) {

	type test struct {
		*Duration
	}

	before := &test{Duration: NewDuration(12*time.Hour + 10*time.Minute)}
	encoded, err := icalendar.Marshal(before)
	c.Assert(err, IsNil)

	after := new(test)
	err = icalendar.Unmarshal(encoded, after)
	c.Assert(err, IsNil)

	c.Assert(after, DeepEquals, before)

}
