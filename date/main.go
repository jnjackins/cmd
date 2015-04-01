package main

import (
	"flag"
	"fmt"
	"time"
)

var nflag = flag.Bool("num", false, "Print the date and time in a numeric format")
var fflag = flag.String("f", "", "Provide a custom Go time `format` (ref: Jan 2 3:04:05PM 2006 -0700)")

func main() {
	flag.Parse()
	format := "Mon Jan 2 3:04:05 PM MST 2006"
	if *nflag {
		format = "2006-01-02 15:04:05 -07:00"
	}
	if *fflag != "" {
		format = *fflag
	}
	fmt.Println(time.Now().Format(format))
}
