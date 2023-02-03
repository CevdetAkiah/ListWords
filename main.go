package main

import (
	"os"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	reader := os.Stdin
	ui(done, reader)
}
