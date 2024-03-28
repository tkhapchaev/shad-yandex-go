//go:build !solution

package jsonlist

import (
	"encoding/json"
	"io"
	"reflect"
)

func Marshal(w io.Writer, slice interface{}) error {
	enc := json.NewEncoder(w)
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: val.Type()}
	}

	for i := 0; i < val.Len(); i++ {
		if err := enc.Encode(val.Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

func Unmarshal(r io.Reader, slice interface{}) error {
	dec := json.NewDecoder(r)
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: val.Type()}
	}

	elemType := val.Elem().Type().Elem()

	for {
		item := reflect.New(elemType).Interface()

		if err := dec.Decode(item); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		val.Elem().Set(reflect.Append(val.Elem(), reflect.ValueOf(item).Elem()))
	}

	return nil
}
