package main

import (
	"flag"
	"net/http"

	"go.uber.org/zap"

	"pihole/api"
	"pihole/hub"
	"pihole/logging"
)

func main() {
	flagAddr := flag.String("addr", ":http", "")
	flag.Parse()

	logging.MustSetup()

	h := hub.NewHub()

	go func() {
		if err := api.ListenAndServe(api.DefaultAddr, h); err != nil {
			zap.L().Fatal("unable to listen and serve for api",
				zap.Error(err))
		}
	}()

	zap.L().Info("web started",
		zap.String("address", *flagAddr))

	if err := http.ListenAndServe(*flagAddr, logging.WithLog(h)); err != nil {
		zap.L().Fatal("unable to listen and serve for web",
			zap.Error(err))
	}
}
