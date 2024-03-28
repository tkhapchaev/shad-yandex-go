//go:build !solution

package testequal

import (
	"bytes"
	"fmt"
	"reflect"
)

func ErrorFormatted(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	message := ""

	if len(msgAndArgs) > 0 {
		if format, ok := msgAndArgs[0].(string); ok {
			message = fmt.Sprintf(format, msgAndArgs[1:]...)
		} else if body, ok := msgAndArgs[0].(string); ok {
			message = body
		}
	}

	format :=
		`
		expected: %v
        actual  : %v
        message : %v`

	t.Errorf(format, expected, actual, message)
}

func Equal(expected, actual interface{}) bool {
	switch o := expected.(type) {
	case uint, int, uint8, int8, uint16, int16, uint64, int64:
		return o == actual
	case struct{}:
		return false
	case map[string]string:
		a, ok := actual.(map[string]string)

		if !ok {
			return false
		}

		if len(o) != len(a) {
			return false
		}

		if len(a) == 0 {
			return false
		}

		for k, v := range o {
			actualValue, ok := a[k]

			if !ok || v != actualValue {
				return false
			}
		}

		return true
	}

	expectedValue, ok := expected.([]byte)

	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	actualValue, ok := actual.([]byte)

	if !ok {
		return false
	}

	if expectedValue == nil || actualValue == nil {
		return expectedValue == nil && actualValue == nil
	}

	return bytes.Equal(expectedValue, actualValue)
}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if Equal(expected, actual) {
		return true
	}

	ErrorFormatted(t, expected, actual, msgAndArgs...)

	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if !Equal(expected, actual) {
		return true
	}

	ErrorFormatted(t, expected, actual, msgAndArgs...)

	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if Equal(expected, actual) {
		return
	}

	ErrorFormatted(t, expected, actual, msgAndArgs...)
	t.FailNow()
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if !Equal(expected, actual) {
		return
	}

	ErrorFormatted(t, expected, actual, msgAndArgs...)
	t.FailNow()
}
