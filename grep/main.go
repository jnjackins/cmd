package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var cflag = flag.Bool("c", false, "Print only a count of matching lines.")
var iflag = flag.Bool("i", false, "Ignore alphabetic case distinctions.")
var lflag = flag.Bool("l", false, "Print only the names of files with selected lines.")
var Lflag = flag.Bool("L", false, "Print only the names of files with no selected lines.")
var nflag = flag.Bool("n", false, "Give line number for each matching line.")
var qflag = flag.Bool("q", false, "Same as s. Provided for (poor) compatibility.")
var sflag = flag.Bool("s", false, "Produce no output, but return status.")
var vflag = flag.Bool("v", false, "Print lines that do not match the pattern.")

func main() {
	elog := log.New(os.Stderr, "grep: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: grep [options] pattern [file ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	s := flag.Arg(0)
	if *iflag {
		s = strings.ToLower(s)
	}
	pattern, err := regexp.CompilePOSIX(s)
	if err != nil {
		elog.Fatal(err)
	}
	matches := 0
	files := flag.Args()[1:]
	if len(files) == 0 {
		files = append(files, "/dev/stdin")
	}
	multi := len(files) > 1
	for i := range files {
		fmatches, err := grep(pattern, files[i], multi)
		if err != nil {
			elog.Printf("error scanning %s: %s", files[i], err)
		}
		if (*lflag || *Lflag) && fmatches > 0 {
			fmt.Println(files[i])
			matches++
			continue
		}
		if *cflag {
			if multi {
				fmt.Print(files[i] + ":")
			}
			fmt.Println(fmatches)
		}
		matches += fmatches
	}
	if matches == 0 {
		os.Exit(1)
	}
}

func grep(pattern *regexp.Regexp, fname string, multi bool) (int, error) {
	f, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(f)
	n, matches := 0, 0
	for scanner.Scan() {
		n++
		line := scanner.Text()
		match := false
		if *iflag {
			match = pattern.MatchString(strings.ToLower(line))
		} else {
			match = pattern.MatchString(line)
		}
		if !*vflag && match || *vflag && !match {
			matches++
			if *lflag {
				return 1, scanner.Err()
			}
			if *Lflag {
				return 0, scanner.Err()
			}
			if *sflag || *qflag {
				os.Exit(0)
			}
			if *cflag {
				continue
			}
			if multi {
				fmt.Print(fname + ":")
			}
			if *nflag {
				fmt.Printf("%d:", n)
			}
			fmt.Println(line)
		}
	}
	if *Lflag {
		return 1, scanner.Err()
	}
	return matches, scanner.Err()
}
