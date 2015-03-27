// TODO: handle sigint correctly

//go:generate go tool yacc -p "sh" sh.y

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.SetPrefix("sh: ")
	log.SetFlags(0)
	setupEnv()
	in := bufio.NewReader(os.Stdin)
	tty := isTTY(os.Stdin)
	for {
		if tty {
			fmt.Print(os.Getenv("prompt"))
		}
		line, err := in.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
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
