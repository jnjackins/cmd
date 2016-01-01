package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var fflag = flag.String("f", "", "List of `fields` to cut.")

func main() {
	elog := log.New(os.Stderr, "cut: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cut -f list [file ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *fflag == "" {
		elog.Print("illegal list value")
		flag.Usage()
		os.Exit(1)
	}
	fields := make([]int, 0, 1)
	for _, r := range strings.Split(*fflag, ",") {
		dash := strings.IndexRune(r, '-')
		if dash > -1 {
			low, err1 := strconv.Atoi(r[:dash])
			high, err2 := strconv.Atoi(r[dash+1:])
			if err1 != nil || err2 != nil {
				elog.Fatalf("bad range: %s", r)
			}
			for n := low; n <= high; n++ {
				fields = append(fields, n)
			}
		} else {
			n, err := strconv.Atoi(r)
			if err != nil {
				elog.Fatalf("bad field: %s", r)
			}
			fields = append(fields, n)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out := cut(scanner.Text(), fields)
		fmt.Println(strings.Join(out, " "))
	}
	if err := scanner.Err(); err != nil {
		elog.Print(err)
	}
}

func cut(s string, fields []int) []string {
	all := strings.Fields(s)
	out := make([]string, len(fields))
	for i, n := range fields {
		if len(all) > n-1 {
			out[i] = all[n-1]
		}
	}
	return out
}
