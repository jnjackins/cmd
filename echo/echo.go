package main

import (
	"flag"
	"fmt"
	"strings"
)

var nflag = flag.Bool("n", false, "omit trailing newline")

func main() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), " "))
	if !*nflag {
		fmt.Println()
	}
}
