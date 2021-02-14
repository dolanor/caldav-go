package caldav

import (
	"github.com/dolanor/caldav-go/utils"
	"github.com/dolanor/caldav-go/webdav"
)

// a server that accepts CalDAV requests
type Server webdav.Server

// NewServer creates a reference to a CalDAV server.
// host is the url to access the server, stopping at the port:
// https://user:password@host:port/
func NewServer(host string) (*Server, error) {
	if s, err := webdav.NewServer(host); err != nil {
		return nil, utils.NewError(NewServer, "unable to create WebDAV server", host, err)
	} else {
		return (*Server)(s), nil
	}
}

// downcasts the server to the WebDAV interface
func (s *Server) WebDAV() *webdav.Server {
	return (*webdav.Server)(s)
}

// creates a new CalDAV request object
func (s *Server) NewRequest(method string, path string, icaldata ...interface{}) (*Request, error) {
	return NewRequest(method, s.WebDAV().Http().AbsUrlStr(path), icaldata...)
}
