package icalendar

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	Newline = "\r\n"
)

type validatable interface {
	ValidateICalValue() error
}

type encodable interface {
	EncodeICalValue() (string, error)
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

var propNameSanitizer = strings.NewReplacer(
	"_", "-",
	":", "",
)

var propValueSanitizer = strings.NewReplacer(
	"\"", "'",
	"\\", "\\\\",
	";", "\\;",
	",", "\\,",
	"\n", "\\n",
)

func marshalProperty(name, value string) string {
	name = propNameSanitizer.Replace(name)
	value = propValueSanitizer.Replace(value)
	return fmt.Sprintf("%s:%s", strings.ToUpper(name), value)
}

// converts an iCalendar component into its string representation
func Marshal(target interface{}) (string, error) {

	if v, ok := target.(validatable); ok {
		if err := v.ValidateICalValue(); err != nil {
			return "", NewError(Marshal, "interface failed validation", target, err)
		}
	}

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

	// encode arrays early
	if vkind == reflect.Slice || vkind == reflect.Array {

		for i, n := 0, v.Len(); i < n; i++ {
			if encoded, err := Marshal(v.Index(i)); err != nil {
				msg := fmt.Sprintf("unable to encode component at index %d", i)
				return "", NewError(Marshal, msg, target, err)
			} else if encoded != "" {
				out = append(out, encoded)
			}
		}

		return strings.Join(out, Newline), nil

		// fail early on non-structs
	} else if vkind != reflect.Struct {
		return "", NewError(Marshal, "only structs and enumerations are encodable", target, nil)
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
		var required = false
		if ftag != "" {
			tags := strings.Split(ftag, ",")
			name = tags[0]
			if len(tags) > 1 {
				if tags[1] == "omitempty" {
					omitempty = true
				} else if tags[1] == "required" {
					required = true
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

		// now check to see if we have anything to encode
		if fv.IsValid() && !isEmptyValue(fv) {

			// check for collections of components
			fkind := fv.Kind()
			if fkind == reflect.Slice || fkind == reflect.Array {
				for i, n := 0, fv.Len(); i < n; i++ {
					if encoded, err := Marshal(fv.Index(i)); err != nil {
						msg := fmt.Sprintf("unable to encode field %s at index %d", fs.Name, i)
						return "", NewError(Marshal, msg, target, err)
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
					msg := fmt.Sprintf("unable to encode field %s", fs.Name)
					return "", NewError(Marshal, msg, target, err)
				} else if encoded != "" {
					out = append(out, encoded)
				}
				continue
			}

			// check to override default
			var encoded string
			if encoder, ok := fi.(encodable); ok {
				var err error
				if encoded, err = encoder.EncodeICalValue(); err != nil {
					msg := fmt.Sprintf("unable to encode field %s", fs.Name)
					return "", NewError(Marshal, msg, target, err)
				}
			} else {
				encoded = fmt.Sprintf("%s", fi)
			}

			if encoded != "" {
				value = encoded
			}

		}

		// check empty values for required or empty
		if value == "" {
			if required {
				msg := fmt.Sprintf("missing value for required field %s", fs.Name)
				return "", NewError(Marshal, msg, target, nil)
			} else if omitempty {
				continue
			}
		}

		// encode in the property
		out = append(out, marshalProperty(name, value))

	}

	name := "V" + strings.ToUpper(vtype.Name())
	out = append([]string{marshalProperty("begin", name)}, out...)
	out = append(out, marshalProperty("end", name))

	return strings.Join(out, Newline), nil

}
