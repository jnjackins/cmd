// awk built-in statements and functions

package main

import (
	"fmt"
	"strings"
)

func nop(record) {}

func printFn(args []*symbol) func(record) {
	dprintf("generating print function: args=%v", args)
	return func(rec record) {
		if len(args) == 0 {
			fmt.Println(rec)
		} else {
			printSyms(args)
		}
	}
}

func printSyms(args []*symbol) {
	ss := make([]string, len(args))
	for i := range args {
		ss[i] = args[i].getString()
	}
	fmt.Println(strings.Join(ss, ""))
}
