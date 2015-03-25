// TODO other types
// TODO file properties
// TODO handle dot

package main

import (
	"archive/tar"
	"bufio"
	"io"
	"os"

	"sigint.ca/die"
)

func main() {
	inbuf := bufio.NewReader(os.Stdin)
	r := tar.NewReader(inbuf)

	header, err := r.Next()
	for err == nil {
		info := header.FileInfo()
		if info.IsDir() {
			innerErr := os.Mkdir(header.Name, info.Mode())
			die.On(innerErr, "untar: error creating directory")
		} else {
			f, innerErr := os.Create(header.Name)
			die.On(innerErr, "untar: error creating file")

			fbuf := bufio.NewWriter(f)
			_, innerErr = io.Copy(fbuf, r)
			die.On(innerErr, "untar: error copying data")
			fbuf.Flush()
			f.Close()
		}
		header, err = r.Next()
	}
	if err != io.EOF {
		die.On(err, "untar: error advancing to next entry")
	}
}
