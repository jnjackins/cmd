// TODO other types
// TODO file properties
// TODO handle dot

package main

import (
	"archive/tar"
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	elog := log.New(os.Stderr, "tar: ", 0)
	outbuf := bufio.NewWriter(os.Stdout)
	defer outbuf.Flush()
	w := tar.NewWriter(outbuf)
	defer w.Close()
	for _, arg := range os.Args[1:] {
		walkfunc := func(path string, info os.FileInfo, err error) error {
			f, err := os.Open(path)
			if err != nil {
				elog.Fatal(err)
			}
			defer f.Close()
			info, err = f.Stat()
			if err != nil {
				elog.Fatal(err)
			}
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				elog.Fatal(err)
			}
			header.Name = path
			if err := w.WriteHeader(header); err != nil {
				elog.Fatal(err)
			}
			if info.Mode().IsRegular() {
				fbuf := bufio.NewReader(f)
				if _, err = io.Copy(w, fbuf); err != nil {
					elog.Fatal(err)
				}
			}
			return nil
		}
		if err := filepath.Walk(arg, walkfunc); err != nil {
			elog.Fatal(err)
		}
	}
}
