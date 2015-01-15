package providers

import (
	"fmt"
	"github.com/taviti/caldav-go/caldav"
	"net/url"
)

// A provider for a CalendarServer.org server
type calendarServerProvider struct {
	*url.URL
}

func (c *calendarServerProvider) AbsURL(path string) *url.URL {
	safe := *c.URL
	safe.Path = path
	return &safe
}

func (c *calendarServerProvider) CurrentUsername() string {
	return c.URL.User.Username()
}

func (c *calendarServerProvider) UserCalendarPath(username string) string {
	return fmt.Sprintf("/calendars/users/%s/calendar/", username)
}

// Creates a provider for a CalendarServer.org server
// See http://calendarserver.org/
func NewCalendarServerProvider(hosturl string) (caldav.Provider, error) {
	var err error
	p := new(calendarServerProvider)
	if p.URL, err = url.Parse(hosturl); err != nil {
		err = fmt.Errorf("unable to parse calendar server host URL '%s'. %s", hosturl, err)
	}
	return p, err
}
