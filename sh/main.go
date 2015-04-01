// TODO: handle sigint correctly

//go:generate go tool yacc -p "sh" parse.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var lflag = flag.String("l", "", "read commands from `file` before reading normal input")

func main() {
	log.SetPrefix("sh: ")
	log.SetFlags(0)
	setupEnv()
	flag.Parse()
	if *lflag != "" {
		f, err := os.Open(*lflag)
		if err != nil {
			log.Print(err)
		} else {
			parse(f)
		}
	}
	parse(os.Stdin)
}

func parse(f *os.File) {
	in := bufio.NewReader(f)
	tty := isTTY(f)
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
