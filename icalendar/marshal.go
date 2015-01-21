package icalendar

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	Newline = "\r\n"
)

type canEncodeTag interface {
	EncodeICalTag() string
}

type canValidateValue interface {
	ValidateICalValue() error
}

type canEncodeValue interface {
	EncodeICalValue() string
}

type canEncodeName interface {
	EncodeICalName() string
}

type canEncodeParams interface {
	EncodeICalParams() map[string]string
}

type encoder func(reflect.Value) (string, error)

type property struct {
	Name, Value, DefaultValue string
	Params                    map[string]string
	OmitEmpty, Required       bool
}

var propNameSanitizer = strings.NewReplacer(
	"_", "-",
	":", "\\:",
)

var propValueSanitizer = strings.NewReplacer(
	"\"", "'",
	"\\", "\\\\",
	";", "\\;",
	"\n", "\\n",
)

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

func isInvalidOrEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func tagAndJoinValue(v reflect.Value, in []string) string {
	var tag string
	if encoder, ok := v.Interface().(canEncodeTag); ok {
		tag = encoder.EncodeICalTag()
	}
	if tag == "" {
		tag = fmt.Sprintf("v%s", v.Type().Name())
	}
	tag = strings.ToUpper(tag)
	var out []string
	out = append(out, marshalProperty(&property{Name: "begin", Value: tag}))
	out = append(out, in...)
	out = append(out, marshalProperty(&property{Name: "end", Value: tag}))
	return strings.Join(out, Newline)
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

func dereferencePointerValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func propertyFromInterface(target interface{}) (p *property, err error) {

	if va, ok := target.(canValidateValue); ok {
		if ierr := va.ValidateICalValue(); ierr != nil {
			err = NewError(propertyFromInterface, "interface failed validation", target, ierr)
			return
		}
	}

	p = new(property)

	if enc, ok := target.(canEncodeName); ok {
		p.Name = enc.EncodeICalName()
	}

	if enc, ok := target.(canEncodeParams); ok {
		p.Params = enc.EncodeICalParams()
	}

	if enc, ok := target.(canEncodeValue); ok {
		p.Value = enc.EncodeICalValue()
	}

	return

}

func marshalCollection(v reflect.Value) (string, error) {

	var out []string

	for i, n := 0, v.Len(); i < n; i++ {
		vi := v.Index(i).Interface()
		if encoded, err := Marshal(vi); err != nil {
			msg := fmt.Sprintf("unable to encode interface at index %d", i)
			return "", NewError(marshalCollection, msg, vi, err)
		} else if encoded != "" {
			out = append(out, encoded)
		}
	}

	return strings.Join(out, Newline), nil

}

func marshalStruct(v reflect.Value) (string, error) {

	var out []string

	// iterate over all fields
	vtype := v.Type()
	n := vtype.NumField()

	for i := 0; i < n; i++ {

		// keep a reference to the field value and definition
		fv := v.Field(i)
		fs := vtype.Field(i)
		fi := fv.Interface()

		// use the field definition to extract out property defaults
		p := propertyFromStructField(fs)
		if p == nil {
			continue // skip explicitly ignored fields and private members
		}

		// some fields are not properties, but actually nested objects.
		// detect those early using the property and object encoder...
		if _, ok := fi.(canEncodeValue); !ok && !isInvalidOrEmptyValue(fv) {
			if encoded, err := encode(fv, objectEncoder); err != nil {
				msg := fmt.Sprintf("unable to encode field %s", fs.Name)
				return "", NewError(marshalStruct, msg, v.Interface(), err)
			} else if encoded != "" {
				// encoding worked! no need to process as a property
				out = append(out, encoded)
				continue
			}
		}

		// now check to see if the field value overrides the defaults...
		if !isInvalidOrEmptyValue(fv) {
			// first, check the field value interface for overrides...
			if overrides, err := propertyFromInterface(fi); err != nil {
				msg := fmt.Sprintf("field %s failed validation", fs.Name)
				return "", NewError(marshalStruct, msg, v.Interface(), err)
			} else if p.Merge(overrides); p.Value == "" {
				// then, if we couldn't find an override from the interface,
				// try the simple string encoder...
				if p.Value, err = stringEncoder(fv); err != nil {
					msg := fmt.Sprintf("unable to encode field %s", fs.Name)
					return "", NewError(marshalStruct, msg, v.Interface(), err)
				}
			}
		}

		// make sure we have a value by this point
		if !p.HasNameAndValue() {
			if p.OmitEmpty {
				continue
			} else if p.DefaultValue != "" {
				p.Value = p.DefaultValue
			} else if p.Required {
				msg := fmt.Sprintf("missing value for required field %s", fs.Name)
				return "", NewError(Marshal, msg, v.Interface(), nil)
			}
		}

		// encode in the property
		out = append(out, marshalProperty(p))

	}

	// wrap the fields in the enclosing struct tags
	return tagAndJoinValue(v, out), nil

}

func objectEncoder(v reflect.Value) (string, error) {

	// decompose the value into its interface parts
	v = dereferencePointerValue(v)

	// encode the value based off of its type
	switch v.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return marshalCollection(v)
	case reflect.Struct:
		return marshalStruct(v)
	}

	return "", nil

}

func stringEncoder(v reflect.Value) (string, error) {
	return fmt.Sprintf("%v", v.Interface()), nil
}

func propertyEncoder(v reflect.Value) (string, error) {

	vi := v.Interface()
	if p, err := propertyFromInterface(vi); err != nil {

		// return early if interface fails its own validation
		return "", err

	} else if p.HasNameAndValue() {

		// if an interface encodes its own name and value, it's a property
		return marshalProperty(p), nil

	}

	return "", nil

}

func encode(v reflect.Value, encoders ...encoder) (string, error) {

	for _, encode := range encoders {
		if encoded, err := encode(v); err != nil {
			return "", err
		} else if encoded != "" {
			return encoded, nil
		}
	}

	return "", nil

}

// converts an iCalendar component into its string representation
func Marshal(target interface{}) (string, error) {

	// don't do anything with invalid interfaces
	v := reflect.ValueOf(target)
	if isInvalidOrEmptyValue(v) {
		return "", NewError(Marshal, "unable to marshal empty or invalid values", target, nil)
	}

	if encoded, err := encode(v, propertyEncoder, objectEncoder); err != nil {
		return "", err
	} else if encoded == "" {
		return "", NewError(Marshal, "unable to encode interface, all methods exhausted", v.Interface(), nil)
	} else {
		return encoded, nil
	}

}
