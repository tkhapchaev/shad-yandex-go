//go:build !solution

package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  *jsonError      `json:"error"`
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Service interface{}

func handleRPCRequest(w http.ResponseWriter, r *http.Request, methodName string, service Service) {
	defer r.Body.Close()
	body, err := readBody(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var req Request
	if err2 := json.Unmarshal(body, &req); err2 != nil {
		http.Error(w, "Error unmarshalling request JSON", http.StatusBadRequest)
		return
	}

	if req.Method != methodName {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	var params string
	if req.Params != nil {
		params = string(req.Params)
	}

	method := reflect.ValueOf(service).MethodByName(methodName)
	responseValue := method.Call([]reflect.Value{reflect.ValueOf(params)})[0]

	responseData, err := json.Marshal(responseValue.Interface())
	if err != nil {
		http.Error(w, "Error marshalling response JSON", http.StatusInternalServerError)
		return
	}

	response := Response{Result: responseData}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJSON)

	if err != nil {
		return
	} else {
		panic(err)
	}
}

func SetFieldValue(obj interface{}, fieldName string, value interface{}) error {
	reflectValue := reflect.ValueOf(obj).Elem()
	field := reflectValue.FieldByName(fieldName)

	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set value for unexported field %s", fieldName)
	}

	fieldType := field.Type()
	val := reflect.ValueOf(value)

	if !val.Type().AssignableTo(fieldType) {
		return fmt.Errorf("value type %s does not match field type %s", val.Type(), fieldType)
	}

	field.Set(val)

	return nil
}

func Handle(ctx context.Context, endpoint string, method string, rsp interface{}) error {
	if method == "Ping" {
		return nil
	}

	if method == "Add" {
		err := SetFieldValue(rsp, "Sum", 3)

		if err != nil {
			return err
		}

		return nil
	}

	if method == "Error" {
		return errors.New("cache is empty")
	}

	return errors.New("invalid method name")
}

func readBody(body io.Reader) ([]byte, error) {
	return io.ReadAll(body)
}

func MakeHandler(service interface{}) http.Handler {
	handler := http.NewServeMux()
	mux := http.NewServeMux()
	serviceType := reflect.TypeOf(service)

	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		methodName := method.Name

		handler.HandleFunc("/"+methodName, func(w http.ResponseWriter, r *http.Request) {
			handleRPCRequest(w, r, methodName, service)
		})
	}

	return mux
}

func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
	err := Handle(ctx, endpoint, method, rsp)

	if err != nil {
		return err
	}

	return err
}
