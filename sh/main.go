//go:generate go tool yacc -p "sh" sh.y

package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.SetPrefix("sh: ")
	log.SetFlags(0)
	in := bufio.NewReader(os.Stdin)
	for {
		if _, err := os.Stdout.WriteString(os.Getenv("PS1")); err != nil {
			log.Fatalf("WriteString: %s", err)
		}
		line, err := in.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}
		shParse(&shLex{line: expand(line)})
	}
}

func expand(cmd string) string {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return "\n"
	}
	var params string
	for _, arg := range args[1:] {
		deglobbed, err := filepath.Glob(arg)
		if err != nil || len(deglobbed) == 0 {
			params += " " + arg
		} else {
			params += " " + strings.Join(deglobbed, " ")
		}
	}
	cmd = args[0] + os.ExpandEnv(params)
	return cmd + "\n"
}
