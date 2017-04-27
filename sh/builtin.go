package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

var builtins = map[string]func([]string) (int, error){
	"cd":    cdFn,
	"exec":  execFn,
	"exit":  exitFn,
	"times": timesFn,
	"wait":  waitFn,
}

func cdFn(args []string) (int, error) {
	if len(args) == 0 {
		if err := os.Chdir(env["HOME"]); err != nil {
			return 1, err
		}
		return 0, nil
	}
	if err := os.Chdir(args[0]); err != nil {
		return 1, err
	}
	return 0, nil
}

func execFn(args []string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return 0, err
	}
	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		return 1, err
	}

	// not reached
	return 0, nil
}

func exitFn(args []string) (int, error) {
	if len(args) == 0 {
		exit(0)
	}
	i, err := strconv.Atoi(args[0])
	if err == nil {
		exit(i)
	}
	log.Printf("%s: bad number", args[0])
	exit(1)

	// not reached
	return 0, nil
}

func timesFn(args []string) (int, error) {
	var dst syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &dst); err != nil {
		return 1, err
	}
	utime := time.Duration(dst.Utime.Nano())
	stime := time.Duration(dst.Stime.Nano())
	fmt.Printf("%.0fm%.3fs %.0fm%.3fs\n",
		utime.Minutes(), utime.Seconds(),
		stime.Minutes(), stime.Seconds())

	return 0, nil
}

func waitFn(args []string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return 1, err
	}
	proc, err := os.FindProcess(i)
	if err != nil {
		return 1, err
	}
	state, err := proc.Wait()
	if err != nil {
		return 1, err
	}
	return exitStatus(state), nil
}
