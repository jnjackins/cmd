package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	network = flag.String("network", "tcp", "specify network")
	timeout = flag.Duration("timeout", 5*time.Second, "specify a timeout duration")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] hostname service/port ...\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 2 {
		usage()
		os.Exit(1)
	}

	hostname := flag.Arg(0)
	failed := 0
	for _, service := range flag.Args()[1:] {
		port, err := net.LookupPort(*network, service)
		if err != nil {
			fmt.Printf("unknown service: %s\n", service)
			failed++
			continue
		}
		fmt.Printf("connecting to %s on port %d... ", hostname, port)

		address := net.JoinHostPort(hostname, service)
		conn, err := net.DialTimeout(*network, address, *timeout)
		if err != nil {
			failed++
			fmt.Printf("failed: %v\n", err)
			continue
		}
		fmt.Printf("ok\n")
		conn.Close()
	}
	os.Exit(failed)
}
