package entities

import (
	"encoding/xml"
	"github.com/taviti/caldav-go/webdav/entities"
)

// a CalDAV Property resource
type Prop struct {
	XMLName        xml.Name               `xml:"DAV: prop"`
	GetContentType string                 `xml:"getcontenttype,omitempty"`
	DisplayName    string                 `xml:"displayname,omitempty"`
	CalendarData   *CalendarData          `xml:",omitempty"`
	ResourceType   *entities.ResourceType `xml:",omitempty"`
	CTag           string                 `xml:"http://calendarserver.org/ns/ getctag,omitempty"`
	ETag           string                 `xml:"http://calendarserver.org/ns/ getetag,omitempty"`
}

// used to restrict properties returned in calendar data
type PropertyName struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:caldav prop"`
	Name    string   `xml:"name,attr"`
}
