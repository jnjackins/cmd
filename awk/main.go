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

var (
	debug  bool     // debug mode
	lineno int      // line number in awk program
	args   []string // command line arguments
)

var (
	dflag = flag.Bool("d", false, "Debug mode")
	fflag = flag.String("f", "", "Specifies a program `file`")
	Fflag = flag.String("F", "", "Specifies a field `separator`")
	vflag = flag.String("v", "", "Initializes `var=value`")
)

const usage = "Usage: %s [-F fieldsep] [-v var=value] [-f programfile | 'program'] [file ...]\n"

func init() {
	cmdname := os.Args[0]
	log.SetOutput(os.Stderr)
	log.SetPrefix(cmdname + ": ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, cmdname)
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	args = flag.Args()
	if len(args) == 0 && *fflag == "" {
		flag.Usage()
		os.Exit(1)
	}
	debug = *dflag

	initSymbols()

	r := progReader()
	prog, err := compileProg(r)
	r.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rec := splitRecord(scanner.Text())
		updateSymbols(rec)
		prog.exec(rec)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func progReader() io.ReadCloser {
	if *fflag != "" {
		f, err := os.Open(*fflag)
		if err != nil {
			log.Fatal(err)
		}
		return f
	} else {
		return nopCloser{strings.NewReader(args[0])}
	}
}

type nopCloser struct{ io.Reader }

func (nopCloser) Close() error { return nil }

func dprintf(format string, args ...interface{}) {
	if debug {
		log.Printf("[ "+format+" ]", args...)
	}
}
