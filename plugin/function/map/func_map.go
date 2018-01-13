package funcMap

import (
	"reflect"
	"unicode"
)

const (
	ERROR_CODE_ARGS_NUM = 1000000 + iota
	ERROR_CODE_ARGS_TYPE
)

type FuncMap struct {
	funcParamsMap map[string][]string
	funcOwnerMap  map[string]interface{}
}

func NewFuncMap() *FuncMap {
	return &FuncMap{
		funcParamsMap: make(map[string][]string),
		funcOwnerMap:  make(map[string]interface{}),
	}
}

func IsExportedName(name string) bool {
	return name != "" && unicode.IsUpper(rune(name[0]))
}

func GetMethods(v interface{}) map[string][]string {
	funcMap := make(map[string][]string)
	reflectType := reflect.TypeOf(v)
	for i := 0; i < reflectType.NumMethod(); i++ {
		funcList := make([]string, 0)
		method := reflectType.Method(i)
		methodType := method.Type
		methodName := method.Name
		if !IsExportedName(methodName) {
			continue
		}
		for j := 1; j < methodType.NumIn(); j++ {
			params := methodType.In(j)
			funcList = append(funcList, params.String())
		}
		funcMap[methodName] = funcList
	}
	return funcMap
}

func (f FuncMap) Invoke(funcName string, args ...interface{}) (result []interface{}, err interface{}) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = err1
		}
	}()
	if receiver, ok := f.funcOwnerMap[funcName]; ok {
		if pts, ok := f.funcParamsMap[funcName]; ok {
			if len(pts) != len(args) {
				panic(ERROR_CODE_ARGS_NUM)
			}
			for index, paramType := range pts {
				if reflect.TypeOf(args[index]).String() != paramType {
					panic(ERROR_CODE_ARGS_TYPE)
				}
			}

			inputs := make([]reflect.Value, len(args))
			for i, _ := range args {
				inputs[i] = reflect.ValueOf(args[i])
			}
			rv := reflect.ValueOf(receiver).MethodByName(funcName).Call(inputs)
			result = make([]interface{}, len(rv))
			for k, v := range rv {
				result[k] = v.Interface()
			}
		}
	}

	return
}

func (f FuncMap) Register(v interface{}) {
	fs := GetMethods(v)
	for key, value := range fs {
		f.funcParamsMap[key] = value
		f.funcOwnerMap[key] = v
	}
}
