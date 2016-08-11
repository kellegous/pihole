package main

import (
	"flag"
	"log"
)

func main() {
	flagAddr := flag.String("addr", ":80", "")
	flagCert := flag.String("cert", "", "")
	flag.Parse()

	log.Printf("addr = %s", *flagAddr)
	log.Printf("cert = %s", *flagCert)
}
