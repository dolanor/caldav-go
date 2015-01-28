package values

import (
	"fmt"
	"github.com/taviti/caldav-go/utils"
	"strconv"
	"strings"
)

// a representation of a geographical point for iCalendar
type Geo []float64

// creates a new icalendar geo representation
func NewGeo(lat, lng float64) *Geo {
	return &Geo{lat, lng}
}

// returns the latitude encoded into the geo point
func (g *Geo) Lat() float64 {
	return (*g)[0]
}

// returns the longitude encoded into the geo point
func (g *Geo) Lng() float64 {
	return (*g)[1]
}

// validates the geo value against the iCalendar specification
func (g *Geo) ValidateICalValue() error {

	args := []float64(*g)

	if len(args) != 2 {
		return utils.NewError(g.ValidateICalValue, "geo value must have length of 2", g, nil)
	}

	if g.Lat() < -90 || g.Lat() > 90 {
		return utils.NewError(g.ValidateICalValue, "geo latitude must be between -90 and 90 degrees", g, nil)
	}

	if g.Lng() < -180 || g.Lng() > 180 {
		return utils.NewError(g.ValidateICalValue, "geo longitude must be between -180 and 180 degrees", g, nil)
	}

	return nil

}

// encodes the geo value for the iCalendar specification
func (g *Geo) EncodeICalValue() (string, error) {
	return fmt.Sprintf("%f %f", g.Lat(), g.Lng()), nil
}

// decodes the geo value from the iCalendar specification
func (g *Geo) DecodeICalValue(value string) error {
	if latlng := strings.Split(value, " "); len(latlng) < 2 {
		return utils.NewError(g.DecodeICalValue, "geo value must have both a latitude and longitude component", g, nil)
	} else if lat, err := strconv.ParseFloat(latlng[0], 64); err != nil {
		return utils.NewError(g.DecodeICalValue, "unable to decode latitude component", g, err)
	} else if lng, err := strconv.ParseFloat(latlng[1], 64); err != nil {
		return utils.NewError(g.DecodeICalValue, "unable to decode latitude component", g, err)
	} else {
		g = (*Geo)(&[]float64{lat, lng})
		return nil
	}
}
