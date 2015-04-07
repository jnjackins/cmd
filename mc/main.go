package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"

	"sigint.ca/text/column"
)

var cflag = flag.Int("c", 80, "maximum columns per line")

func main() {
	elog := log.New(os.Stderr, "mc: ", 0)
	flag.Parse()
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()
	colWriter := column.NewWriter(out, *cflag)
	_, err := io.Copy(colWriter, in)
	if err != nil {
		elog.Print(err)
	}
	if err := colWriter.Flush(); err != nil {
		elog.Print(err)
	}
}
