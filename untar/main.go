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
)

func main() {
	elog := log.New(os.Stderr, "untar: ", 0)
	inbuf := bufio.NewReader(os.Stdin)
	r := tar.NewReader(inbuf)
	header, err := r.Next()
	for err == nil {
		info := header.FileInfo()
		if info.IsDir() {
			if err := os.Mkdir(header.Name, info.Mode()); err != nil {
				elog.Fatal(err)
			}
		} else {
			f, err := os.Create(header.Name)
			if err != nil {
				elog.Fatal(err)
			}

			fbuf := bufio.NewWriter(f)
			if _, err = io.Copy(fbuf, r); err != nil {
				log.Fatal(err)
			}
			fbuf.Flush()
			f.Close()
		}
		header, err = r.Next()
	}
	if err != io.EOF {
		log.Fatal(err)
	}
}
