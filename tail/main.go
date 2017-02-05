package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var fflag = flag.Bool("f", false, "After printing the tail as usual, follow additions to the file and print them.")
var nflag = flag.Int("n", 10, "Set the number of `lines` at the end of file to be printed.")

func main() {
	log.SetPrefix("tailf: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: tail [options] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}

	path := "/dev/stdin"
	if flag.NArg() == 1 {
		path = flag.Arg(1)
	}
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	// try seek to the last byte of the file
	if pos, err := f.Seek(-1, 2); err != nil || pos < 0 {
		// couldn't seek, read the whole thing and traverse backwards
		buf, err := getbuf(f, *nflag)
		if err != nil {
			log.Fatalf("getbuf: %v", err)
		}
		os.Stdout.Write(buf)
	} else {
		if err := backtrack(f, *nflag); err != nil {
			log.Fatalf("backtrack: %v", err)
		}
		buf := make([]byte, 8192)
		io.CopyBuffer(os.Stdout, f, buf)
		if *fflag {
			for range time.NewTicker(500 * time.Millisecond).C {
				io.CopyBuffer(os.Stdout, f, buf)
			}
		}
	}
}

func getbuf(f *os.File, n int) ([]byte, error) {
	buf, err := ioutil.ReadAll(f)
	if err != nil && err != io.EOF {
		return nil, err
	}
	sep := []byte("\n")
	p := buf
	i := len(p) - 1
	for lines := 0; lines <= *nflag; lines++ {
		i = bytes.LastIndex(p, sep)
		if i < 0 {
			break
		}
		p = p[:i]
	}
	i += 1 // don't need the leading newline
	return buf[i:], nil
}

// Count n lines backwards from the end of the file.
func backtrack(f *os.File, n int) error {
	buf := make([]byte, 1)

	// stop only after encountering the *nflag+1th newline, so that the line
	// containing the 10th encountered newline is also printed.
	remaining := n + 1

	var pos int64
	for {
		n, err := f.Read(buf)
		if err != nil {
			return fmt.Errorf("read byte: %v", err)
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
			return nil
		}
		if pos == 0 {
			break
		}
	}
	return nil
}
