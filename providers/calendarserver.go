package providers

import (
	"fmt"
	"github.com/taviti/caldav-go/caldav"
	"io"
	"net/http"
	"net/url"
)

// A provider for a CalendarServer.org server
type calendarServerProvider struct {
	*url.URL
}

func (c *calendarServerProvider) NewRequest(method, path string, body io.ReadCloser) (r *http.Request) {
	cURL := *c.URL
	cURL.Path = path
	r = new(http.Request)
	r.Method = method
	r.URL = &cURL
	r.Proto = "HTTP/1.1"
	r.ProtoMajor = 1
	r.ProtoMinor = 1
	r.Header = make(http.Header)
	r.Host = cURL.Host
	r.Body = nil
	return
}

func (c *calendarServerProvider) NewVerifyRequest() *http.Request {
	return c.NewRequest("OPTIONS", "/", nil)
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
