package caldav

import (
	"bytes"
	"github.com/taviti/caldav-go/http"
	"github.com/taviti/caldav-go/icalendar"
	"github.com/taviti/caldav-go/utils"
	"github.com/taviti/caldav-go/webdav"
	"io"
	"io/ioutil"
	"strings"
)

// an CalDAV request object
type Request webdav.Request

// downcasts the request to the WebDAV interface
func (r *Request) WebDAV() *webdav.Request {
	return (*webdav.Request)(r)
}

// creates a new CalDAV request object
func NewRequest(method string, urlstr string, icaldata ...interface{}) (*Request, error) {
	if buffer, err := icalToReadCloser(icaldata); err != nil {
		return nil, utils.NewError(NewRequest, "unable to encode icalendar data", icaldata, err)
	} else if r, err := http.NewRequest(method, urlstr, buffer); err != nil {
		return nil, utils.NewError(NewRequest, "unable to create request", urlstr, err)
	} else {
		if buffer != nil {
			// set the content type to XML if we have a body
			r.Native().Header.Set("Content-Type", "text/calendar; charset=UTF-8")
		}
		return (*Request)(r), nil
	}
}

func icalToReadCloser(icaldata ...interface{}) (io.ReadCloser, error) {
	var buffer []string
	for _, icaldatum := range icaldata {
		if encoded, err := icalendar.Marshal(icaldatum); err != nil {
			return nil, utils.NewError(icalToReadCloser, "unable to encode as xml", icaldatum, err)
		} else {
			buffer = append(buffer, string(encoded))
		}
	}
	if len(buffer) > 0 {
		var encoded = strings.Join(buffer, "\n")
		return ioutil.NopCloser(bytes.NewBuffer([]byte(encoded))), nil
	} else {
		return nil, nil
	}
}
