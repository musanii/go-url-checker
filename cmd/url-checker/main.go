package main

import (
	"errors"
	"os"
)

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("at least one URL is required")
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		panic(err)
	}
}
