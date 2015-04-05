package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	elog := log.New(os.Stderr, "ps: ", 0)
	f, err := os.Open("/proc")
	if err != nil {
		elog.Fatal(err)
	}
	pids, err := f.Readdirnames(0)
	if err != nil {
		elog.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "%4s %6s %s %s\n", "PID", "RSS", "STATE", "CMD")
	for _, pid := range pids {
		stat, err := ioutil.ReadFile("/proc/" + pid + "stat")
		if err != nil {
			log.Print(err)
			continue
		}
		fields := strings.Fields(string(stat))
		name := strings.Trim(fields[1], "()")
		state := fields[4]
		rss := fields[25]
		fmt.Printf("%4d %6d %s %s\n", pid, rss, state, name)
	}
}
