package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Build ...
type Build struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

// Desc ...
func (b *Build) Desc() string {
	return fmt.Sprintf("%s-%s", b.OS, b.Arch)
}

// Info ...
type Info struct {
	Servers map[string]*Build
	Clients map[string]*Build
}

// ReadFile ...
func (i *Info) ReadFile(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	return i.Read(r)
}

// Read ...
func (i *Info) Read(r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}
