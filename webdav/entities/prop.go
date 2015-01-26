package entities

import "encoding/xml"

// a property of a resource
type Prop struct {
	XMLName        xml.Name      `xml:"prop"`
	GetContentType string        `xml:"getcontenttype,omitempty"`
	DisplayName    string        `xml:"displayname,omitempty"`
	ResourceType   *ResourceType `xml:",omitempty"`
	CTag           string        `xml:"http://calendarserver.org/ns/ getctag,omitempty"`
}

// the type of a resource
type ResourceType struct {
	XMLName    xml.Name                `xml:"resourcetype"`
	Collection *ResourceTypeCollection `xml:",omitempty"`
	Calendar   *ResourceTypeCalendar   `xml:",omitempty"`
}

// A calendar resource type
type ResourceTypeCalendar struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:caldav calendar"`
}

// A collection resource type
type ResourceTypeCollection struct {
	XMLName xml.Name `xml:"collection"`
}
