package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	elog := log.New(os.Stderr, "time: ", 0)
	args := os.Args[1:]
	var cmd *exec.Cmd
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: time command")
		os.Exit(1)
	}
	cmd = exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	start := time.Now()
	err := cmd.Run()
	if err != nil {
		elog.Fatal(err)
	}
	r := time.Since(start)
	u := cmd.ProcessState.UserTime()
	s := cmd.ProcessState.SystemTime()
	fmt.Fprintf(os.Stderr, "%.2fu %.2fs %.2fr", u.Seconds(), s.Seconds(), r.Seconds())
	fmt.Fprintf(os.Stderr, "\t%s\n", strings.Join(args, " "))
}
