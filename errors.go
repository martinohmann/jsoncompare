package jsoncompare

import "fmt"

// ComparatorError defines an error that can happen during Compare.
type ComparatorError struct {
	path jsonPath
	msg  string
}

// Error implements the error interface.
func (e ComparatorError) Error() string {
	return fmt.Sprintf("%s: %s", e.path, e.msg)
}

func newKeyError(path jsonPath, key string) error {
	msg := fmt.Sprintf("key %q does not exist in haystack", key)
	return ComparatorError{path, msg}
}

func newLengthError(path jsonPath, haystackLen, needleLen int) error {
	msg := fmt.Sprintf("length mismatch, expected %d, got %d", needleLen, haystackLen)
	return ComparatorError{path, msg}
}

func newTypeError(path jsonPath, haystack, needle interface{}) error {
	msg := fmt.Sprintf("type mismatch, expected %T, got %T", needle, haystack)
	return ComparatorError{path, msg}
}

func newValueError(path jsonPath, haystack, needle interface{}) error {
	msg := fmt.Sprintf("value mismatch, expected %v, got %v", needle, haystack)
	return ComparatorError{path, msg}
}
