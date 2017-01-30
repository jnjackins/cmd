package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

var signals = map[string]syscall.Signal{
	"HUP":  unix.SIGHUP,
	"INT":  unix.SIGINT,
	"QUIT": unix.SIGQUIT,
	"ILL":  unix.SIGILL,
	"ABRT": unix.SIGABRT,
	"FPE":  unix.SIGFPE,
	"KILL": unix.SIGKILL,
	"SEGV": unix.SIGSEGV,
	"PIPE": unix.SIGPIPE,
	"ALRM": unix.SIGALRM,
	"TERM": unix.SIGTERM,
	"USR1": unix.SIGUSR1,
	"USR2": unix.SIGUSR2,
	"CHLD": unix.SIGCHLD,
	"CONT": unix.SIGCONT,
	"STOP": unix.SIGSTOP,
	"TSTP": unix.SIGTSTP,
	"TTIN": unix.SIGTTIN,
	"TTOU": unix.SIGTTOU,
}

var sflag = flag.String("s", "TERM", "Specify the `signal` to send.")

func main() {
	log.SetPrefix("kill: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: kill [-s signal] pid")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	pid, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal(err)
	}
	var sig syscall.Signal
	if n, err := strconv.Atoi(*sflag); err == nil {
		sig = syscall.Signal(n)
	} else {
		var ok bool
		if sig, ok = signals[*sflag]; !ok {
			log.Fatalf("invalid signal name: %s", *sflag)
		}
	}
	if err := proc.Signal(sig); err != nil {
		log.Fatal(err)
	}
}
