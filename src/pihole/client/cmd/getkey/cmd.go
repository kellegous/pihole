package getkey

import (
	"flag"
	"fmt"

	"pihole/client/config"

	"github.com/golang/glog"
)

// Name ...
const Name = "get-key"

// Main ...
func Main(args []string) {
	f := flag.NewFlagSet(Name, flag.ExitOnError)
	flagConf := f.String("conf", "pihole", "Path to config")
	f.Parse(args)

	var cfg config.Config
	if err := cfg.Read(*flagConf); err != nil {
		glog.Fatal(err)
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println(pub)
}
