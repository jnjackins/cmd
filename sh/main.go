//go:generate go tool yacc -p "sh" sh.y

package main

import (
	"bufio"
	"io"
	"log"
	"os"
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
		shParse(&shLex{line: os.ExpandEnv(line)})
	}
}
