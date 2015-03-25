package main

import (
	"fmt"
	"os"
)

func main() {
	for _, e := range os.Environ() {
		fmt.Println(e)
	}
}
