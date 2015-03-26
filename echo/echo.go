package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var nflag bool
	args := os.Args
	if len(args) > 1 && args[1] == "-n" {
		nflag = true
		args = args[1:]
	}
	fmt.Print(strings.Join(args[1:], " "))
	if !nflag {
		fmt.Println()
	}
}
