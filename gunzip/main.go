package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func main() {
	elog := log.New(os.Stderr, "gunzip: ", 0)
	inbuf := bufio.NewReader(os.Stdin)
	r, err := gzip.NewReader(inbuf)
	if err != nil {
		elog.Fatal(err)
	}
	defer r.Close()
	w := bufio.NewWriter(os.Stdout)
	_, err = io.Copy(w, r)
	if err != nil {
		elog.Fatal(err)
	}
	err = w.Flush()
	if err != nil {
		elog.Fatal(err)
	}
}
