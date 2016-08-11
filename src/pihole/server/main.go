package main

import (
	"flag"
	"log"
	"net/http"

	"pihole/api"
	"pihole/hub"
)

func main() {
	flagAddr := flag.String("addr", ":http", "")
	flag.Parse()
	h := hub.NewHub()

	go func() {
		log.Panic(api.ListenAndServe(api.DefaultAddr, h))
	}()

	log.Panic(http.ListenAndServe(*flagAddr, h))
}
