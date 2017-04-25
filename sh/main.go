package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

//go:generate goyacc -o y.out.go -p "sh" syntax.y
//go:generate rm y.output

func main() {
	log.SetPrefix("sh: ")
	log.SetFlags(0)

	shDebug = 1

	in := bufio.NewReader(os.Stdin)
	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			log.Fatalf("WriteString: %s", err)
		}
		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}

		shParse(&shLex{line: line})
	}
}
