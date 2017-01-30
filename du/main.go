package main

import (
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

func init() {
	log.SetPrefix("du: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: du [options] [path ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, path := range args {
		du(path)
	}
}

func du(path string) {
	total := walk(path)
	if *sflag {
		print(path, total)
	}
}

func walk(path string) int64 {
	fi, err := os.Lstat(path)
	if err != nil {
		log.Print(err)
		return 0
	}
	var size int64
	if fi.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			log.Print(err)
			return 0
		}
		entries, err := f.Readdirnames(0)
		if err != nil {
			log.Print(err)
			return 0
		}
		f.Close()
		for _, entry := range entries {
			size += walk(path + "/" + entry)
		}
	} else {
		size = fi.Size()
	}
	if !*sflag && (fi.IsDir() || *aflag) {
		print(path, size)
	}
	return size
}

func print(path string, size int64) {
	fmt.Fprintf(os.Stdout, "%d\t%s\n", (size+512)/1024, filepath.Clean(path))
}
