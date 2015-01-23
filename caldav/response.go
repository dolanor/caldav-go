package webdav

import (
	"github.com/taviti/caldav-go/webdav"
)

// a WebDAV response object
type Response webdav.Response

// downcasts the response to the WebDAV interface
func (r *Response) Http() *webdav.Response {
	return (*webdav.Response)(r)
}

// creates a new WebDAV response object
func NewResponse(response *webdav.Response) *Response {
	return (*Response)(response)
}
