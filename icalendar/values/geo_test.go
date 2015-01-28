package values

import (
	"fmt"
	. "github.com/taviti/check"
	"testing"
)

type GeoSuite struct{}

var _ = Suite(new(GeoSuite))

func TestGeo(t *testing.T) { TestingT(t) }

func (s *GeoSuite) TestLatLng(c *C) {
	geo := NewGeo(10, 20)
	c.Assert(geo.Lat(), Equals, float64(10))
	c.Assert(geo.Lng(), Equals, float64(20))
}

func (s *GeoSuite) TestEncode(c *C) {
	geo := NewGeo(10, -20)
	encoded, err := geo.EncodeICalValue()
	c.Assert(err, IsNil)
	expected := fmt.Sprintf("%f %f", geo.Lat(), geo.Lng())
	c.Assert(encoded, Equals, expected)
}

func (s *GeoSuite) TestValidate(c *C) {
	geo := NewGeo(-91, 0)
	c.Assert(geo.ValidateICalValue(), ErrorMatches, "latitude")
	geo = NewGeo(0, 181)
	c.Assert(geo.ValidateICalValue(), ErrorMatches, "longitude")
}
