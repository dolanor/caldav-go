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

type encodableValue interface {
	EncodeICalValue() string
}

type encodableName interface {
	EncodeICalName() string
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
	":", "\\:",
	"\n", "\\n",
)

func marshalProperty(name, value string) string {
	name = propNameSanitizer.Replace(name)
	value = propValueSanitizer.Replace(value)
	if !strings.Contains(name, ";") {
		name = strings.ToUpper(name)
	}
	return fmt.Sprintf("%s:%s", name, value)
}

// converts an iCalendar component into its string representation
func Marshal(target interface{}) (string, error) {

	// don't do anything with invalid interfaces
	v := reflect.ValueOf(target)
	if !v.IsValid() || isEmptyValue(v) {
		return "", NewError(Marshal, "unable to encode empty or invalid interfaces", target, nil)
	}

	// check for self-validating interfaces
	if va, ok := target.(validatable); ok {
		if err := va.ValidateICalValue(); err != nil {
			return "", NewError(Marshal, "interface failed validation", target, err)
		}
	}

	// handle self-encodable properties
	if en, ok := target.(encodableName); ok {
		var value string
		name := en.EncodeICalName()
		if ev, ok := target.(encodableValue); ok {
			value = ev.EncodeICalValue()
		} else {
			value = fmt.Sprintf("%s", target)
		}
		return marshalProperty(name, value), nil
	}

	// drill into interfaces and pointers.
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vkind := v.Kind()
	vtype := v.Type()

	var out []string
	if vkind == reflect.Slice || vkind == reflect.Array {

		// encode arrays early
		for i, n := 0, v.Len(); i < n; i++ {
			if encoded, err := Marshal(v.Index(i).Interface()); err != nil {
				msg := fmt.Sprintf("unable to encode component at index %d", i)
				return "", NewError(Marshal, msg, target, err)
			} else if encoded != "" {
				out = append(out, encoded)
			}
		}

	} else {

		// attempt to encode struct fields
		if vkind == reflect.Struct {

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
				var name, dvalue, value string
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
							dvalue = tags[1]
						}
					}
				}

				// make sure we have a name
				if name == "" {
					name = fs.Name
				}

				// now check to see if we have anything to encode
				if fv.IsValid() && !isEmptyValue(fv) {

					// check if field self-validates
					fi := fv.Interface()
					if validator, ok := fi.(validatable); ok {
						if err := validator.ValidateICalValue(); err != nil {
							msg := fmt.Sprintf("unable to encode field %s", fs.Name)
							return "", NewError(Marshal, msg, target, err)
						}
					}

					// check to see if field name is overridden
					if encoder, ok := fi.(encodableName); ok {
						if override := encoder.EncodeICalName(); override != "" {
							name = override
						}
					}

					// check to see if field value is overridden
					if encoder, ok := fi.(encodableValue); ok {
						if override := encoder.EncodeICalValue(); override != "" {
							value = override
						}
					} else {

						// drill into interfaces and pointers.
						for fv.Kind() == reflect.Interface || fv.Kind() == reflect.Ptr {
							fv = fv.Elem()
						}

						fkind := fv.Kind()

						if fkind == reflect.Slice || fkind == reflect.Array {

							// check for collections of components
							for i, n := 0, fv.Len(); i < n; i++ {
								if encoded, err := Marshal(fv.Index(i).Interface()); err != nil {
									msg := fmt.Sprintf("unable to encode field %s at index %d", fs.Name, i)
									return "", NewError(Marshal, msg, target, err)
								} else if encoded != "" {
									out = append(out, encoded)
								}
							}
							continue

						} else if fkind == reflect.Struct {

							// check for nested components
							if encoded, err := Marshal(fi); err != nil {
								msg := fmt.Sprintf("unable to encode field %s", fs.Name)
								return "", NewError(Marshal, msg, target, err)
							} else if encoded != "" {
								out = append(out, encoded)
							}
							continue

						} else {

							// otherwise just use the string representation
							value = fmt.Sprintf("%v", fi)

						}

					}

				}

				// make sure we have a value by this point
				if value == "" {

					// check for a default value
					if dvalue == "" {

						if required {
							// error out on required values
							msg := fmt.Sprintf("missing value for required field %s", fs.Name)
							return "", NewError(Marshal, msg, target, nil)
						} else if omitempty {
							// or skip omitted values
							continue
						}
					} else {
						value = dvalue
					}
				}

				// encode in the property
				out = append(out, marshalProperty(name, value))

			}

		}

		// all objects will begin and end here
		name := "V" + strings.ToUpper(vtype.Name())
		out = append([]string{marshalProperty("begin", name)}, out...)
		out = append(out, marshalProperty("end", name))

	}

	return strings.Join(out, Newline), nil

}
