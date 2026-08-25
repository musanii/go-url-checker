package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/musanii/go-url-checker/internal/checker"
)

type URLChecker interface {
	CheckMultipleURLs([]string) []checker.CheckResult
}

type checkerService struct{}

func (checkerService) CheckMultipleURLs(urls []string) []checker.CheckResult {
	return checker.CheckMultipleURLs(urls)
}

func run(args []string, urlChecker URLChecker) error {
	if len(args) == 0 {
		return errors.New("at least one URL is required")
	}
	results := urlChecker.CheckMultipleURLs(args)

	for _, result := range results {

		if result.Err != nil {
			fmt.Printf(
				"%s ERROR: %s\n",
				result.URL,
				result.Err,
			)
			return result.Err
		}
		fmt.Printf(
			"%s  %d %s\n",
			result.URL,
			result.StatusCode,
			result.Duration.Round(time.Millisecond),
		)

		if result.StatusCode < 200 || result.StatusCode >= 300 {
			return fmt.Errorf(
				"URL %s return HTTP status %d",
				result.URL,
				result.StatusCode,
			)
		}

	}
	return nil
}

func main() {

	urlChecker := checkerService{}
	if err := run(os.Args[1:], urlChecker); err != nil {
		panic(err)
	}
}
