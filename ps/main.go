package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

var kflag = flag.Bool("k", false, "Show kernel threads (processes with a parent PID of 0).")

func main() {
	elog := log.New(os.Stderr, "ps: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ps [options]")
		flag.PrintDefaults()
	}
	flag.Parse()
	tps := ticksPerSecond()
	f, err := os.Open("/proc")
	if err != nil {
		elog.Fatal(err)
	}
	pids, err := f.Readdirnames(0)
	if err != nil {
		elog.Fatal(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "PID\tPPID\tTIME\tRSS\tSTATE\tCMD")
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
		ppid := fields[3]
		utime, _ := strconv.Atoi(fields[13])
		stime, _ := strconv.Atoi(fields[14])
		time := time.Second * time.Duration((utime+stime)/tps)
		rss := fields[23]
		if !*kflag {
			// TODO: find a better way to identify kernel threads
			if pid == "2" || ppid == "2" {
				continue
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%v\t%s\t%s\t%s\n", pid, ppid, time, rss, state, name)
	}
}
