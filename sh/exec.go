package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"
	"unicode/utf8"
)

func execute(t *treeNode) int {
	switch t.typ {
	// foo; bar
	case ';':
		execute(t.children[0])
		return execute(t.children[1])

	// foo &
	case '&':
		cmd, err := t.mkFork()
		if err != nil {
			log.Print(err)
			return -1
		}
		if err := cmd.Start(); err != nil {
			log.Print(err)
			return -1
		}
		fmt.Printf("%d\n", cmd.Process.Pid)
		return 0

	// foo && bar
	case AND:
		if status := execute(t.children[0]); status == 0 {
			return execute(t.children[1])
		} else {
			return status
		}

	// foo || bar
	case OR:
		if status := execute(t.children[0]); status != 0 {
			return execute(t.children[1])
		} else {
			return status
		}

	// foo | bar
	case '|':
		r, w, err := os.Pipe()
		if err != nil {
			log.Print(err)
			return -1
		}
		t.children[0].io.pipeOut = w
		t.children[1].io.pipeIn = r
		go func() {
			execute(t.children[0])
			w.Close()
		}()
		defer r.Close()
		return execute(t.children[1])

	// foo bar
	case SIMPLE:
		args, vars := expandArgs(t.children[0], true)

		if len(args) == 0 {
			// only variable assignments
			for _, v := range vars {
				parts := strings.SplitN(v, "=", 2)
				setEnv(parts[0], parts[1])
			}
			return 0
		}

		if fn, ok := builtins[args[0]]; ok {
			exit, err := fn(args[1:])
			if err != nil {
				log.Printf("%s: %v", args[0], err)
			}
			return exit
		}

		cmd, err := t.mkCmd(args, vars)
		if err != nil {
			log.Print(err)
			return -1
		}
		//defer closeFds(cmd)

		dprintf("running simple command: %#v", cmd.Args)
		if err := cmd.Run(); err != nil {
			dprintf("run returned error: %v", err)
			return 1
		}
		return exitStatus(cmd.ProcessState)

	case IF:
		if execute(t.children[0]) == 0 {
			return execute(t.children[1])
		}
		return 0

	case FOR:
		assign := t.children[0].string
		in, _ := expandArgs(t.children[1], false)
		for _, s := range in {
			setEnv(assign, s)
			execute(t.children[2])
		}
		return 0

	default:
		log.Printf("not implemented: %#v", t)
		return -1
	}
}

func exitStatus(state *os.ProcessState) int {
	return state.Sys().(syscall.WaitStatus).ExitStatus()
}

func expandArgs(t *treeNode, doAssignments bool) (args, vars []string) {
	prologue := true
	for _, n := range t.children {
		s := n.string

		// don't expand quoted text
		if n.typ == QUOTE {
			args = append(args, s)
			continue
		}

		// variable assignments
		if prologue {
			i := strings.Index(s, "=")
			if i < 0 {
				prologue = false
			} else if readVarName(s[:i]) != "" {
				vars = append(vars, s)
				continue
			} else {
				prologue = false
			}
		}

		// expand variables
		if i := strings.Index(s, "$"); i >= 0 {
			name := readVarName(s[i+len("$"):])
			s = s[:i] + getEnv(name) + s[i+len("$")+len(name):]
		}

		// expand globs
		if strings.ContainsAny(s, "[?*") {
			matches, err := filepath.Glob(s)
			if err == nil {
				args = append(args, matches...)
				continue
			}
		}
		args = append(args, s)
	}
	return args, vars
}

func readVarName(s string) string {
	if len(s) == 0 {
		return ""
	}
	c, _ := utf8.DecodeRuneInString(s)
	if unicode.IsNumber(c) {
		return ""
	}
	var name string
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_' {
			name += string(c)
			continue
		}
		return name
	}
	return name
}

func (t *treeNode) mkCmd(args, vars []string) (*exec.Cmd, error) {
	if t.typ != SIMPLE {
		panic("mkCmd: bad node type")
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		return nil, fmt.Errorf("%s: command not found", args[0])
	}
	cmd := exec.Cmd{
		Path: path,
		Args: args,
	}
	if len(vars) > 0 {
		cmd.Env = append(os.Environ(), vars...)
	}

	if t.io != nil {
		cmd.Stdin = t.io.pipeIn
		cmd.Stdout = t.io.pipeOut
		for fd, redir := range t.io.redirs {
			dprintf("redirecting fd=%d to %s (append=%v)", fd, redir.path, redir.append)
			switch fd {
			case 0:
				f, err := os.Open(redir.path)
				if err != nil {
					return nil, err
				}
				cmd.Stdin = f
			case 1, 2:
				flag := os.O_RDWR
				if redir.append {
					flag |= os.O_APPEND
				} else {
					flag |= os.O_CREATE
				}
				f, err := os.OpenFile(redir.path, flag, 0666)
				if err != nil {
					return nil, err
				}
				if fd == 1 {
					cmd.Stdout = f
				} else if fd == 2 {
					cmd.Stderr = f
				}
			default:
				panic("TODO")
			}
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

func (t *treeNode) mkFork() (*exec.Cmd, error) {
	path, err := os.Executable()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(path, "-c", t.children[0].String())
	cmd.Args[0] = filepath.Base(cmd.Args[0] + "(fork)")
	dprintf("running forked command: %#v", cmd.Args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}
