package types

import (
	"reflect"
	"regexp"

	"github.com/k0kubun/pp"
)

var (
	emailRegexp = regexp.MustCompile(".*(@.*)")
	Registry    = map[string]reflect.Type{} // this is the registry of types by name
)

// add a type to the registry
func RegisterType(t reflect.Type) {
	name := t.Name()
	Registry[name] = t
	pp.Println("registerTypeMap: ", Registry)
}

func GetType(object interface{}) string {
	if t := reflect.TypeOf(object); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
