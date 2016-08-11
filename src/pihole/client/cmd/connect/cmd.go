package connect

import (
	"flag"
	"log"
	"time"

	"pihole/client/config"
	"pihole/client/proxy"
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
		log.Panic(err)
	}

	for {
		if err := proxy.ConnectAndServe(&cfg); err != nil {
			log.Println(err)
		}

		time.Sleep(10 * time.Second)
	}
}
