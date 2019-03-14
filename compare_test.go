package jsoncompare

import "testing"

func TestCompare(t *testing.T) {
	subtreeComparator := NewComparator(MatchSubtree)

	cases := []struct {
		name             string
		comparator       *Comparator
		needle, haystack []byte
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:             "invalid haystack",
			needle:           []byte(`{"foo":"bar"}`),
			haystack:         []byte(`]`),
			expectError:      true,
			expectedErrorMsg: "invalid character ']' looking for beginning of value",
		},
		{
			name:             "invalid needle",
			needle:           []byte(`]`),
			haystack:         []byte(`{"foo":"bar"}`),
			expectError:      true,
			expectedErrorMsg: "invalid character ']' looking for beginning of value",
		},
		{
			name:        "nil",
			needle:      []byte(`null`),
			haystack:    []byte(`null`),
			expectError: false,
		},
		{
			name:        "string and number",
			needle:      []byte(`"string"`),
			haystack:    []byte(`2.0`),
			expectError: true,
		},
		{
			name:        "empty",
			needle:      []byte(`{}`),
			haystack:    []byte(`{}`),
			expectError: false,
		},
		{
			name:        "slice",
			needle:      []byte(`[]`),
			haystack:    []byte(`[]`),
			expectError: false,
		},
		{
			name:        "simple equal",
			needle:      []byte(`{"foo":"bar"}`),
			haystack:    []byte(`{"foo":"bar"}`),
			expectError: false,
		},
		{
			name:             "key mismatch",
			needle:           []byte(`{"foo":"bar"}`),
			haystack:         []byte(`{"bar":"bar"}`),
			expectError:      true,
			expectedErrorMsg: `$: key "foo" does not exist in haystack`,
		},
		{
			name:        "empty haystack",
			needle:      []byte(`{"foo":"bar"}`),
			haystack:    []byte(`{}`),
			expectError: true,
		},
		{
			name:        "empty needle",
			needle:      []byte(`{}`),
			haystack:    []byte(`{"foo":"bar"}`),
			expectError: false,
		},
		{
			name:             "different type",
			needle:           []byte(`{"foo":2}`),
			haystack:         []byte(`{"foo":["asdf"]}`),
			expectError:      true,
			expectedErrorMsg: "$.foo: type mismatch, expected float64, got []interface {}",
		},
		{
			name:             "different value",
			needle:           []byte(`{"foo":2}`),
			haystack:         []byte(`{"foo":3}`),
			expectError:      true,
			expectedErrorMsg: "$.foo: value mismatch, expected 2, got 3",
		},
		{
			name:        "nested",
			needle:      []byte(`{"foo":{"bar":"baz"}}`),
			haystack:    []byte(`{"foo":{"bar":"baz"}}`),
			expectError: false,
		},
		{
			name:        "nested with additional element",
			needle:      []byte(`{"foo":{"bar":"baz"}}`),
			haystack:    []byte(`{"foo":{"bar":"baz","qux":1}}`),
			expectError: false,
		},
		{
			name:        "nested with missing element",
			needle:      []byte(`{"foo":{"bar":"baz","quz":1}}`),
			haystack:    []byte(`{"foo":{"bar":"baz"}}`),
			expectError: true,
		},
		{
			name:        "nested type mismatch",
			needle:      []byte(`{"foo":{"bar":"baz","quz":1}}`),
			haystack:    []byte(`{"foo":2}`),
			expectError: true,
		},
		{
			name:        "slice",
			needle:      []byte(`{"foo":["bar"]}`),
			haystack:    []byte(`{"foo":["bar"]}`),
			expectError: false,
		},
		{
			name:             "type mismatch #1",
			needle:           []byte(`{"foo":["bar"]}`),
			haystack:         []byte(`{"foo":{"bar":"baz"}}`),
			expectError:      true,
			expectedErrorMsg: "$.foo: type mismatch, expected []interface {}, got map[string]interface {}",
		},
		{
			name:             "type mismatch #2",
			needle:           []byte(`{"foo":{"bar":"baz"}}`),
			haystack:         []byte(`{"foo":["bar"]}`),
			expectError:      true,
			expectedErrorMsg: "$.foo: type mismatch, expected map[string]interface {}, got []interface {}",
		},
		{
			name:             "type mismatch #3, deeply nested",
			needle:           []byte(`{"foo":{"bar":[{"baz":"qux"},{"qux":1}]}}`),
			haystack:         []byte(`{"foo":{"bar":[{"baz":"qux"},{"qux":[]}]}}`),
			expectError:      true,
			expectedErrorMsg: "$.foo.bar[1].qux: type mismatch, expected float64, got []interface {}",
		},
		{
			name:        "slice type mismatch",
			needle:      []byte(`{"foo":["bar"]}`),
			haystack:    []byte(`{"foo":[2]}`),
			expectError: true,
		},
		{
			name:        "slice length mismatch",
			needle:      []byte(`{"foo":["bar","baz"]}`),
			haystack:    []byte(`{"foo":["bar"]}`),
			expectError: true,
		},
		{
			name:        "slice order mismatch",
			needle:      []byte(`{"foo":["bar","baz"]}`),
			haystack:    []byte(`{"foo":["baz","bar"]}`),
			expectError: false,
		},
		{
			name:        "slice with additional haystack element",
			needle:      []byte(`{"foo":["bar"]}`),
			haystack:    []byte(`{"foo":["bar","baz"]}`),
			expectError: false,
		},
		{
			name:        "match map len",
			comparator:  NewComparator(MatchMapLen),
			needle:      []byte(`{"foo":"bar"}`),
			haystack:    []byte(`{"foo":"bar","baz":"qux"}`),
			expectError: true,
		},
		{
			name:        "match slice len",
			comparator:  NewComparator(MatchSliceLen),
			needle:      []byte(`{"foo":["bar"]}`),
			haystack:    []byte(`{"foo":["bar","baz"]}`),
			expectError: true,
		},
		{
			name:        "match slice order",
			comparator:  NewComparator(MatchSliceOrder),
			needle:      []byte(`{"foo":["baz","bar"]}`),
			haystack:    []byte(`{"foo":["bar","baz"]}`),
			expectError: true,
		},
		{
			name:        "exact match",
			comparator:  NewComparator(MatchStrict),
			needle:      []byte(`{"foo":["bar","baz"],"bar":1}`),
			haystack:    []byte(`{"foo":["bar","baz"],"bar":1}`),
			expectError: false,
		},
		{
			name:        "exact match, unordered map",
			comparator:  NewComparator(MatchStrict),
			needle:      []byte(`{"bar":1,"foo":["bar","baz"]}`),
			haystack:    []byte(`{"foo":["bar","baz"],"bar":1}`),
			expectError: false,
		},
		{
			name:        "should not match same element twice",
			needle:      []byte(`{"foo":[1, 1]}`),
			haystack:    []byte(`{"foo":[1, 2]}`),
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			comparator := tc.comparator
			if comparator == nil {
				comparator = subtreeComparator
			}

			err := comparator.Compare(tc.haystack, tc.needle)
			if tc.expectError {
				if err == nil {
					t.Fatalf(
						"expected error but got nil for needle %s and haystack %s",
						tc.needle,
						tc.haystack,
					)
				} else if tc.expectedErrorMsg != "" && err.Error() != tc.expectedErrorMsg {
					t.Fatalf(
						"expected error message %q but got %q for needle %s and haystack %s",
						tc.expectedErrorMsg,
						err.Error(),
						tc.needle,
						tc.haystack,
					)
				}
			} else if err != nil {
				t.Fatalf(
					"did not expect error but got %q for needle %s and haystack %s",
					err.Error(),
					tc.needle,
					tc.haystack,
				)
			}
		})
	}
}

func benchmarkCompare(b *testing.B, haystack, needle []byte) {
	for i := 0; i < b.N; i++ {
		Compare(haystack, needle)
	}
}

func BenchmarkScalar(b *testing.B) {
	benchmarkCompare(b, []byte(`"foo"`), []byte(`false`))
}

func BenchmarkObject(b *testing.B) {
	benchmarkCompare(b, []byte(`{"foo":"bar"}`), []byte(`{"foo":"bar"}`))
}

func BenchmarkArray(b *testing.B) {
	benchmarkCompare(b, []byte(`["foo","bar"]`), []byte(`["bar","foo"]`))
}

func BenchmarkDeeplyNested(b *testing.B) {
	benchmarkCompare(
		b,
		[]byte(`{"z":null,"foo":{"bar":[{"c":{"d":1}},{"baz":{"qux":1}}]},"a":1,"b":"2"}`),
		[]byte(`{"foo":{"bar":[{"baz":{"qux":1}},{"c":{"d":1}}]},"b":"2","a":1}`),
	)
}
