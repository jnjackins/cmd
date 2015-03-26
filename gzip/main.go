package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func main() {
	elog := log.New(os.Stderr, "gzip: ", 0)
	r := bufio.NewReader(os.Stdin)
	outbuf := bufio.NewWriter(os.Stdout)
	w := gzip.NewWriter(outbuf)
	_, err := io.Copy(w, r)
	if err != nil {
		elog.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		elog.Fatal(err)
	}
	err = outbuf.Flush()
	if err != nil {
		elog.Fatal(err)
	}
}
