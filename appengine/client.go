package appengine

import (
	"appengine"
	"appengine/urlfetch"
	"github.com/taviti/caldav-go/caldav"
)

// creates a new client for communicating with a CalDAV server
// uses the HTTP client from appengine/urlfetch
func NewClient(c appengine.Context, provider caldav.Provider) caldav.Client {
	return caldav.NewClient(provider, urlfetch.Client(c))
}
