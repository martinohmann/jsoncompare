package jsoncompare

import "strconv"

// rootJsonPath is the path of the topmost level in a json object.
var rootJsonPath = jsonPath{"$"}

// jsonPath is a container to keep track of the current nesting level in a json
// object. This is used to provide the user with the exact path in the object
// where an error during Compare occured.
type jsonPath struct {
	path string
}

// withIndex returns a new jsonPath for the given array index.
func (j jsonPath) withIndex(index int) jsonPath {
	return jsonPath{j.path + "[" + strconv.Itoa(index) + "]"}
}

// withKey returns a new jsonPath for the given map key.
func (j jsonPath) withKey(key string) jsonPath {
	return jsonPath{j.path + "." + key}
}

// String implements the fmt.Stringer interface.
func (j jsonPath) String() string {
	return j.path
}
