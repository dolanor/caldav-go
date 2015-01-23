package webdav

import (
	"github.com/taviti/caldav-go/http"
	"github.com/taviti/caldav-go/utils"
)

// a server that accepts WebDAV requests
type Server http.Server

// creates a reference to an WebDAV server
func NewServer(baseUrlStr string) (*Server, error) {
	if s, err := http.NewServer(baseUrlStr); err != nil {
		return nil, utils.NewError(NewServer, "unable to create WebDAV server", baseUrlStr, err)
	} else {
		return (*Server)(s), nil
	}
}

// downcasts the server to the local HTTP interface
func (s *Server) Http() *http.Server {
	return (*http.Server)(s)
}

// creates a new WebDAV request object
func (s *Server) NewRequest(method string, path string, xmldata ...interface{}) (*Request, error) {
	return NewRequest(method, s.Http().AbsUrlStr(path), xmldata...)
}
