package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Remote ...
type Remote struct {
	Addr string
}

// NewRemote ...
func NewRemote(addr string) *Remote {
	return &Remote{Addr: addr}
}

func inheritIO(c *exec.Cmd) {
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
}

// Install ...
func (r *Remote) Install(
	src, dst string,
	perm os.FileMode,
	owner string) error {

	tmp := filepath.Join("/tmp", NewID(16))

	if err := r.SCP(src, tmp); err != nil {
		return err
	}

	if err := r.Chown(tmp, owner, false); err != nil {
		return err
	}

	if err := r.Chmod(tmp, perm, false); err != nil {
		return err
	}

	return r.SSH([]string{
		"sudo",
		"mv",
		tmp,
		dst,
	})
}

// Chmod ...
func (r *Remote) Chmod(dst string, mode os.FileMode, recursive bool) error {
	args := []string{
		"sudo",
		"chmod",
	}

	if recursive {
		args = append(args, "-R")
	}

	args = append(args, fmt.Sprintf("%03o", mode), dst)

	return r.SSH(args)
}

// Chown ...
func (r *Remote) Chown(dst string, owner string, recursive bool) error {
	args := []string{
		"sudo",
		"chown",
	}

	if recursive {
		args = append(args, "-R")
	}

	args = append(args, owner, dst)

	return r.SSH(args)
}

// SCP ...
func (r *Remote) SCP(src, dst string) error {
	c := exec.Command("scp", src, fmt.Sprintf("%s:%s", r.Addr, dst))
	inheritIO(c)
	return c.Run()
}

// SSH ...
func (r *Remote) SSH(cmd []string) error {
	args := []string{
		r.Addr,
	}

	for _, arg := range cmd {
		args = append(args, arg)
	}

	c := exec.Command("ssh", args...)
	inheritIO(c)
	return c.Run()
}
