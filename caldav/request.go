package caldav

import (
	"github.com/taviti/caldav-go/utils"
	"github.com/taviti/caldav-go/webdav"
)

// an CalDAV request object
type Request webdav.Request

// downcasts the request to the WebDAV interface
func (r *Request) WebDAV() *webdav.Request {
	return &webdav.Request(r)
}

// creates a new CalDAV request object
func NewRequest(method string, urlstr string, xmldata ...interface{}) (*Request, error) {
	if r, err := webdav.NewRequest(method, urlstr, xmldata...); err != nil {
		return nil, utils.NewError(NewRequest, "unable to create request", urlstr, err)
	} else {
		return &Request(r), nil
	}
}
