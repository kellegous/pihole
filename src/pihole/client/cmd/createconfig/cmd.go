package createconfig

import (
	"fmt"
	"net/url"
	"os"

	"go.uber.org/zap"

	"pihole/client/config"
)

const bitsIntKey = 4096

// Name ...
const Name = "create-config"

func usage() {
	fmt.Fprintf(os.Stderr,
		"usage: %s %s hub to-url host...\n",
		os.Args[0],
		Name)
	os.Exit(1)
}

func create(cfg *config.Config, addr, toURL string, hosts []string) error {
	u, err := url.Parse(toURL)
	if err != nil {
		return err
	}

	cfg.Hosts = hosts
	cfg.ToURL = u
	cfg.Hub.Addr = addr
	cfg.Hub.User = "pihole"

	return nil
}

// Main ...
func Main(conf string, args []string) {
	// create-config pihole.com localhost kpi.kellegous.com kpi.kellego.us
	if len(args) < 3 {
		usage()
	}

	var cfg config.Config
	if err := create(&cfg, args[0], args[1], args[2:]); err != nil {
		zap.L().Fatal("unable to create config",
			zap.Error(err))
	}

	if err := cfg.Create(conf); err != nil {
		zap.L().Fatal("unable to create config",
			zap.Error(err))
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		zap.L().Fatal("unable to load public key",
			zap.Error(err))
	}

	fmt.Println(pub)
}
