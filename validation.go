package goutils

import "reflect"

// IsPtr is pointer
func IsPtr(i interface{}) bool {
	if reflect.ValueOf(i).Kind() == reflect.Ptr {
		return true
	}

	return false
}

// IsArr is array
func IsArr(val interface{}) bool {
	v := reflect.ValueOf(val)
	return v.Kind() == reflect.Array
}
