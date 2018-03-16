package tools

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID ...
func NewID(n int) string {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
