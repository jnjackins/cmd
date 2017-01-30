package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var fflag = flag.Bool("f", false, "After printing the tail as usual, follow additions to the file and print them.")
var nflag = flag.Int("n", 10, "Set the number of `lines` at the end of file to be printed.")

func init() {
	log.SetPrefix("tailf: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: tail [options] file")
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	seek(f, *nflag)
	buf := make([]byte, 8192)
	tcopy(f, buf)
	if *fflag {
		for range time.NewTicker(500 * time.Millisecond).C {
			tcopy(f, buf)
		}
	}
}

// Count n lines backwards from the end of the file.
func seek(f *os.File, n int) {
	// seek to just before the last byte of the file, so we can compare it to '\n'.
	if _, err := f.Seek(-1, 2); err != nil {
		return
	}
	buf := make([]byte, 1)

	// stop only after encountering the *nflag+1th newline, so that the line
	// containing the 10th encountered newline is also printed.
	remaining := n + 1

	var pos int64
	for {
		n, err := f.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		pos += int64(n)
		if buf[0] == '\n' {
			remaining--
		}
		if remaining == 0 {
			break
		}
		// seek back 2 bytes: one for what we read, one for backwards progress
		pos, err = f.Seek(-2, 1)
		if err != nil {
			// if that fails, at least try to seek back what we just read
			f.Seek(-1, 1)
			return
		}
		if pos == 0 {
			break
		}
	}
}

func tcopy(f *os.File, buf []byte) {
	_, err := io.CopyBuffer(os.Stdout, f, buf)
	if err != nil {
		log.Fatal(err)
	}
}
