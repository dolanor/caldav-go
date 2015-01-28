package icalendar

import (
	"fmt"
	"github.com/taviti/caldav-go/utils"
	"reflect"
	"strings"
)

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

func dereferencePointerValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func extractTagFromValue(v reflect.Value) (string, error) {

	var tag string

	if encoder, ok := v.Interface().(canEncodeTag); ok {
		if t, err := encoder.EncodeICalTag(); err != nil {
			return "", utils.NewError(extractTagFromValue, "unable to extract tag from interface", v.Interface(), err)
		} else {
			tag = t
		}
	}

	if tag == "" {
		tag = fmt.Sprintf("v%s", dereferencePointerValue(v).Type().Name())
	}

	return strings.ToUpper(tag), nil

}
