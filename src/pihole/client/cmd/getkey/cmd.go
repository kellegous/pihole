package getkey

import (
	"flag"
	"fmt"
	"log"

	"pihole/client/config"
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
		log.Panic(err)
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(pub)
}
