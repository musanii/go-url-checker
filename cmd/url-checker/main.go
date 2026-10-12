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

type ContextURLChecker interface {
	CheckMultipleURLsWithContext(context.Context, []string) []checker.CheckResult
}

type URLState string

const (
	URLUp   URLState = "UP"
	URLDown URLState = "DOWN"
)

type URLMonitor struct {
	states    map[string]URLState
	latencies map[string]time.Duration
	history   map[string][]CheckHistory
}

type CheckHistory struct {
	Timestamp time.Time
	State     URLState
	Latency   time.Duration
}

type Alert struct {
	URL       string
	Previous  URLState
	Current   URLState
	Timestamp time.Time
}

func NewURLMonitor() *URLMonitor {
	return &URLMonitor{
		states:    make(map[string]URLState),
		latencies: make(map[string]time.Duration),
		history:   make(map[string][]CheckHistory),
	}
}

func (m *URLMonitor) recordState(url string, state URLState) (URLState, bool) {
	previousState, ok := m.states[url]
	changed := ok && previousState != state
	m.states[url] = state

	return previousState, changed
}

type checkerService struct{}

func (checkerService) CheckMultipleURLs(urls []string) []checker.CheckResult {
	return checker.CheckMultipleURLs(urls)
}

func (checkerService) CheckMultipleURLsWithContext(ctx context.Context, urls []string) []checker.CheckResult {
	return checker.CheckMultipleURLsWithContext(ctx, urls)
}

func monitor(ctx context.Context, urls []string, interval time.Duration, checks int, urlChecker URLChecker, urlMonitor *URLMonitor, alert func(Alert)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < checks; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		contextChecker, ok := urlChecker.(ContextURLChecker)
		if !ok {
			return
		}
		results := contextChecker.CheckMultipleURLsWithContext(ctx, urls)

		for _, result := range results {
			state := determineState(result)
			previousState, changed := urlMonitor.recordState(result.URL, state)

			urlMonitor.latencies[result.URL] = result.Duration

			urlMonitor.history[result.URL] = append(
				urlMonitor.history[result.URL],
				CheckHistory{
					Timestamp: time.Now(),
					State:     state,
					Latency:   result.Duration,
				},
			)

			if changed {
				alert(Alert{
					URL:       result.URL,
					Previous:  previousState,
					Current:   state,
					Timestamp: time.Now(),
				})
			}

		}

		if i < checks-1 {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if ctx.Err() != nil {
					return
				}
			}
		}
	}
}

func determineState(result checker.CheckResult) URLState {
	if result.Err == nil && result.StatusCode < 400 {
		return URLUp
	}
	return URLDown
}

func isSlowResponse(duration time.Duration) bool {
	return duration >= 500*time.Millisecond
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
