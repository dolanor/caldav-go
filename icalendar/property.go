package icalendar

import (
	"fmt"
	"github.com/taviti/caldav-go/utils"
	"log"
	"reflect"
	"strings"
)

var _ = log.Print

const (
	Newline = "\r\n"
)

var propNameSanitizer = strings.NewReplacer(
	"_", "-",
	":", "\\:",
)

var propValueSanitizer = strings.NewReplacer(
	"\"", "'",
	"\\", "\\\\",
	"\n", "\\n",
)

var propNameDesanitizer = strings.NewReplacer(
	"-", "_",
	"\\:", ":",
)

var propValueDesanitizer = strings.NewReplacer(
	"'", "\"",
	"\\\\", "\\",
	"\\n", "\n",
)

type canValidateValue interface {
	ValidateICalValue() error
}

type property struct {
	Name, Value, DefaultValue string
	Params                    map[string]string
	OmitEmpty, Required       bool
}

func (p *property) HasNameAndValue() bool {
	return p.Name != "" && p.Value != ""
}

func (p *property) Merge(override *property) {
	if override.Name != "" {
		p.Name = override.Name
	}
	if override.Value != "" {
		p.Value = override.Value
	}
	if override.Params != nil {
		p.Params = override.Params
	}
}

func propertyFromStructField(fs reflect.StructField) (p *property) {

	ftag := fs.Tag.Get("ical")
	if fs.PkgPath != "" || ftag == "-" {
		return
	}

	p = new(property)

	// parse the field tag
	if ftag != "" {
		tags := strings.Split(ftag, ",")
		p.Name = tags[0]
		if len(tags) > 1 {
			if tags[1] == "omitempty" {
				p.OmitEmpty = true
			} else if tags[1] == "required" {
				p.Required = true
			} else {
				p.DefaultValue = tags[1]
			}
		}
	}

	// make sure we have a name
	if p.Name == "" {
		p.Name = fs.Name
	}

	p.Name = strings.ToUpper(p.Name)

	return

}

func marshalProperty(p *property) string {
	name := strings.ToUpper(propNameSanitizer.Replace(p.Name))
	value := propValueSanitizer.Replace(p.Value)
	keys := []string{name}
	for name, value := range p.Params {
		name = strings.ToUpper(propNameSanitizer.Replace(name))
		value = propValueSanitizer.Replace(value)
		keys = append(keys, fmt.Sprintf("%s=%s", name, value))
	}
	name = strings.Join(keys, ";")
	return fmt.Sprintf("%s:%s", name, value)
}

func propertyFromInterface(target interface{}) (p *property, err error) {

	var ierr error
	if va, ok := target.(canValidateValue); ok {
		if ierr = va.ValidateICalValue(); ierr != nil {
			err = utils.NewError(propertyFromInterface, "interface failed validation", target, ierr)
			return
		}
	}

	p = new(property)

	if enc, ok := target.(canEncodeName); ok {
		if p.Name, ierr = enc.EncodeICalName(); ierr != nil {
			err = utils.NewError(propertyFromInterface, "interface failed name encoding", target, ierr)
			return
		}
	}

	if enc, ok := target.(canEncodeParams); ok {
		if p.Params, ierr = enc.EncodeICalParams(); ierr != nil {
			err = utils.NewError(propertyFromInterface, "interface failed params encoding", target, ierr)
			return
		}
	}

	if enc, ok := target.(canEncodeValue); ok {
		if p.Value, ierr = enc.EncodeICalValue(); ierr != nil {
			err = utils.NewError(propertyFromInterface, "interface failed value encoding", target, ierr)
			return
		}
	}

	return

}

func unmarshalProperty(line string) *property {
	nvp := strings.SplitN(line, ":", 2)
	prop := new(property)
	if len(nvp) > 1 {
		prop.Value = strings.TrimSpace(nvp[1])
	}
	npp := strings.Split(nvp[0], ";")
	if len(npp) > 1 {
		prop.Params = make(map[string]string, 0)
		for i := 1; i < len(npp); i++ {
			var key, value string
			kvp := strings.Split(npp[i], "=")
			key = strings.TrimSpace(kvp[0])
			key = propNameDesanitizer.Replace(key)
			if len(kvp) > 1 {
				value = strings.TrimSpace(kvp[1])
				value = propValueDesanitizer.Replace(value)
			}
			prop.Params[key] = value
		}
	}
	prop.Name = strings.TrimSpace(npp[0])
	prop.Name = propNameDesanitizer.Replace(prop.Name)
	prop.Value = propValueDesanitizer.Replace(prop.Value)
	return prop
}
