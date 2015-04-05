package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var pagesize = flag.Int("l", 22, "`lines` per page")

var (
	stdin, tty *bufio.Reader
	elog       log.Logger
)

func main() {
	elog := log.New(os.Stderr, "p: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: p [options] [file ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	stdin = bufio.NewReader(os.Stdin)
	f, err := os.Open("/dev/tty")
	if err != nil {
		f, err = os.Open("/dev/console")
		if err != nil {
			elog.Fatal("couldn't open /dev/tty or /dev/console")
		}
	}
	defer f.Close()
	tty = bufio.NewReader(f)
	if len(flag.Args()) == 0 {
		page(stdin)
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			elog.Print("error opening \"%s\"; skipping", path)
		}
		page(f)
		f.Close()
	}
}

func page(r io.Reader) {
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		fmt.Print(scanner.Text())
		i++
		if i == *pagesize {
			tty.ReadString('\n')
			i = 0
		} else {
			fmt.Println()
		}
	}
	if err := scanner.Err(); err != nil {
		elog.Print(err)
	}
}
