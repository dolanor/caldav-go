package caldav

import (
	"net/url"
)

// A mechanism for defining, authenticating and connecting to a CalDav server
type Provider interface {

	// creates a URL object for a CalDAV service provider
	AbsURL(path string) *url.URL

	// returns the user name of the currently authenticated user
	CurrentUsername() string

	// returns the calendar path for a given user account
	UserCalendarPath(username string) string
}
