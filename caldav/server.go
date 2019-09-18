package caldav

import (
	"github.com/dolanor/caldav-go/utils"
	"github.com/dolanor/caldav-go/webdav"
)

// a server that accepts CalDAV requests
type Server webdav.Server

// creates a reference to a CalDAV server
func NewServer(baseUrlStr string) (*Server, error) {
	if s, err := webdav.NewServer(baseUrlStr); err != nil {
		return nil, utils.NewError(NewServer, "unable to create WebDAV server", baseUrlStr, err)
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
