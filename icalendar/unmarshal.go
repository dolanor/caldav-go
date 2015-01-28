package icalendar

import (
	"fmt"
	"github.com/taviti/caldav-go/utils"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var _ = log.Print

type canDecodeValue interface {
	DecodeICalValue(value string) error
}

type canDecodeParams interface {
	DecodeICalParams(value map[string]string) error
}

type token struct {
	name       string
	components map[string][]*token
	properties map[string][]*property
}

func tokenize(encoded string) (*token, error) {
	if encoded = strings.TrimSpace(encoded); encoded == "" {
		return nil, utils.NewError(tokenize, "no content to tokenize", encoded, nil)
	}
	return tokenizeSlice(strings.Split(encoded, Newline))
}

func tokenizeSlice(slice []string, name ...string) (*token, error) {

	tok := new(token)
	size := len(slice)

	if len(name) > 0 {
		tok.name = name[0]
	} else if size <= 0 {
		return nil, utils.NewError(tokenizeSlice, "token has no content", slice, nil)
	}

	tok.properties = make(map[string][]*property, 0)
	tok.components = make(map[string][]*token, 0)

	for i := 0; i < size; i++ {

		line := slice[i]
		prop := unmarshalProperty(line)

		if strings.EqualFold(prop.Name, "begin") {
			for j := i; j < size; j++ {
				end := strings.Replace(line, "BEGIN", "END", 1)
				if slice[j] == end {
					if component, err := tokenizeSlice(slice[i+1:j], prop.Value); err != nil {
						msg := fmt.Sprintf("unable to tokenize %s component", prop.Value)
						return nil, utils.NewError(tokenizeSlice, msg, slice, err)
					} else {
						if existing, ok := tok.components[prop.Value]; ok {
							tok.components[prop.Value] = []*token{component}
						} else {
							tok.components[prop.Value] = append(existing, component)
						}
						i = j
						break
					}
				}
			}
		} else if existing, ok := tok.properties[prop.Name]; ok {
			tok.properties[prop.Name] = []*property{prop}
		} else {
			tok.properties[prop.Name] = append(existing, prop)
		}

	}

	return tok, nil

}

func hydrateProperty(v reflect.Value, prop *property) error {

	// create a new object to hold the property value
	var vdref = dereferencePointerValue(v)
	var vkind = vdref.Kind()
	var vtype = vdref.Type()
	var isArray = vkind == reflect.Array || vkind == reflect.Slice

	var vnew reflect.Value
	if isArray {
		vnew = reflect.New(vtype.Elem())
	} else {
		vnew = reflect.New(vtype)
	}

	var hasValue = false
	var vint interface{}
	var vnewdref = dereferencePointerValue(vnew)

	if decoder, ok := v.Interface().(canDecodeValue); ok {
		vint = decoder
		if err := decoder.DecodeICalValue(prop.Value); err != nil {
			return utils.NewError(hydrateProperty, "error decoding property value", prop, err)
		} else {
			// interface handled decoding, no need for new value
			hasValue = true
		}
	} else if decoder, ok := vnew.Interface().(canDecodeValue); ok {
		vint = decoder
		// pointer value handles decoding...
		if err := decoder.DecodeICalValue(prop.Value); err != nil {
			return utils.NewError(hydrateProperty, "error decoding property value", prop, err)
		}
	} else {
		// we handle decoding...
		switch vnewdref.Kind() {
		case reflect.Bool:
			if i, err := strconv.ParseBool(prop.Value); err != nil {
				return utils.NewError(hydrateProperty, "unable to decode bool", prop, err)
			} else {
				vnewdref.SetBool(i)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if i, err := strconv.ParseInt(prop.Value, 10, 64); err != nil {
				return utils.NewError(hydrateProperty, "unable to decode int", prop, err)
			} else {
				vnewdref.SetInt(i)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if i, err := strconv.ParseUint(prop.Value, 10, 64); err != nil {
				return utils.NewError(hydrateProperty, "unable to decode uint", prop, err)
			} else {
				vnewdref.SetUint(i)
			}
		case reflect.Float32, reflect.Float64:
			if i, err := strconv.ParseFloat(prop.Value, 64); err != nil {
				return utils.NewError(hydrateProperty, "unable to decode float", prop, err)
			} else {
				vnewdref.SetFloat(i)
			}
		case reflect.String:
			vnewdref.SetString(prop.Value)
		}
	}

	// decode any params, if supported
	if len(prop.Params) > 0 {
		if decoder, ok := vint.(canDecodeParams); ok {
			if err := decoder.DecodeICalParams(prop.Params); err != nil {
				return utils.NewError(hydrateProperty, "error decoding property parameters", prop, err)
			}
		}
	}

	// finish with any validation
	if validator, ok := vint.(canValidateValue); ok {
		if err := validator.ValidateICalValue(); err != nil {
			return utils.NewError(hydrateProperty, "error validating property value", prop, err)
		}
	}

	// set the pointer to the new value
	if !hasValue {
		if isArray {
			vdref.Set(reflect.Append(vdref, vnewdref))
		} else {
			vdref.Set(vnewdref)
		}
	}

	return nil

}

func hydrateNestedComponent(v reflect.Value, component *token) error {

	// create a new object to hold the property value
	var vdref = dereferencePointerValue(v)
	var vkind = vdref.Kind()
	var vtype = vdref.Type()
	var isArray = vkind == reflect.Array || vkind == reflect.Slice

	var vnew reflect.Value
	if isArray {
		vnew = reflect.New(vtype.Elem())
	} else {
		vnew = reflect.New(vtype)
	}

	if err := hydrateValue(vnew, component); err != nil {
		return utils.NewError(hydrateProperty, "unable to decode component", component, err)
	}

	// set the pointer to the new value
	if isArray {
		v.Set(reflect.Append(vdref, vnew))
	} else {
		v.Set(vnew)
	}

	return nil

}

func hydrateProperties(v reflect.Value, component *token) error {

	vdref := dereferencePointerValue(v)
	vtype := vdref.Type()
	vkind := vdref.Kind()

	if vkind != reflect.Struct {
		return utils.NewError(hydrateProperties, "unable to hydrate properties of non-struct", v, nil)
	}

	n := vtype.NumField()
	for i := 0; i < n; i++ {

		prop := propertyFromStructField(vtype.Field(i))
		if prop == nil {
			continue // skip if field is ignored
		}

		if properties, ok := component.properties[prop.Name]; ok {
			// hydrate property values
			for _, prop := range properties {
				if err := hydrateProperty(vdref.Field(i), prop); err != nil {
					msg := fmt.Sprintf("unable to hydrate property %s", prop.Name)
					return utils.NewError(hydrateProperties, msg, v, err)
				}
			}
		} else if components, ok := component.components[prop.Name]; ok {
			// hydrate nested components
			for _, comp := range components {
				if err := hydrateNestedComponent(vdref.Field(i), comp); err != nil {
					msg := fmt.Sprintf("unable to hydrate property %s", prop.Name)
					return utils.NewError(hydrateProperties, msg, v, err)
				}
			}
		}

	}

	return nil

}

func hydrateComponent(v reflect.Value, component *token) error {
	if tag, err := extractTagFromValue(v); err != nil {
		return utils.NewError(hydrateComponent, "error extracting tag from value", component, err)
	} else if tag != component.name {
		msg := fmt.Sprintf("expected %s and found %s", tag, component.name)
		return utils.NewError(hydrateComponent, msg, component, nil)
	} else if err := hydrateProperties(v, component); err != nil {
		return utils.NewError(hydrateComponent, "unable to hydrate properties", component, err)
	}
	return nil
}

func hydrateComponents(v reflect.Value, components []*token) error {
	vdref := dereferencePointerValue(v)
	for i, component := range components {
		velem := reflect.New(vdref.Type().Elem())
		if err := hydrateComponent(velem, component); err != nil {
			msg := fmt.Sprintf("unable to hydrate component %d", i)
			return utils.NewError(hydrateComponent, msg, component, err)
		} else {
			v.Set(reflect.Append(vdref, velem))
		}
	}
	return nil
}

func hydrateValue(v reflect.Value, component *token) error {

	if !v.IsValid() || v.Kind() != reflect.Ptr {
		return utils.NewError(hydrateValue, "unmarshal target must be a valid pointer", v, nil)
	}

	// handle any encodable properties
	if encoder, isprop := v.Interface().(canEncodeName); isprop {
		if name, err := encoder.EncodeICalName(); err != nil {
			return utils.NewError(hydrateValue, "unable to lookup property name", v, err)
		} else if properties, found := component.properties[name]; !found || len(properties) == 0 {
			return utils.NewError(hydrateValue, "no matching propery values found for "+name, v, nil)
		} else if len(properties) > 1 {
			return utils.NewError(hydrateValue, "more than one property value matches single property interface", v, nil)
		} else {
			return hydrateProperty(v, properties[0])
		}
	}

	// handle components
	vkind := dereferencePointerValue(v).Kind()
	if tag, err := extractTagFromValue(v); err != nil {
		return utils.NewError(hydrateValue, "unable to extract component tag", v, err)
	} else if components, found := component.components[tag]; !found || len(components) == 0 {
		msg := fmt.Sprintf("unable to find matching component for %s", tag)
		return utils.NewError(hydrateValue, msg, v, nil)
	} else if vkind == reflect.Array || vkind == reflect.Slice {
		return hydrateComponents(v, components)
	} else if len(components) > 1 {
		return utils.NewError(hydrateValue, "non-array interface provided but more than one component found!", v, nil)
	} else {
		return hydrateComponent(v, components[0])
	}

}

// decodes encoded icalendar data into a native interface
func Unmarshal(encoded string, into interface{}) error {
	if component, err := tokenize(encoded); err != nil {
		return utils.NewError(Unmarshal, "unable to tokenize encoded data", encoded, err)
	} else {
		return hydrateValue(reflect.ValueOf(into), component)
	}
}
