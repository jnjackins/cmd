package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	rflag  = flag.Bool("r", false, "Recursive mode.")
	fflag  = flag.Bool("f", false, "Force mode.")
	rfflag = flag.Bool("rf", false, "-r and -f, for compatibility.")
	iflag  = flag.Bool("i", false, "Interactive mode.")
)

func main() {
	log.SetPrefix("rm: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: rm [options] file ...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *rfflag {
		*rflag = true
		*fflag = true
	}
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	var status int
	for _, path := range flag.Args() {
		if *iflag {
			fmt.Printf("remove %s? ", path)
			s, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			if !strings.HasPrefix(strings.ToLower(s), "y") {
				continue
			}
		}
		if *rflag {
			err := os.RemoveAll(path)
			if err != nil && !*fflag {
				log.Print(err)
				status++
			}
		} else {
			err := os.Remove(path)
			if err != nil && !*fflag {
				log.Print(err)
				status++
			}
		}
	}
	os.Exit(status)
}
