package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "PID\tRSS\tSTATE\tCMD")
	defer w.Flush()
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
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", pid, rss, state, name)
	}
}
