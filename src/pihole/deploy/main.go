package main

import (
	"flag"
	"os"
	"path/filepath"

	"pihole/build"
	"pihole/deploy/client"
	"pihole/deploy/config"
	"pihole/deploy/server"
	"pihole/logging"

	"go.uber.org/zap"
)

func findRoot(root string) (string, error) {
	if root != "" {
		if filepath.IsAbs(root) {
			return root, nil
		}

		return filepath.Abs(root)
	}

	return filepath.Abs(
		filepath.Join(filepath.Dir(os.Args[0]), "../.."))
}

func main() {
	flagRoot := flag.String("root", "",
		"project root directory")
	flagConf := flag.String("conf", "deploy.conf",
		"config file")
	flagToClients := flag.Bool("to-clients", true,
		"whether to deploy clients")
	flagToServers := flag.Bool("to-servers", true,
		"whether to deploy servers")
	flag.Parse()

	logging.MustSetup()

	root, err := findRoot(*flagRoot)
	if err != nil {
		zap.L().Fatal("unable to find root",
			zap.String("flag", *flagRoot),
			zap.String("path", os.Args[0]))
	}

	zap.L().Info("deploying",
		zap.String("sha", build.SHA),
		zap.String("ref", build.Ref),
		zap.String("name", build.Name()))

	var cfg config.Info
	if err := cfg.ReadFile(*flagConf); err != nil {
		zap.L().Fatal("unable to read config",
			zap.String("filename", *flagConf),
			zap.Error(err))
	}

	if *flagToServers {
		for name, build := range cfg.Servers {
			if err := server.Deploy(root, name, build); err != nil {
				zap.L().Fatal("unable to deploy server",
					zap.String("name", name))
			}
		}
	}

	if *flagToClients {
		for name, build := range cfg.Clients {
			if err := client.Deploy(root, name, build); err != nil {
				zap.L().Fatal("unable to deploy client",
					zap.String("name", name))
			}
		}
	}
}
