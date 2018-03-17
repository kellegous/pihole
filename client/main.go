package main

import (
	"flag"

	"github.com/kellegous/pihole/client/cmd/connect"
	"github.com/kellegous/pihole/client/cmd/createconfig"
	"github.com/kellegous/pihole/client/cmd/getkey"
	"github.com/kellegous/pihole/client/cmd/version"
	"github.com/kellegous/pihole/logging"
)

func main() {
	flagConf := flag.String(
		"conf",
		"pihole",
		"")
	flag.Parse()

	logging.MustSetup()

	args := flag.Args()

	if flag.NArg() == 0 {
		connect.Main(*flagConf, args)
		return
	}

	switch args[0] {
	case connect.Name:
		connect.Main(*flagConf, args[1:])
	case createconfig.Name:
		createconfig.Main(*flagConf, args[1:])
	case getkey.Name:
		getkey.Main(*flagConf, args[1:])
	case version.Name:
		version.Main(*flagConf, args[1:])
	}
}
