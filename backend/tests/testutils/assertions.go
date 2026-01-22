// Package testutils provides test utilities for integration tests
//
// File: assertions.go
// Description: Custom test assertions
package testutils

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		t.Errorf("expected values to be different, both are %v", actual)
	}
}

// AssertNil asserts that a value is nil
func AssertNil(t *testing.T, actual interface{}) {
	t.Helper()

	if actual != nil && !reflect.ValueOf(actual).IsNil() {
		t.Errorf("expected nil, got %v", actual)
	}
}

// AssertNotNil asserts that a value is not nil
func AssertNotNil(t *testing.T, actual interface{}) {
	t.Helper()

	if actual == nil || reflect.ValueOf(actual).IsNil() {
		t.Error("expected non-nil value, got nil")
	}
}

// AssertNoError asserts that an error is nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// AssertError asserts that an error is not nil
func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// AssertErrorContains asserts that an error contains a specific message
func AssertErrorContains(t *testing.T, err error, msg string) {
	t.Helper()

	if err == nil {
		t.Errorf("expected error containing %q, got nil", msg)
		return
	}

	if !strings.Contains(err.Error(), msg) {
		t.Errorf("expected error containing %q, got %q", msg, err.Error())
	}
}

// AssertTrue asserts that a value is true
func AssertTrue(t *testing.T, actual bool) {
	t.Helper()

	if !actual {
		t.Error("expected true, got false")
	}
}

// AssertFalse asserts that a value is false
func AssertFalse(t *testing.T, actual bool) {
	t.Helper()

	if actual {
		t.Error("expected false, got true")
	}
}

// AssertLen asserts the length of a slice, map, or string
func AssertLen(t *testing.T, obj interface{}, expected int) {
	t.Helper()

	val := reflect.ValueOf(obj)
	switch val.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
		if val.Len() != expected {
			t.Errorf("expected length %d, got %d", expected, val.Len())
		}
	default:
		t.Errorf("cannot get length of %T", obj)
	}
}

// AssertContains asserts that a string contains a substring
func AssertContains(t *testing.T, str, substr string) {
	t.Helper()

	if !strings.Contains(str, substr) {
		t.Errorf("expected %q to contain %q", str, substr)
	}
}

// AssertNotContains asserts that a string does not contain a substring
func AssertNotContains(t *testing.T, str, substr string) {
	t.Helper()

	if strings.Contains(str, substr) {
		t.Errorf("expected %q to not contain %q", str, substr)
	}
}

// AssertStatusCode asserts the HTTP status code
func AssertStatusCode(t *testing.T, expected int, actual int) {
	t.Helper()

	if expected != actual {
		t.Errorf("expected status code %d, got %d", expected, actual)
	}
}

// AssertJSONEqual asserts that two JSON values are equal
func AssertJSONEqual(t *testing.T, expected, actual string) {
	t.Helper()

	var expectedObj, actualObj interface{}

	if err := json.Unmarshal([]byte(expected), &expectedObj); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualObj); err != nil {
		t.Fatalf("failed to unmarshal actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedObj, actualObj) {
		t.Errorf("JSON not equal:\nexpected: %s\nactual: %s", expected, actual)
	}
}

// AssertGreater asserts that a value is greater than another
func AssertGreater(t *testing.T, value, than int) {
	t.Helper()

	if value <= than {
		t.Errorf("expected %d to be greater than %d", value, than)
	}
}

// AssertGreaterOrEqual asserts that a value is greater than or equal to another
func AssertGreaterOrEqual(t *testing.T, value, than int) {
	t.Helper()

	if value < than {
		t.Errorf("expected %d to be greater than or equal to %d", value, than)
	}
}

// AssertLess asserts that a value is less than another
func AssertLess(t *testing.T, value, than int) {
	t.Helper()

	if value >= than {
		t.Errorf("expected %d to be less than %d", value, than)
	}
}

// AssertLessOrEqual asserts that a value is less than or equal to another
func AssertLessOrEqual(t *testing.T, value, than int) {
	t.Helper()

	if value > than {
		t.Errorf("expected %d to be less than or equal to %d", value, than)
	}
}

// AssertEmpty asserts that a slice, map, or string is empty
func AssertEmpty(t *testing.T, obj interface{}) {
	t.Helper()

	val := reflect.ValueOf(obj)
	switch val.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
		if val.Len() != 0 {
			t.Errorf("expected empty, got length %d", val.Len())
		}
	default:
		t.Errorf("cannot check if %T is empty", obj)
	}
}

// AssertNotEmpty asserts that a slice, map, or string is not empty
func AssertNotEmpty(t *testing.T, obj interface{}) {
	t.Helper()

	val := reflect.ValueOf(obj)
	switch val.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
		if val.Len() == 0 {
			t.Error("expected non-empty value")
		}
	default:
		t.Errorf("cannot check if %T is not empty", obj)
	}
}
