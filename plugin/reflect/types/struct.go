package types

import (
	"reflect"
	// "github.com/fatih/structs"
	// "github.com/k0kubun/pp"
)

// create a new object by name, returning it as interface{}
func NewDataStructByName(name string) interface{} {
	t, found := Registry[name]
	if !found {
		panic("name not found!")
	}
	return reflect.New(t).Elem().Interface()
}

func IsCustomType(t reflect.Type) bool {
	if t.PkgPath() != "" {
		return true
	}
	if k := t.Kind(); k == reflect.Array || k == reflect.Chan || k == reflect.Map ||
		k == reflect.Ptr || k == reflect.Slice {
		return IsCustomType(t.Elem()) || k == reflect.Map && IsCustomType(t.Key())
	} else if k == reflect.Struct {
		for i := t.NumField() - 1; i >= 0; i-- {
			if IsCustomType(t.Field(i).Type) {
				return true
			}
		}
	}
	return false
}

// indirect returns the actual value if the given value is a pointer
func Indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}
