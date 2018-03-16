package client

import (
	"os"
	"path/filepath"
	"pihole/deploy/config"
	"pihole/deploy/tools"

	"go.uber.org/zap"
)

// Deploy ...
func Deploy(root, name string, b *config.Build) error {
	zap.L().Info("deploying client",
		zap.String("name", name),
		zap.String("os", b.OS),
		zap.String("arch", b.Arch))

	dir, err := tools.Make(root, b)
	if err != nil {
		return err
	}

	rem := tools.NewRemote(name)

	if err := rem.Install(
		filepath.Join(dir, "client"),
		"/usr/local/bin/pihole",
		os.FileMode(0755),
		"root"); err != nil {
		return err
	}

	if err := rem.Install(
		filepath.Join(root, "etc/pihole.service"),
		"/lib/systemd/system/pihole.service",
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
		"pihole.service",
	})
}
