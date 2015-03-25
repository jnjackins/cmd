package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"

	"sigint.ca/die"
)

func main() {
	inbuf := bufio.NewReader(os.Stdin)
	r, err := gzip.NewReader(inbuf)
	die.On(err, "gunzip: error creating gzip reader from stdin")
	defer r.Close()
	w := bufio.NewWriter(os.Stdout)
	_, err = io.Copy(w, r)
	die.On(err, "gunzip: error copying decompressed data to stdout")
	err = w.Flush()
	die.On(err, "gunzip: error flushing stdout")
}
