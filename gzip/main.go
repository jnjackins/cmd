package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"

	"sigint.ca/die"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	outbuf := bufio.NewWriter(os.Stdout)
	w := gzip.NewWriter(outbuf)
	_, err := io.Copy(w, r)
	die.On(err, "gzip: error copying compressed data to stdout")
	err = w.Close()
	die.On(err, "gzip: error closing gzip writer")
	err = outbuf.Flush()
	die.On(err, "gzip: error flushing stdout")
}
