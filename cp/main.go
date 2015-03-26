package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
)

var usage = `Usage: cp source target
       cp file ... directory`

func main() {
	elog := log.New(os.Stderr, "cp: ", 0)
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}
	if len(args) == 2 {
		stat, err := os.Stat(args[1])
		if err != nil {
			elog.Fatal(err)
		}
		if stat.IsDir() {
			err := cp(args[0], args[1]+"/"+path.Base(args[0]))
			if err != nil {
				elog.Fatal(err)
			}
		} else {
			err := cp(args[0], args[1])
			if err != nil {
				elog.Fatal(err)
			}
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
		var wg sync.WaitGroup
		for _, fname := range args[:len(args)-1] {
			wg.Add(1)
			// copy files concurrently
			go func(fname string) {
				defer wg.Done()
				err := cp(fname, dir+"/"+path.Base(fname))
				if err != nil {
					elog.Print(err)
				}
			}(fname)
		}
		wg.Wait()
	}
}

func cp(from, to string) error {
	stat, err := os.Stat(from)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.New(from + " is a directory (not copied)")
	}
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()
	dest, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dest.Close()
	r, w := bufio.NewReader(source), bufio.NewWriter(dest)
	_, err = io.Copy(w, r)
	return err
}
