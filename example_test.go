package jsoncompare_test

import (
	"fmt"

	"github.com/martinohmann/jsoncompare"
)

func ExampleCompare() {
	needle := []byte(`{"foo": [1, 2]}`)
	haystack := []byte(`{"foo": ["1", "2"]}`)

	if err := jsoncompare.Compare(haystack, needle); err != nil {
		fmt.Println(err)
	}

	// Output:
	// $.foo[0]: type mismatch, expected float64, got string
}

func ExampleNewComparator_matchSubtree() {
	needle := []byte(`{"foo": [2, 1]}`)
	haystack := []byte(`{"foo": [1, 2, 3]}`)
	comparator := jsoncompare.NewComparator(jsoncompare.MatchSubtree)

	if err := comparator.Compare(haystack, needle); err == nil {
		fmt.Println("haystack contains needle")
	}

	// Output:
	// haystack contains needle
}

func ExampleNewComparator_matchSliceLen() {
	needle := []byte(`{"foo": [3, 2, 1]}`)
	haystack := []byte(`{"foo": [1, 2, 3]}`)
	comparator := jsoncompare.NewComparator(jsoncompare.MatchSliceLen)

	if err := comparator.Compare(haystack, needle); err == nil {
		fmt.Println("haystack contains needle")
	}

	// Output:
	// haystack contains needle
}

func ExampleNewComparator_matchLenStrict() {
	needle := []byte(`{"foo": [2, 1]}`)
	haystack := []byte(`{"foo": [1, 2, 3]}`)
	comparator := jsoncompare.NewComparator(jsoncompare.MatchLenStrict)

	if err := comparator.Compare(haystack, needle); err != nil {
		fmt.Println(err)
	}

	// Output:
	// $.foo: length mismatch, expected 2, got 3
}

func ExampleNewComparator_matchSliceStrict() {
	needle := []byte(`{"foo": [3, 2, 1]}`)
	haystack := []byte(`{"bar": "baz", "foo": [1, 2, 3]}`)
	comparator := jsoncompare.NewComparator(jsoncompare.MatchSliceStrict)

	if err := comparator.Compare(haystack, needle); err != nil {
		fmt.Println(err)
	}

	// Output:
	// $.foo[0]: value mismatch, expected 3, got 1
}
