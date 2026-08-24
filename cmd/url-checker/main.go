package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/musanii/go-url-checker/internal/checker"
)

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("at least one URL is required")
	}
	results := checker.CheckMultipleURLs(args)

	for _, result := range results {
		if result.Err != nil {
			fmt.Printf(
				"%s error: %v\n",
				result.URL,
				result.StatusCode,
			)
			continue

		}
		fmt.Printf(
			"%s  %d %s\n",
			result.URL,
			result.StatusCode,
			result.Duration.Round(time.Millisecond),
		)

	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		panic(err)
	}
}
