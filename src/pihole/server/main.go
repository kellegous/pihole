package main

import (
	"flag"
	"net/http"

	"pihole/api"
	"pihole/hub"
	"pihole/logging"

	"github.com/golang/glog"
)

func main() {
	flagAddr := flag.String("addr", ":http", "")
	flag.Parse()
	h := hub.NewHub()

	go func() {
		glog.Fatal(api.ListenAndServe(api.DefaultAddr, h))
	}()

	glog.Infof("Web: %s", *flagAddr)
	glog.Fatal(http.ListenAndServe(
		*flagAddr,
		logging.WithLog(h)))
}
