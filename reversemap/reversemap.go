//go:build !solution

package reversemap

import "reflect"

func ReverseMap(forward interface{}) interface{} {
	forwardValue := reflect.ValueOf(forward)

	if forwardValue.Kind() != reflect.Map {
		panic("input is not a map")
	}

	keyType := forwardValue.Type().Key()
	valueType := forwardValue.Type().Elem()

	reversedValue := reflect.MakeMap(reflect.MapOf(valueType, keyType))
	keys := forwardValue.MapKeys()

	for _, key := range keys {
		value := forwardValue.MapIndex(key)
		reversedValue.SetMapIndex(value, key)
	}

	return reversedValue.Interface()
}
