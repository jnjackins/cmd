package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var usage = `Usage: ln [-s] source [target]
       ln [-s] file ... dir`

var sflag = flag.Bool("s", false, "Create a symbolic link.")

func main() {
	elog := log.New(os.Stderr, "ln: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
	flag.Parse()
	link := os.Link
	if *sflag {
		link = os.Symlink
	}
	args := flag.Args()
	switch len(args) {
	case 0:
		flag.Usage()
		os.Exit(1)
	case 1:
		if err := link(args[0], filepath.Base(args[0])); err != nil {
			elog.Fatal(err)
		}
	default:
		last := args[len(args)-1]
		info, err := os.Stat(last)
		if err != nil {
			if !os.IsNotExist(err) {
				elog.Fatal(err)
			}
		}
		if info != nil && info.IsDir() {
			var fail bool
			for _, arg := range args[:len(args)-1] {
				if err := link(arg, last+"/"+filepath.Base(arg)); err != nil {
					fail = true
					elog.Print(err)
				}
			}
			if fail {
				os.Exit(1)
			}
		} else {
			if len(args) != 2 {
				flag.Usage()
				os.Exit(1)
			}
			if err := link(args[0], args[1]); err != nil {
				elog.Fatal(err)
			}
		}
	}
}
