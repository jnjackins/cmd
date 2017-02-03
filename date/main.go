package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	stamp  = flag.Bool("stamp", false, "Predefined format: "+time.Stamp)
	rfc339 = flag.Bool("rfc3339", false, "Predefined format: "+time.RFC3339)
	custom = flag.String("format", "", "Provide a custom Go time format (ref: Mon Jan 2 3:04:05PM 2006 -0700)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [format option]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nSee https://golang.org/pkg/time/#pkg-constants for formatting details.")
	}
	flag.Parse()
	format := time.UnixDate
	if *stamp {
		format = time.Stamp
	}
	if *rfc339 {
		format = time.RFC3339
	}
	if *custom != "" {
		format = *custom
	}
	fmt.Println(time.Now().Format(format))
}
