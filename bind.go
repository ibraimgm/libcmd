package libcfg

import (
	"fmt"
	"reflect"
	"strings"
)

type bindingData struct {
	opt       *optEntry //long, short,help, val
	variables []string  //env var list
	val       *variant  // variant***
}

func collectBindings(target interface{}) ([]bindingData, error) {
	v := reflect.ValueOf(target)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, fmt.Errorf("cannot bind to a value that is not a pointer or is nil")
	}

	v = v.Elem()
	t := reflect.TypeOf(target).Elem()
	data := make([]bindingData, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// ignore if e can't modify it
		if !value.CanSet() {
			continue
		}

		// build a variant
		val := extractVariant(field, value)

		// use existing variant in opt field
		opt := extractOpt(field)
		if opt != nil {
			opt.val = val
		}

		variables := extractEnv(field)

		data = append(data, bindingData{
			val:       val,
			opt:       opt,
			variables: variables,
		})
	}

	return data, nil
}

func extractVariant(field reflect.StructField, value reflect.Value) *variant {
	defaultValue := reflect.Zero(value.Type())

	if defaultStr := strings.TrimSpace(field.Tag.Get("default")); defaultStr != "" {
		if converted, err := valueAsKind(defaultStr, value.Kind(), value.Type()); err == nil {
			defaultValue = converted
		}
	}

	return varFromReflect(value, defaultValue)
}

func extractOpt(field reflect.StructField) *optEntry {
	long := strings.TrimSpace(field.Tag.Get("long"))
	short := strings.TrimSpace(field.Tag.Get("short"))
	help := strings.TrimSpace(field.Tag.Get("help"))

	if long == "" && short == "" && help == "" {
		return nil
	}

	return &optEntry{
		long:  long,
		short: short,
		help:  help,
	}
}

func extractEnv(field reflect.StructField) []string {
	xs := strings.Split(field.Tag.Get("env"), ",")

	if len(xs) == 0 {
		return xs
	}

	result := make([]string, 0, len(xs))
	for i := range xs {
		x := strings.TrimSpace(xs[i])

		if x != "" {
			result = append(result, x)
		}
	}

	return result
}
