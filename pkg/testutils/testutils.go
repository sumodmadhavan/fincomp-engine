package testutils

import (
	"testing"
)

// AssertEqual checks if two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

// AssertNoError checks if the error is nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

// AssertError checks if the error is not nil
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
