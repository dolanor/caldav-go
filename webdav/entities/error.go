package entities

import "encoding/xml"

// a WebDAV error
type Error struct {
	XMLName     xml.Name `xml:"DAV: error"`
	Description string   `xml:"error-description,omitempty"`
}

func (e *Error) Error() string {
	return e.Description
}
