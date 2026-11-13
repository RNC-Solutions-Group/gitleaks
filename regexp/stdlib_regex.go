//go:build !gore2regex

package regexp

import (
	re "regexp"
)

const Version = "stdlib"

type Regexp = re.Regexp

// Compile compiles a user-supplied expression and returns the parse error
// instead of panicking on malformed configuration.
func Compile(str string) (*re.Regexp, error) {
	return re.Compile(str)
}

func MustCompile(str string) *re.Regexp {
	return re.MustCompile(str)
}
