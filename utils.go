package gosmock

import "reflect"

func isNullable(t reflect.Type) bool {
	if t == nil {
		return true
	}
	k := t.Kind()
	if k == reflect.Ptr || k == reflect.Interface || k == reflect.Map ||
		k == reflect.Slice || k == reflect.Chan || k == reflect.Func {
		return true
	}
	return false
}
