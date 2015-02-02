package entities

import "encoding/xml"

// a WebDAV error
type Error struct {
	XMLName     xml.Name `xml:"DAV: error"`
	Description string   `xml:"error-description,omitempty"`
	Message     string   `xml:"message,omitempty"`
}

func (e *Error) Error() string {
	if e.Description != "" {
		return e.Description
	} else {
		return e.Message
	}
}
