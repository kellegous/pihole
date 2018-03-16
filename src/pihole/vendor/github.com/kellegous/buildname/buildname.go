package buildname

import (
	"fmt"
	"hash/fnv"
	"strings"
)

func hash(b []byte) uint64 {
	h := fnv.New64()
	h.Write(b)
	return h.Sum64()
}

// FromVersion generates a build name from a given version string.
func FromVersion(version string) string {
	na := uint64(len(adjectives))
	nn := uint64(len(nouns))
	n := nn * na
	h := hash([]byte(version)) % n

	ns := nouns[h%nn]
	as := adjectives[h/nn]

	return fmt.Sprintf("%s %s",
		strings.Title(as),
		strings.Title(ns))
}
