package main

import (
	"io"
	"log"
	"os"
)

func main() {
	log.SetPrefix("cat: ")
	log.SetFlags(0)
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "/dev/stdin")
	}
	var errc int
	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			log.Print(err)
			errc++
			continue
		}
		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			log.Print(err)
			errc++
		}
		f.Close()
	}
	os.Exit(errc)
}
