package main

import (
	"flag"
	"pihole/logging"

	"pihole/client/cmd/connect"
	"pihole/client/cmd/createconfig"
	"pihole/client/cmd/getkey"
	"pihole/client/cmd/version"
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
