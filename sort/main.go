package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

var nflag = flag.Bool("n", false, "sort numerically")

type numericStrings []string

func main() {
	elog := log.New(os.Stderr, "sort: ", 0)
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		for _, path := range args {
			f, err := os.Open(path)
			if err != nil {
				elog.Print(err)
			}
			if err := sortReader(f, *nflag); err != nil {
				elog.Print(err)
			}
		}
	} else {
		if err := sortReader(os.Stdin, *nflag); err != nil {
			elog.Print(err)
		}
	}
}

func sortReader(r io.Reader, numeric bool) error {
	text, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	lines := strings.Split(string(text), "\n")
	if numeric {
		sort.Sort(numericStrings(lines))
	} else {
		sort.Strings(lines)
	}
	for _, l := range lines {
		fmt.Println(l)
	}
	return nil
}

func (s numericStrings) Len() int      { return len(s) }
func (s numericStrings) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s numericStrings) Less(i, j int) bool {
	var num1, num2 float64
	n, _ := fmt.Sscanf(s[i], "%f%_\n", &num1)
	m, _ := fmt.Sscanf(s[j], "%f%_\n", &num2)
	if n == 1 && m == 1 {
		if num1 == num2 {
			return s[i] < s[j]
		}
		return num1 < num2
	}
	return s[i] < s[j]
}
