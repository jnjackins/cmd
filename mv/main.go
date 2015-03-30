package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

var usage = `Usage: mv source target
       mv file ... directory`

func main() {
	elog := log.New(os.Stderr, "mv: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	if len(args) == 2 {
		from, to := args[0], args[1]
		stat, err := os.Stat(to)
		if err != nil && !os.IsNotExist(err) {
			elog.Fatal(err)
		}
		if stat != nil && stat.IsDir() {
			to = to + "/" + path.Base(from)
		}
		if err := os.Rename(from, to); err != nil {
			elog.Fatal(err)
		}
	} else {
		dir := args[len(args)-1]
		stat, err := os.Stat(dir)
		if err != nil {
			elog.Fatal(err)
		}
		if !stat.IsDir() {
			elog.Fatal("multi-file target must be directory\n" + usage)
		}
		for _, fname := range args[:len(args)-1] {
			err := os.Rename(fname, dir+"/"+path.Base(fname))
			if err != nil {
				elog.Print(err)
			}
		}
	}
}
