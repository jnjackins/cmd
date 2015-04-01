// TODO other types
// TODO file properties
// TODO handle dot

package main

import (
	"archive/tar"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var usage = `Usage:
  tar c path ...	- [c]reate tarball and write to stdout
  tar x			- e[x]tract tarball from stdin`

func main() {
	elog := log.New(os.Stderr, "tar: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	switch args[0] {
	case "c":
		paths := args[1:]
		if len(paths) == 0 {
			flag.Usage()
			os.Exit(1)
		}
		if err := create(paths); err != nil {
			elog.Fatal(err)
		}
	case "x":
		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		if err := extract(); err != nil {
			elog.Fatal(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func create(paths []string) error {
	outbuf := bufio.NewWriter(os.Stdout)
	defer outbuf.Flush()
	w := tar.NewWriter(outbuf)
	defer w.Close()

	walk := func(path string, info os.FileInfo, err error) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		info, err = f.Stat()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = path
		if err := w.WriteHeader(header); err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			fbuf := bufio.NewReader(f)
			if _, err = io.Copy(w, fbuf); err != nil {
				return err
			}
		}
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			return err
		}
	}
	return nil
}

func extract() error {
	inbuf := bufio.NewReader(os.Stdin)
	r := tar.NewReader(inbuf)
	header, err := r.Next()
	for err == nil {
		info := header.FileInfo()
		if info.IsDir() {
			if err := os.Mkdir(header.Name, info.Mode()); err != nil {
				return err
			}
		} else {
			f, err := os.Create(header.Name)
			if err != nil {
				return err
			}

			fbuf := bufio.NewWriter(f)
			if _, err = io.Copy(fbuf, r); err != nil {
				return err
			}
			fbuf.Flush()
			f.Close()
		}
		header, err = r.Next()
	}
	if err != io.EOF {
		return err
	}
	return nil
}
