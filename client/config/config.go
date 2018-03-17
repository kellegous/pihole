package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/ssh"
)

const (
	confFile   = "conf"
	prvKeyFile = "id_rsa"
	bitsInKey  = 4096
)

// Config ...
type Config struct {
	ID    string
	Hosts []string
	ToURL *url.URL
	Hub   struct {
		User string
		Addr string
	}
	dir string
}

type config struct {
	ID    string   `toml:"id"`
	Addr  string   `toml:"hub-addr"`
	Hosts []string `toml:"from-hosts"`
	ToURL string   `toml:"to-url"`
}

// PrivateKey ...
func (c *Config) PrivateKey() string {
	return filepath.Join(c.dir, prvKeyFile)
}

// PublicKey ...
func (c *Config) PublicKey() (string, error) {
	b, err := ioutil.ReadFile(c.PrivateKey())
	if err != nil {
		return "", err
	}

	blk, _ := pem.Decode(b)

	prv, err := x509.ParsePKCS1PrivateKey(blk.Bytes)
	if err != nil {
		return "", err
	}

	pub, err := ssh.NewPublicKey(&prv.PublicKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"ssh-rsa %s %s",
		base64.StdEncoding.EncodeToString(pub.Marshal()),
		c.ID), nil
}

func expandPath(dir, path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	p, err := filepath.Abs(filepath.Join(dir, path))
	if err != nil {
		return "", err
	}

	return p, nil
}

func parseAddr(addr string) (string, string, error) {
	var username string

	ix := strings.Index(addr, "@")
	if ix < 0 {
		u, err := user.Current()
		if err != nil {
			return "", "", err
		}
		username = u.Username
	} else {
		username = addr[:ix]
		addr = addr[ix+1:]
	}

	if !strings.Contains(addr, ":") {
		addr += ":22"
	}

	return username, addr, nil
}

func exportPrvKey(dst string, prv *rsa.PrivateKey) error {
	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	return pem.Encode(w, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(prv),
	})
}

func newID() string {
	var buf [8]byte
	if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
		log.Panic(err)
	}
	return hex.EncodeToString(buf[:])
}

// Create ...
func (c *Config) Create(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s exists", dir)
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	key, err := rsa.GenerateKey(rand.Reader, bitsInKey)
	if err != nil {
		return err
	}

	if err := exportPrvKey(filepath.Join(dir, prvKeyFile), key); err != nil {
		return err
	}

	s := config{
		ID:    newID(),
		Addr:  fmt.Sprintf("%s@%s", c.Hub.User, c.Hub.Addr),
		ToURL: c.ToURL.String(),
		Hosts: c.Hosts,
	}

	w, err := os.Create(filepath.Join(dir, confFile))
	if err != nil {
		return err
	}
	defer w.Close()

	if err := toml.NewEncoder(w).Encode(&s); err != nil {
		return err
	}

	return c.Read(dir)
}

// ReadFile ...
func (c *Config) Read(dir string) error {
	var s config
	var err error

	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}

	if _, err = toml.DecodeFile(filepath.Join(dir, confFile), &s); err != nil {
		return err
	}

	c.Hub.User, c.Hub.Addr, err = parseAddr(s.Addr)
	if err != nil {
		return err
	}

	c.ToURL, err = url.Parse(s.ToURL)
	if err != nil {
		return err
	}

	c.ID = s.ID
	c.Hosts = s.Hosts
	c.dir = dir

	return nil
}
