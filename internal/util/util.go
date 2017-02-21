package util

import (
	"reflect"
	"regexp"
	"strings"
)

func Underscore(name string) string {
	var words []string
	var camelCase = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

	for _, submatch := range camelCase.FindAllStringSubmatch(name, -1) {
		if submatch[1] != "" {
			words = append(words, submatch[1])
		}
		if submatch[2] != "" {
			words = append(words, submatch[2])
		}
	}

	return strings.ToLower(strings.Join(words, "_"))
}

func DereferencePtrType(t reflect.Type) reflect.Type {
	kind := t.Kind()
	if kind == reflect.Ptr {
		return DereferencePtrType(t.Elem())
	}
	return t
}

func DereferencePtrValue(v reflect.Value) reflect.Value {
	kind := v.Kind()
	if kind == reflect.Ptr {
		return DereferencePtrValue(v.Elem())
	}
	return v
}
