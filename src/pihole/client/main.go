package main

import (
	"os"
	"strings"

	"pihole/client/cmd/connect"
	"pihole/client/cmd/createconfig"
	"pihole/client/cmd/getkey"
)

func main() {
	args := os.Args

	if len(args) == 1 || strings.HasPrefix(os.Args[1], "-") {
		connect.Main(args[1:])
		return
	}

	switch args[1] {
	case connect.Name:
		connect.Main(args[2:])
	case createconfig.Name:
		createconfig.Main(args[2:])
	case getkey.Name:
		getkey.Main(args[2:])
	}
}
