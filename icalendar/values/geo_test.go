package values

import (
	"fmt"
	"github.com/dolanor/caldav-go/icalendar"
	. "gopkg.in/check.v1"
	"testing"
)

type GeoSuite struct{}

type geoTestObj struct {
	*Geo
}

var _ = Suite(new(GeoSuite))

func TestGeo(t *testing.T) { TestingT(t) }

func (s *GeoSuite) TestLatLng(c *C) {
	geo := NewGeo(10, 20)
	c.Assert(geo.Lat(), Equals, float64(10))
	c.Assert(geo.Lng(), Equals, float64(20))
}

func (s *GeoSuite) TestEncode(c *C) {
	gto := new(geoTestObj)
	gto.Geo = NewGeo(10, -20)
	encoded, err := icalendar.Marshal(gto)
	c.Assert(err, IsNil)
	expected := fmt.Sprintf("BEGIN:VGEOTESTOBJ\r\nGEO:%f %f\r\nEND:VGEOTESTOBJ", gto.Geo.Lat(), gto.Geo.Lng())
	c.Assert(encoded, Equals, expected)
}

func (s *GeoSuite) TestValidate(c *C) {
	geo := NewGeo(-91, 0)
	c.Assert(geo.ValidateICalValue(), ErrorMatches, "(?s).*latitude.*")
	geo = NewGeo(0, 181)
	c.Assert(geo.ValidateICalValue(), ErrorMatches, "(?s).*longitude.*")
}

func (s *GeoSuite) TestIdentity(c *C) {

	before := &geoTestObj{Geo: NewGeo(10, 20)}
	encoded, err := icalendar.Marshal(before)
	c.Assert(err, IsNil)

	after := new(geoTestObj)
	if err := icalendar.Unmarshal(encoded, after); err != nil {
		c.Fatal(err.Error())
	}

	c.Assert(after, DeepEquals, before)

}
