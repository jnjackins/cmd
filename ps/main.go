package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	fmt.Fprintf(os.Stderr, "%5s %6s %s %s\n", "PID", "RSS", "STATE", "CMD")
	for _, pid := range pids {
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}
		stat, err := ioutil.ReadFile("/proc/" + pid + "/stat")
		if err != nil {
			elog.Print(err)
			continue
		}
		fields := strings.Fields(string(stat))
		name := strings.Trim(fields[1], "()")
		state := fields[2]
		rss := fields[23]
		fmt.Printf("%5s %6s %s %s\n", pid, rss, state, name)
	}
}
