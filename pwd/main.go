package main

import (
	"fmt"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
	}
	fmt.Println(wd)
}
