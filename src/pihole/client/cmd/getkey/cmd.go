package getkey

import (
	"fmt"

	"pihole/client/config"

	"github.com/golang/glog"
)

// Name ...
const Name = "get-key"

// Main ...
func Main(conf string, args []string) {
	var cfg config.Config
	if err := cfg.Read(conf); err != nil {
		glog.Fatal(err)
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println(pub)
}
