package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func execute(n node) int {
	switch t := n.(type) {
	case *listNode:
		switch t.typ {
		case typeListSequence:
			execute(t.left)
			return execute(t.right)
		case typeListAnd:
			if status := execute(t.left); status == 0 {
				return execute(t.right)
			} else {
				return status
			}
		case typeListOr:
			if status := execute(t.left); status != 0 {
				return execute(t.right)
			} else {
				return status
			}
		default:
			panic("TODO")
		}

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

	case *simpleNode:
		cmd, err := t.mkCmd()
		if err != nil {
			log.Print(err)
			return -1
		}
		cmd.Run()
		return exitStatus(cmd)

	default:
		log.Printf("not handled: %T", t)
		printTree(n)
		return -1
	}
	panic("unreached")
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
