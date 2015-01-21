package values

import (
	"net/url"
)

// a representation of duration for iCalendar
type Url struct {
	u url.URL
}

// encodes the URL into iCalendar format
func (u *Url) EncodeICalValue() string {
	return u.u.String()
}

// creates a new iCalendar duration representation
func NewUrl(u url.URL) *Url {
	return &Url{u: u}
}
