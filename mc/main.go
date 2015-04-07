package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"

	"sigint.ca/text/column"
)

var wflag = flag.Int("w", 80, "maximum `width` of columnated text")

func main() {
	elog := log.New(os.Stderr, "mc: ", 0)
	flag.Parse()
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()
	colWriter := column.NewWriter(out, *wflag)
	_, err := io.Copy(colWriter, in)
	if err != nil {
		elog.Print(err)
	}
	if err := colWriter.Flush(); err != nil {
		elog.Print(err)
	}
}
