package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	aflag = flag.Bool("a", false, "Print entries for all files.")
	sflag = flag.Bool("s", false, "Print only a final summary.")
)

var (
	elog *log.Logger
	out  *bufio.Writer
)

func init() {
	elog = log.New(os.Stderr, "du: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: du [options] [path ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	out = bufio.NewWriter(os.Stdout)
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, path := range args {
		du(path)
	}
	out.Flush()
}

func du(path string) {
	total := walk(path)
	if *sflag {
		fmt.Println(total)
	}
}

func walk(path string) int64 {
	fi, err := os.Lstat(path)
	if err != nil {
		elog.Print(err)
		return 0
	}
	var size int64
	if fi.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			elog.Print(err)
			return 0
		}
		entries, err := f.Readdirnames(0)
		if err != nil {
			elog.Print(err)
			return 0
		}
		f.Close()
		for _, entry := range entries {
			size += walk(path + "/" + entry)
		}
	} else {
		size = (fi.Size() + 1023) / 1024
	}
	if !*sflag && (fi.IsDir() || *aflag) {
		print(path, size)
	}
	return size
}

func print(path string, size int64) {
	_, err := fmt.Fprintf(out, "%d\t%s\n", size, filepath.Clean(path))
	if err != nil {
		elog.Print(err)
	}
}
