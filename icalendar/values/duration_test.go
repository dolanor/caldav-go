package values

import (
	. "github.com/taviti/check"
	"testing"
	"time"
)

type DurationSuite struct{}

var _ = Suite(new(DateTimeSuite))

func TestDuration(t *testing.T) { TestingT(t) }

func (s *DurationSuite) TestIsPast(c *C) {
	d := NewDuration(-time.Second)
	c.Assert(d.IsPast(), Equals, true)
}

func (s *DurationSuite) TestDecompose(c *C) {
	d := NewDuration(5*7*24*time.Hour + 4*24*time.Hour + 3*time.Hour + 2*time.Minute + time.Second)
	weeks, days, hours, minutes, seconds := d.Decompose()
	c.Assert(weeks, Equals, 5)
	c.Assert(days, Equals, 4)
	c.Assert(hours, Equals, 3)
	c.Assert(minutes, Equals, 2)
	c.Assert(seconds, Equals, 1)
}

func (s *DurationSuite) TestEncode(c *C) {
	d := NewDuration(-5*7*24*time.Hour + 4*24*time.Hour + 3*time.Hour + 2*time.Minute + time.Second)
	encoded := d.EncodeICalValue()
	c.Assert(encoded, Equals, "-P5W4DT3H2M1S")
	d = NewDuration(7*24*time.Hour + 2*24*time.Hour)
	encoded = d.EncodeICalValue()
	c.Assert(encoded, Equals, "P1W2D")
	d = NewDuration(time.Hour + time.Minute + time.Second)
	encoded = d.EncodeICalValue()
	c.Assert(encoded, Equals, "PT1H1M1S")
}
