package server

import (
	"os"
	"path/filepath"

	"github.com/kellegous/pihole/deploy/config"
	"github.com/kellegous/pihole/deploy/tools"

	"go.uber.org/zap"
)

// Deploy ...
func Deploy(root, name string, b *config.Build) error {
	zap.L().Info("deploying server",
		zap.String("name", name),
		zap.String("os", b.OS),
		zap.String("arch", b.Arch))

	dir, err := tools.Make(root, b)
	if err != nil {
		return err
	}

	rem := tools.NewRemote(name)

	if err := rem.Install(
		filepath.Join(dir, "server"),
		"/usr/local/bin/piholed",
		os.FileMode(0755),
		"root"); err != nil {
		return err
	}

	if err := rem.Install(
		filepath.Join(root, "/etc/piholed.service"),
		"/lib/systemd/system/piholed.service",
		os.FileMode(0644),
		"root"); err != nil {
		return err
	}

	if err := rem.SSH([]string{
		"sudo",
		"systemctl",
		"daemon-reload",
	}); err != nil {
		return err
	}

	return rem.SSH([]string{
		"sudo",
		"systemctl",
		"restart",
		"piholed.service",
	})
}
