package icalendar

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	Newline = "\r\n"
)

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

func marshalProperty(name, value string) string {
	return fmt.Sprintf("%s:%s", strings.ToUpper(name), value)
}

// converts an icalendar component into its string representation
// adapted from: https://github.com/icalendar/icalendar
func Marshal(target interface{}) (string, error) {

	var out []string
	v := reflect.ValueOf(target)

	if !v.IsValid() || isEmptyValue(v) {
		return "", nil
	}

	// drill into interfaces and pointers.
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vkind := v.Kind()
	vtype := v.Type()

	if vkind != reflect.Struct {
		return "", nil
	}

	// iterate over all fields
	n := vtype.NumField()
	for i := 0; i < n; i++ {

		// ignore non-exported or explicitly ignored fields
		fv := v.Field(i)
		fs := vtype.Field(i)
		ftag := vtype.Field(i).Tag.Get("ical")
		if fs.PkgPath != "" || ftag == "-" {
			continue
		}

		// parse the field tag
		var name, value string
		var omitempty = false
		if ftag != "" {
			tags := strings.Split(ftag, ",")
			name = tags[0]
			if len(tags) > 1 {
				if tags[1] == "omitempty" {
					omitempty = true
				} else {
					value = tags[1]
				}
			}
		}

		// make sure we have a name
		if name == "" {
			name = fs.Name
		}

		// drill into interfaces and pointers.
		for fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}

		// check for collections of components
		fkind := fv.Kind()
		if fkind == reflect.Slice || fkind == reflect.Array {
			for i, n := 0, fv.Len(); i < n; i++ {
				if encoded, err := Marshal(fv.Index(i)); err != nil {
					return "", fmt.Errorf("unable to encode field %s[%d], %s", name, i, err)
				} else if encoded != "" {
					out = append(out, encoded)
				}
			}
			continue
		}

		// check for nested components
		fi := fv.Interface()
		if fkind == reflect.Struct {
			if encoded, err := Marshal(fi); err != nil {
				return "", fmt.Errorf("unable to encode field %s, %s", name, err)
			} else if encoded != "" {
				out = append(out, encoded)
			}
			continue
		}

		// check to override default
		fvalue := fmt.Sprintf("%s", fi)
		if fvalue != "" {
			value = fvalue
		}

		// omit empty values if requested
		if value == "" && omitempty {
			continue
		}

		// encode in the property
		out = append(out, marshalProperty(name, value))

	}

	name := "V" + strings.ToUpper(vtype.Name())
	out = append([]string{marshalProperty("BEGIN", name)}, out...)
	out = append(out, marshalProperty("END", name))

	return strings.Join(out, Newline), nil

}
