//go:generate go tool yacc -p "hoc" hoc.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var elog *log.Logger

var eflag = flag.String("e", "", "Evaluate `expression` instead of reading from file or stdin.")

func main() {
	elog = log.New(os.Stderr, "hoc: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: hoc [ file ... ] [ -e expression ]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *eflag != "" {
		parse(strings.NewReader(*eflag))
	} else {
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
