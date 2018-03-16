package build

import (
	"github.com/kellegous/buildname"
)

var (
	// SHA ...
	SHA string

	// Ref ...
	Ref string
)

// Name ...
func Name() string {
	return buildname.FromVersion(SHA)
}
