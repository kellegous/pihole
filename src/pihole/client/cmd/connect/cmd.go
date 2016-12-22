package connect

import (
	"time"

	"pihole/client/config"
	"pihole/client/proxy"

	"github.com/golang/glog"
)

// Name ...
const Name = "connect"

// Main ...
func Main(conf string, args []string) {
	var cfg config.Config
	if err := cfg.Read(conf); err != nil {
		glog.Fatal(err)
	}

	for {
		if err := proxy.ConnectAndServe(&cfg); err != nil {
			glog.Error(err)
		}

		glog.Infoln("reconnecting in 10 secs...")
		time.Sleep(10 * time.Second)
	}
}
