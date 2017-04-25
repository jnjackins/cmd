package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
)

//go:generate goyacc -o y.out.go -p "sh" syntax.y
//go:generate rm y.output

var (
	cflag = flag.String("c", "", "Read comands from `string`.")
	debug = flag.Bool("d", false, "Enable debug mode.")
)

func main() {
	flag.Parse()

	log.SetPrefix("sh: ")
	log.SetFlags(0)

	setupEnv()

	var input io.Reader = os.Stdin
	prompt := true
	if *cflag != "" {
		input = bytes.NewBufferString(*cflag + "\n")
		prompt = false
	}
	in := bufio.NewReader(input)
	for {
		if prompt {
			if _, err := os.Stdout.WriteString(env["PS1"]); err != nil {
				log.Fatalf("WriteString: %s", err)
			}
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

func dprintf(format string, args ...interface{}) {
	if *debug {
		log.Printf("debug: "+format, args...)
	}
}
