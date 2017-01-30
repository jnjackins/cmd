package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"sigint.ca/user"
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
	fmt.Fprintln(w, "PID\tUID\tTIME\tRSS\tSTATE\tCMD")
	defer w.Flush()
	for _, pid := range pids {
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}
		dir := "/proc/" + pid
		dirstat, err := os.Stat(dir)
		if err != nil {
			elog.Println(err)
			continue
		}
		uid := strconv.Itoa(int(dirstat.Sys().(*syscall.Stat_t).Uid))
		u, err := user.LookupId(uid)
		if err == nil {
			uid = u.Username
		}
		stat, err := ioutil.ReadFile(dir + "/stat")
		if err != nil {
			elog.Println(err)
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
			// could also check for empty /proc/<pid>/cmdline
			if pid == "2" || ppid == "2" {
				continue
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%v\t%s\t%s\t%s\n", pid, uid, time, rss, state, name)
	}
}
