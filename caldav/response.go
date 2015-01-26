package caldav

import (
	"github.com/taviti/caldav-go/webdav"
	"io/ioutil"
	"log"
)

// a WebDAV response object
type Response webdav.Response

// downcasts the response to the WebDAV interface
func (r *Response) WebDAV() *webdav.Response {
	return (*webdav.Response)(r)
}

// decodes a CalDAV iCalendar response into the provided interface
func (r *Response) Decode(into interface{}) error {
	if body := r.Body; body == nil {
		return nil
	} else {
		out, err := ioutil.ReadAll(body)
		log.Printf("\nOUT: %s", out)
		return err
	}
}

// creates a new WebDAV response object
func NewResponse(response *webdav.Response) *Response {
	return (*Response)(response)
}
