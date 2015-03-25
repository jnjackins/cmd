package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"sigint.ca/die"
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
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: sed <command>")
		os.Exit(1)
	}
	cmd := flag.Arg(0)
	switch cmd[0] {
	case 's':
		sep := cmd[1]
		args := strings.Split(cmd, string(sep))
		if len(args) != 4 {
			fmt.Fprintln(os.Stderr, "sed: invalid substitution command")
			os.Exit(1)
		}
		args = args[1:]
		var global bool
		switch args[2] {
		case "g":
			global = true
		case "":
		default:
			fmt.Fprintln(os.Stderr, "sed: invalid substitution flag")
		}
		from, to := args[0], args[1]
		for {
			line, err := in.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					die.On(err, "sed: error reading from stdio")
				}
				break
			}
			sub(line, from, to, global)
		}
		out.Flush()

	default:
		fmt.Fprintln(os.Stderr, "sed: unimplemented command: "+string(cmd[0]))
		os.Exit(1)
	}
}

func sub(line, from, to string, global bool) {
	re, err := regexp.Compile(from)
	die.On(err, "sed: error parsing pattern \""+from+"\"")
	if global {
		fmt.Fprint(out, re.ReplaceAllString(line, to))
	} else {
		indices := re.FindStringIndex(line)
		if indices == nil {
			fmt.Fprint(out, line)
			return
		}
		i, j := indices[0], indices[1]
		fmt.Fprint(out, line[:i] + to + line[j:])
	}
}
