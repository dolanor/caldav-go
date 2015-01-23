package webdav

import (
	"encoding/xml"
	"github.com/taviti/caldav-go/http"
	"github.com/taviti/caldav-go/utils"
	"strings"
)

// a WebDAV response object
type Response http.Response

// downcasts the response to the local HTTP interface
func (r *Response) Http() *http.Response {
	return (*http.Response)(r)
}

// returns a list of WebDAV features found in the response
func (r *Response) Features() (features []string) {
	if dav := r.Header.Get("DAV"); dav != "" {
		features = strings.Split(dav, " ")
	}
	return
}

// decodes a WebDAV XML response into the provided interface
func (r *Response) Decode(into interface{}) error {
	if body := r.Body; body == nil {
		return nil
	} else if decoder := xml.NewDecoder(body); decoder == nil {
		return nil
	} else if err := decoder.Decode(into); err != nil {
		return utils.NewError(r.Decode, "unable to decode response body", r, err)
	} else {
		return nil
	}
}

// creates a new WebDAV response object
func NewResponse(response *http.Response) *Response {
	return (*Response)(response)
}
