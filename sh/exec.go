package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func execute(n node) int {
	switch t := n.(type) {
	// foo bar &
	case *forkNode:
		path, err := os.Executable()
		if err != nil {
			log.Print(err)
			return -1
		}
		cmd := exec.Command(path, "-c", t.tree.String())
		cmd.Args[0] = filepath.Base(cmd.Args[0] + "(fork)")
		dprintf("forking: %#v", cmd.Args)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Print(err)
			return -1
		}
		fmt.Printf("%d\n", cmd.Process.Pid)
		return 0

	case *listNode:
		switch t.typ {
		// foo; bar
		case typeListSequence:
			execute(t.left)
			return execute(t.right)
		// foo && bar
		case typeListAnd:
			if status := execute(t.left); status == 0 {
				return execute(t.right)
			} else {
				return status
			}
		// foo || bar
		case typeListOr:
			if status := execute(t.left); status != 0 {
				return execute(t.right)
			} else {
				return status
			}

		default:
			panic("bad listNode")
		}

	// foo | bar
	case *pipeNode:
		r, w, err := os.Pipe()
		if err != nil {
			log.Print(err)
			return -1
		}
		t.left.setStdout(w)
		t.right.setStdin(r)
		go func() {
			execute(t.left)
			w.Close()
		}()
		defer r.Close()
		return execute(t.right)

	// foo bar
	case *simpleNode:
		cmd, err := t.mkCmd()
		if err != nil {
			log.Print(err)
			return -1
		}
		dprintf("running simple command: %#v", cmd.Args)
		cmd.Run()
		return exitStatus(cmd)

	default:
		log.Printf("not handled: %T", t)
		printTree(n)
		return -1
	}
}

func exitStatus(cmd *exec.Cmd) int {
	return cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
}

func (c *simpleNode) mkCmd() (*exec.Cmd, error) {
	var args []string
	for p := c.args; p != nil; p = p.next {
		args = append(args, p.val)
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return nil, fmt.Errorf("%s: command not found", args[0])
	}
	cmd := exec.Cmd{
		Path:   path,
		Args:   args,
		Stdin:  c.stdin,
		Stdout: c.stdout,
		Stderr: c.stderr,
	}

	for fd, path := range c.redirs {
		switch fd {
		case 0:
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			cmd.Stdin = f
		case 1:
			f, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			cmd.Stdout = f
		case 2:
			f, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			cmd.Stderr = f
		default:
			panic("TODO")
		}
	}
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	return &cmd, nil
}
