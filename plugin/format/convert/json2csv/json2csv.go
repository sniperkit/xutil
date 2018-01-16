// Package json2csv provides JSON to CSV functions.
package json2csv

import (
	"errors"
	"reflect"
	// "github.com/thoas/go-funk"
	// ureflect "github.com/XieZhendong/ureflect"
	// "github.com/oleiade/reflections"
)

// JSON2CSV converts JSON to CSV.
func JSON2CSV(data interface{}) ([]KeyValue, error) {
	results := []KeyValue{}
	v := valueOf(data)
	//v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Map:
		if v.Len() > 0 {
			result, err := flatten(v)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}

	case reflect.Slice:
		if isObjectArray(v) {
			for i := 0; i < v.Len(); i++ {
				result, err := flatten(v.Index(i))
				if err != nil {
					return nil, err
				}
				results = append(results, result)
			}
		} else if v.Len() > 0 {
			result, err := flatten(v)
			if err != nil {
				return nil, err
			}
			if result != nil {
				results = append(results, result)
			}
		}

	default:
		return nil, errors.New("Unsupported JSON structure.")

	}

	return results, nil
}

func ToMap(data interface{}) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)
	v := valueOf(data)
	//v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Map:
		if v.Len() > 0 {
			result, err := flatten2(v)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}

	case reflect.Slice:
		if isObjectArray(v) {
			for i := 0; i < v.Len(); i++ {
				result, err := flatten2(v.Index(i))
				if err != nil {
					return nil, err
				}
				results = append(results, result)
			}
		} else if v.Len() > 0 {
			result, err := flatten2(v)
			if err != nil {
				return nil, err
			}
			if result != nil {
				results = append(results, result)
			}
		}

	default:
		return nil, errors.New("Unsupported JSON structure.")

	}

	return results, nil
}

func isObjectArray(obj interface{}) bool {
	value := valueOf(obj)
	if value.Kind() != reflect.Slice {
		return false
	}

	len := value.Len()
	if len == 0 {
		return false
	}
	for i := 0; i < len; i++ {
		if valueOf(value.Index(i)).Kind() != reflect.Map {
			return false
		}
	}

	return true
}
