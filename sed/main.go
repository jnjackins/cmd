package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	in  *bufio.Reader
	out *bufio.Writer
)

func init() {
	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
}

func main() {
	elog := log.New(os.Stderr, "sed: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: sed <command>")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	cmd := flag.Arg(0)
	switch cmd[0] {
	case 's':
		sep := cmd[1]
		args := strings.Split(cmd, string(sep))
		if len(args) != 4 {
			elog.Fatal("invalid substitution command")
		}
		args = args[1:]
		var global bool
		switch args[2] {
		case "g":
			global = true
		case "":
		default:
			elog.Fatalf("invalid substitution flag: %s", args[2])
		}
		from, to := args[0], args[1]
		for {
			line, err := in.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					elog.Fatal(err)
				}
				break
			}
			if err := sub(line, from, to, global); err != nil {
				elog.Fatal(err)
			}
		}
		out.Flush()

	default:
		elog.Fatalf("unimplemented command: %s", string(cmd[0]))
	}
}

func sub(line, from, to string, global bool) error {
	re, err := regexp.Compile(from)
	if err != nil {
		return err
	}
	if global {
		fmt.Fprint(out, re.ReplaceAllString(line, to))
	} else {
		indices := re.FindStringIndex(line)
		if indices == nil {
			fmt.Fprint(out, line)
			return nil
		}
		i, j := indices[0], indices[1]
		fmt.Fprint(out, line[:i]+to+line[j:])
	}
	return nil
}
