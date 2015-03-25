// TODO other types
// TODO file properties
// TODO handle dot

package main

import (
	"archive/tar"
	"bufio"
	"io"
	"os"
	"path/filepath"

	"sigint.ca/die"
)

func main() {
	outbuf := bufio.NewWriter(os.Stdout)
	defer outbuf.Flush()
	w := tar.NewWriter(outbuf)
	defer w.Close()
	for _, arg := range os.Args[1:] {
		walkfunc := func(path string, info os.FileInfo, err error) error {
			f, err := os.Open(path)
			die.On(err, "tar: error opening file")
			defer f.Close()

			info, err = f.Stat()
			die.On(err, "tar: error getting file info")

			header, err := tar.FileInfoHeader(info, "")
			die.On(err, "tar: error making header")
			header.Name = path

			err = w.WriteHeader(header)
			die.On(err, "tar: error writing header")

			if info.Mode().IsRegular() {
				fbuf := bufio.NewReader(f)
				_, err = io.Copy(w, fbuf)
				die.On(err, "tar: error copying file data to stdout")
			}

			return nil
		}
		err := filepath.Walk(arg, walkfunc)
		die.On(err, "tar: error walking tree at "+arg)
	}
}
