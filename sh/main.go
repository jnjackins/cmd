//go:generate go tool yacc -p "sh" sh.y

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

func main() {
	log.SetPrefix("sh: ")
	log.SetFlags(0)
	interrupt := make(chan os.Signal)
	go func() {
		for {
			_ = <-interrupt
			fmt.Println()
			fmt.Print(os.Getenv("PS1"))
		}
	}()
	signal.Notify(interrupt, os.Interrupt)
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(os.Getenv("PS1"))
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
