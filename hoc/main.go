//go:generate go tool yacc -p "hoc" hoc.y

package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

var elog *log.Logger

func main() {
	elog = log.New(os.Stderr, "hoc: ", 0)
	flag.Parse()
	f := os.Stdin
	args := flag.Args()
	if len(args) > 0 {
		var err error
		f, err = os.Open(args[0])
		if err != nil {
			elog.Fatal(err)
		}
	}
	parse(f)
}

func parse(r io.Reader) {
	in := bufio.NewScanner(bufio.NewReader(r))
	var sr *strings.Reader
	for in.Scan() {
		sr = strings.NewReader(in.Text() + "\n")
		hocParse(&hocLex{r: sr})
	}
	if err := in.Err(); err != nil {
		elog.Print(err)
	}
}
