package connect

import (
	"flag"
	"time"

	"pihole/client/config"
	"pihole/client/proxy"

	"github.com/golang/glog"
)

// Name ...
const Name = "connect"

// Main ...
func Main(args []string) {
	f := flag.NewFlagSet(Name, flag.PanicOnError)
	flagConfg := f.String(
		"conf",
		"pihole",
		"")
	f.Parse(args)

	var cfg config.Config
	if err := cfg.Read(*flagConfg); err != nil {
		glog.Fatal(err)
	}

	for {
		if err := proxy.ConnectAndServe(&cfg); err != nil {
			glog.Error(err)
		}

		time.Sleep(10 * time.Second)
	}
}
