// TODO: handle sigint correctly

//go:generate go tool yacc -p "sh" sh.y

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
		shParse(&shLex{line: os.ExpandEnv(line)})
	}
}
