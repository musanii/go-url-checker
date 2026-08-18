package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/musanii/go-url-checker/internal/checker"
)

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("at least one URL is required")
	}
	results := checker.CheckMultipleURLs(args)

	for _,result := range results{
		fmt.Printf(
			"%s %d\n",
			result.URL,
			result.StatusCode,
		)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		panic(err)
	}
}
