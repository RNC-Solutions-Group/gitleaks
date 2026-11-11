//go:build gore2regex

package regexp

import (
	re "github.com/wasilibs/go-re2"
)

const Version = "github.com/wasilibs/go-re2"

type Regexp = re.Regexp

// Compile compiles a user-supplied expression and returns the parse error
// instead of panicking on malformed configuration.
func Compile(str string) (*re.Regexp, error) {
	return re.Compile(str)
}

func MustCompile(str string) *re.Regexp {
	return re.MustCompile(str)
}
