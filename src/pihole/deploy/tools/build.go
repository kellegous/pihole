package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"pihole/deploy/config"
)

func setEnv(env []string, k, v string) []string {
	p := k + "="
	for i, ev := range env {
		if strings.HasPrefix(ev, p) {
			env[i] = fmt.Sprintf("%s=%s", k, v)
			return env
		}
	}
	return append(env, fmt.Sprintf("%s=%s", k, v))
}

func getEnv(b *config.Build) []string {
	env := os.Environ()
	env = setEnv(env, "GOOS", b.OS)
	env = setEnv(env, "GOARCH", b.Arch)
	return env
}

// Make ...
func Make(root string, b *config.Build) (string, error) {
	c := exec.Command("make")

	c.Dir = root
	c.Env = getEnv(b)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return "", err
	}

	return filepath.Join(root, "bin", b.Desc()), nil
}
