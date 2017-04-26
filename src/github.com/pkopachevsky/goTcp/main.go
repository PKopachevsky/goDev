package main

import (
	"flag"
	"fmt"
	"os"
	"log"
)

var (
	remoteHost 	= flag.String("host", "", "Remote host")
	rempotePort 	= flag.Int("port", 0, "Remote port")
	listen 		= flag.String("listen", ":4242", "Local address to listen")
	dump 		= flag.String("dump", "", "Write dump to file")
	skipHealthcheck = flag.Bool("skip-healthcheck", false, "Skip heathcheck")
)

func main()  {
	flag.Parse()
	remoteAddr := fmt.Sprintf("%s:%d", *remoteHost, *rempotePort)
	proxy := &proxyServer{localAddr: *listen, remoteAddr:remoteAddr, dumpTo: dumpTo(*dump)}

	var err error
	if !*skipHealthcheck {
		err = proxy.healthCheck()
		if err !=nil {
			log.Fatal(err)
		}
		log.Println("Healthcheck to %s OK", remoteAddr)
	}
	err = proxy.start()
	if err != nil {
		log.Fatal(err)
	}
}

func dumpTo(filename string) *os.File {
	dumpTo := os.Stdout;
	if len(filename) > 0 {
		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Fail to open file %s, fallback to stdout", filename)
		} else {
			dumpTo = file
		}
	}
	return dumpTo
}