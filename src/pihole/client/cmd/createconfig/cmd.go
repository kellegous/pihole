package createconfig

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"pihole/client/config"

	"github.com/golang/glog"
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
func Main(args []string) {
	// create-config pihole.com localhost kpi.kellegous.com kpi.kellego.us
	f := flag.NewFlagSet(Name, flag.ExitOnError)
	flagConf := f.String("conf", "pihole", "")
	f.Parse(args)

	if f.NArg() < 3 {
		usage()
	}

	var cfg config.Config
	if err := create(&cfg, args[0], args[1], args[2:]); err != nil {
		glog.Fatal(err)
	}

	if err := cfg.Create(*flagConf); err != nil {
		glog.Fatal(err)
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println(pub)
}
