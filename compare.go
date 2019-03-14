// Package jsoncompare is a helper utility to compare two byte slices
// containing json with each other. The matching behaviour is configurable.
// This is mostly useful for assertions in tests where we want to validate that
// a given haystack contains a needle but we do not mind if the haystack
// contains additional data not present in needle.
package jsoncompare

import (
	"encoding/json"
	"reflect"
)

// MatchMode defines the type of flags to configure the matching behaviour.
type MatchMode uint

const (
	// MatchSubtree will just match that the given needle is present, ignoring
	// any excess slice and map elements in haystack. Also the order of slice
	// elements is ignored.
	MatchSubtree MatchMode = 0

	// MatchMapLen will assert that maps in needle and haystack have the same
	// length.
	// MatchSliceLen works the same as MatchMapLen, only for slices.
	// MatchSliceOrder will make sure that slice ordering is that same in
	// needle and haystack.
	MatchMapLen MatchMode = 1 << iota
	MatchSliceLen
	MatchSliceOrder

	// MatchStrict will enforce exact matches of needle and haystack.
	MatchStrict = MatchMapLen | MatchSliceLen | MatchSliceOrder

	// MatchLenStrict will enforce that map and slice lengths are identical in
	// needle and haystack.
	MatchLenStrict = MatchMapLen | MatchSliceLen

	// MatchSliceStrict will enforce that slices are of the same lengths and
	// have the same order of elements.
	MatchSliceStrict = MatchSliceLen | MatchSliceOrder
)

// Comparator compares two json byte slices.
type Comparator struct {
	matchMode MatchMode
}

// DefaultComparator is a Comparator that strictly compares needle and haystack
// for equality.
var DefaultComparator = NewComparator(MatchStrict)

// NewComparator creates a new Comparator with MatchMode.
func NewComparator(m MatchMode) *Comparator {
	return &Comparator{matchMode: m}
}

// Compare checks if haystack contains needle. It uses the DefaultComparator.
func Compare(haystack, needle []byte) error {
	return DefaultComparator.Compare(haystack, needle)
}

// Compare checks if haystack contains needle. It takes the Comparator options
// into account when making assertions.
func (c *Comparator) Compare(haystack, needle []byte) error {
	var nval, hval interface{}

	if err := json.Unmarshal(needle, &nval); err != nil {
		return err
	}

	if err := json.Unmarshal(haystack, &hval); err != nil {
		return err
	}

	return c.compare(hval, nval, rootJsonPath)
}

// hasMode returns true if the Comparator is configured with the given
// MatchMode.
func (c *Comparator) hasMode(m MatchMode) bool {
	return c.matchMode&m == m
}

func (c *Comparator) compare(haystack, needle interface{}, path jsonPath) error {
	switch nval := needle.(type) {
	case map[string]interface{}:
		if hval, ok := haystack.(map[string]interface{}); ok {
			return c.compareMap(hval, nval, path)
		}
	case []interface{}:
		if hval, ok := haystack.([]interface{}); ok {
			return c.compareSlice(hval, nval, path)
		}
	default:
		if reflect.TypeOf(haystack) == reflect.TypeOf(needle) {
			if haystack != needle {
				return newValueError(path, haystack, needle)
			}

			return nil
		}
	}

	return newTypeError(path, haystack, needle)
}

func (c *Comparator) compareMap(haystack, needle map[string]interface{}, path jsonPath) error {
	hlen, nlen := len(haystack), len(needle)
	if hlen < nlen || (c.hasMode(MatchMapLen) && hlen != nlen) {
		return newLengthError(path, hlen, nlen)
	}

	for key, nval := range needle {
		if hval, ok := haystack[key]; !ok {
			return newKeyError(path, key)
		} else if err := c.compare(hval, nval, path.withKey(key)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Comparator) compareSlice(haystack, needle []interface{}, path jsonPath) error {
	hlen, nlen := len(haystack), len(needle)
	if hlen < nlen || (c.hasMode(MatchSliceLen) && hlen != nlen) {
		return newLengthError(path, hlen, nlen)
	}

	if c.hasMode(MatchSliceOrder) {
		for i, val := range needle {
			if err := c.compare(haystack[i], val, path.withIndex(i)); err != nil {
				return err
			}
		}

		return nil
	}

	for i, val := range needle {
		var err error

		for j, hval := range haystack {
			if err = c.compare(hval, val, path.withIndex(i)); err == nil {
				haystack = append(haystack[:j], haystack[j+1:]...)
				break
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
