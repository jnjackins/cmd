package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
)

var signals = map[string]syscall.Signal{
	"HUP":  syscall.SIGHUP,
	"INT":  syscall.SIGINT,
	"QUIT": syscall.SIGQUIT,
	"ILL":  syscall.SIGILL,
	"ABRT": syscall.SIGABRT,
	"FPE":  syscall.SIGFPE,
	"KILL": syscall.SIGKILL,
	"SEGV": syscall.SIGSEGV,
	"PIPE": syscall.SIGPIPE,
	"ALRM": syscall.SIGALRM,
	"TERM": syscall.SIGTERM,
	"USR1": syscall.SIGUSR1,
	"USR2": syscall.SIGUSR2,
	"CHLD": syscall.SIGCHLD,
	"CONT": syscall.SIGCONT,
	"STOP": syscall.SIGSTOP,
	"TSTP": syscall.SIGTSTP,
	"TTIN": syscall.SIGTTIN,
	"TTOU": syscall.SIGTTOU,
}

func main() {
	elog := log.New(os.Stderr, "kill: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: kill pid [signal]")
	}
	flag.Parse()
	if len(flag.Args()) == 0 || len(flag.Args()) > 2 {
		flag.Usage()
		os.Exit(1)
	}
	pid, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		elog.Fatal(err)
	}
	sig := syscall.SIGTERM
	if len(flag.Args()) == 2 {
		var ok bool
		if sig, ok = signals[flag.Arg(1)]; !ok {
			elog.Fatalf("invalid signal: %s", flag.Arg(1))
		}
	}
	elog.Printf("sending signal #%d to %d\n", sig, pid)
	if err := syscall.Kill(pid, sig); err != nil {
		elog.Fatal(err)
	}
}
