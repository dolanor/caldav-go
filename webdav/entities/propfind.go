package entities

import "encoding/xml"

// a request to find properties on an an entity or collection
type Propfind struct {
	XMLName xml.Name `xml:"DAV: propfind"`
	AllProp *AllProp `xml:",omitempty"`
	Props   []*Prop  `xml:"prop,omitempty"`
}

// a propfind property representing all properties
type AllProp struct {
	XMLName xml.Name `xml:"allprop"`
}

// a convenience method for searching all properties
func NewAllPropsFind() *Propfind {
	return &Propfind{AllProp: new(AllProp)}
}
