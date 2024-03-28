//go:build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func SetPrivateField(obj interface{}, name string, value interface{}) {
	val := reflect.ValueOf(obj).Elem().FieldByName(name)
	val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()
	val.Set(reflect.ValueOf(value))
}
