package caldav

import (
	"github.com/taviti/caldav-go/webdav"
	"github.com/taviti/caldav-go/webdav/entities"
)

// a client for communicating with a CalDAV server
type Client interface {

	// verifies that the CalDAV capabilities exist on a path
	// you may supply an optional list of extra capabilities to test the server for
	// returns an error if the provider does not support CalDAV or any provided extra capabilities
	Validate(path string, capabilities ...string) error

	// executes a propfind request on a particular URI, returns a multistatus response
	// the path is the the location of the entity to run propfind on
	// the propfind entity is used to describe the propfind request
	// the depth tells the propfind request how deep to search down the hierarchy
	Propfind(path string, depth webdav.Depth, pf *entities.Propfind) (*entities.Multistatus, error)

	// returns the provider interface used to configure the client
	Provider() Provider
}
