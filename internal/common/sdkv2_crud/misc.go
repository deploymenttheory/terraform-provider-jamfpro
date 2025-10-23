package common

import (
	"fmt"
	"reflect"
	"strconv"
)

// getIDField returns the value of the ID field in a response.
func getIDField(response any) (any, error) {
	v := reflect.ValueOf(response).Elem()

	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("ID field not found in response")
	}

	str, ok := idField.Interface().(string)
	if ok {
		return str, nil
	}

	integer, ok := idField.Interface().(int)
	if ok {
		return strconv.Itoa(integer), nil
	}
	return nil, fmt.Errorf("unsupported type")
}
