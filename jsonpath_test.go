package jsoncompare

import "testing"

func TestJsonPath(t *testing.T) {
	path := rootJsonPath.
		withKey("foo").
		withIndex(2).
		withKey("bar").
		withKey("baz").
		withIndex(3).
		withIndex(0)

	expected := "$.foo[2].bar.baz[3][0]"

	if path.String() != expected {
		t.Fatalf("expected json path to be %q, got %q", expected, path.String())
	}
}
