package icalendar

import (
	"fmt"
	"reflect"
	"strings"
)

const Newline = "\r\n"

// used for icalendar components that do not have children
var Leaf []Component

// an icalendar component
type Component interface {
	Name() string
	Children() []Component
	Validate() error
}

func isEmptyValue(v reflect.Value) bool {
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

func marshalFields(target Component) (string, error) {

	if err := target.Validate(); err != nil {
		return "", fmt.Errorf("unable to encode %s component, %s", target.Name(), err)
	}

	var out []string
	v := reflect.ValueOf(target)

	if !v.IsValid() || isEmptyValue(v) {
		return "", nil
	}

	// Drill into interfaces and pointers.
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vkind := v.Kind()
	vtype := v.Type()

	if vkind == reflect.Struct {

		// iterate over all fields
		n := vtype.NumField()
		for i := 0; i < n; i++ {

			// ignore non-exported or explicitly ignored fields
			f := vtype.Field(i)
			ftag := f.Tag.Get("ical")
			if f.PkgPath != "" || ftag == "-" {
				continue
			}

			// check for components
			fi := f.Interface()
			if fc, ok := fi.(Component); ok {
				if encoded, err := Marshal(fc); err != nil {
					return "", fmt.Errorf("unable to encode field %s, %s", fc.Name(), err)
				} else if encoded != "" {
					out = append(out, encoded)
				}
				continue
			}

			// parse the field tag
			var name string
			var omitempty = false
			if ftag != "" {
				name = strings.Split(ftag, ",")[0]
				omitempty = strings.Contains(ftag, "omitempty")
			}

			// make sure we have a name
			if name == "" {
				name = f.Name
			}

			// omit empty values if requested
			var value = fmt.Sprintf("%s", fi)
			if value == "" && omitempty {
				continue
			}

			// encode in the property
			if encoded, err := marshalProperty(name, value); err != nil {
				return "", err
			} else {
				out = append(out, encoded)
			}

		}
	}

	return strings.Join(out, Newline), nil

}

func marshalChildren(target Component) (string, error) {
	var out []string
	for i, child := range target.Children() {
		if encoded, err := Marshal(child); err != nil {
			return "", fmt.Errorf("error marshaling child %s at index %d, %s", target.Name(), i, err)
		} else {
			out = append(out, encoded)
		}
	}
	return strings.Join(out, Newline), nil
}

func marshalProperty(name, value string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("unable to encode empty property name")
	} else {
		return fmt.Sprintf("%s:%s", strings.ToUpper(name), value), nil
	}
}

// converts an icalendar component into its string representation
// adapted from: https://github.com/icalendar/icalendar
func Marshal(target Component) (string, error) {
	var out []string
	if encoded, err := marshalProperty("BEGIN", target.Name()); err != nil {
		return "", err
	} else {
		out = append(out, encoded)
	}
	if encoded, err := marshalFields(target); err != nil {
		return "", fmt.Errorf("unable to marshal fields for %s, %s", target.Name(), err)
	} else if encoded != "" {
		out = append(out, encoded)
	}
	if encoded, err := marshalChildren(target); err != nil {
		return "", fmt.Errorf("unable to marshal children for %s, %s", target.Name(), err)
	} else if encoded != "" {
		out = append(out, encoded)
	}
	if encoded, err := marshalProperty("END", target.Name()); err != nil {
		return "", err
	} else {
		out = append(out, encoded)
	}
	return strings.Join(out, Newline), nil
}
