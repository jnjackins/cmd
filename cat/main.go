package main

import (
	"bufio"
	"io"
	"os"

	"sigint.ca/die"
)

func main() {
	outbuf := bufio.NewWriter(os.Stdout)
	defer outbuf.Flush()
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "/dev/stdin")
	}
	for _, fname := range os.Args[1:] {
		f, err := os.Open(fname)
		die.On(err, "cat: error opening file")
		fbuf := bufio.NewReader(f)
		io.Copy(outbuf, fbuf)
		f.Close()
	}
}
