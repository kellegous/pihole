package main

import (
	"flag"
	"pihole/logging"

	"pihole/client/cmd/connect"
	"pihole/client/cmd/createconfig"
	"pihole/client/cmd/getkey"
)

func main() {
	flagConf := flag.String(
		"conf",
		"pihole",
		"")
	flag.Parse()

	if err := logging.Setup(); err != nil {
		panic(err)
	}

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
	}
}
