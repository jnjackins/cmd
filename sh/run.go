package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func runLine(line [][]*exec.Cmd) {
	// During parsing, it is convenient to append to the list
	// of pipes from the last to the first, so iterate over them
	// in reverse to restore the correct order.
	for i := len(line) - 1; i >= 0; i-- {
		pipe := line[i]
		for j := 0; j < len(pipe); j++ {
			cmd := pipe[j]
			if fn, ok := builtins[cmd.Path]; ok {
				fn(cmd.Path, cmd.Args)
			} else {
				if err := start(cmd); err != nil {
					log.Print(err)
					break
				} else {
					defer wait(cmd)
				}
			}
		}
	}
}

// Start reports whether it started cmd.
func start(cmd *exec.Cmd) error {
	if filepath.Base(cmd.Path) == cmd.Path {
		if lp, err := exec.LookPath(cmd.Path); err != nil {
			return err
		} else {
			cmd.Path = lp
		}
	}
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	for key, val := range env {
		cmd.Env = append(cmd.Env, key+"="+val)
		delete(env, key)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func wait(cmd *exec.Cmd) {
	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}
	status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	if err := os.Setenv("status", strconv.Itoa(status)); err != nil {
		log.Print(err)
	}
}

func connect(cmd1, cmd2 *exec.Cmd) {
	stdout, err := cmd1.StdoutPipe()
	if err != nil {
		log.Print(err)
	}
	cmd2.Stdin = stdout
}

func open(path string, mode int) *os.File {
	switch mode {
	case 'r':
		mode = os.O_RDONLY
	case 'w':
		mode = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case 'a':
		mode = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		panic("open: invalid mode")
	}
	f, err := os.OpenFile(path, mode, 0666)
	if err != nil {
		log.Print(err)
	}
	return f
}

func close(closer interface{}) {
	switch c := closer.(type) {
	case io.Closer:
		err := c.Close()
		if err != nil {
			log.Print(err)
		}
	default:
		panic("sh: close: argument is not an io.Closer")
	}
}
