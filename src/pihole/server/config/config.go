package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	Store struct {
		Path string `toml:"path"`
	} `toml:"store"`

	Web struct {
		Addr string `toml:"addr"`
	} `toml:"web"`

	API struct {
		Addr string `toml:"addr"`
	} `toml:"api"`
}

func initialize(c *Config, base string) error {
	if filepath.IsAbs(c.Store.Path) {
		return nil
	}

	p, err := filepath.Abs(filepath.Join(base, c.Store.Path))
	if err != nil {
		return err
	}

	c.Store.Path = p

	return nil
}

// LoadFile ...
func (c *Config) LoadFile(filename string) error {
	if _, err := toml.DecodeFile(filename, c); err != nil {
		return err
	}

	p, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return err
	}

	return initialize(c, p)
}
