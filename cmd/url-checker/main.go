package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/musanii/go-url-checker/internal/checker"
)

type URLChecker interface {
	CheckMultipleURLs([]string) []checker.CheckResult
}

type URLState string

const (
	URLUp   URLState = "UP"
	URLDown URLState = "DOWN"
)

type URLMonitor struct {
	states map[string]URLState
}

func NewURLMonitor() *URLMonitor {
	return &URLMonitor{
		states: make(map[string]URLState),
	}
}

func (m *URLMonitor) recordState(url string, state URLState) {
	m.states[url] = state
}

type checkerService struct{}

func (checkerService) CheckMultipleURLs(urls []string) []checker.CheckResult {
	return checker.CheckMultipleURLs(urls)
}

func monitor(ctx context.Context, urls []string, interval time.Duration, urlChecker URLChecker) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		urlChecker.CheckMultipleURLs(urls)

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
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
