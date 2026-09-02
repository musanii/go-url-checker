package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/musanii/go-url-checker/internal/checker"
)

type fakeURLChecker struct {
	results  []checker.CheckResult
	sequence []checker.CheckResult
	calls    int
}

func (f *fakeURLChecker) CheckMultipleURLs(urls []string) []checker.CheckResult {
	f.calls++
	if len(f.sequence) > 0 {
		result := f.sequence[f.calls-1]
		return []checker.CheckResult{result}

	}
	return f.results
}

func TestRunRequiresURLs(t *testing.T) {

	urlChecker := checkerService{}
	err := run([]string{}, urlChecker)

	if err == nil {
		t.Fatal("expected URL argument error, got nil")
	}

	expected := "at least one URL is required"

	if err.Error() != expected {
		t.Fatalf(
			"expected error %q, got %q",
			expected,
			err.Error(),
		)

	}
}

func captureOutput(fn func()) string {
	old := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestRunChecksURLs(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()
	urlChecker := checkerService{}

	output := captureOutput(func() {
		err := run([]string{server.URL}, urlChecker)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	if !strings.Contains(output, server.URL) {
		t.Fatalf(
			"expected output to contain %q, got %q",
			server.URL,
			output,
		)
	}

	if !strings.Contains(output, "200") {
		t.Fatalf(
			"expected output to contain status 200, got %q",
			output,
		)
	}

	err := run([]string{server.URL}, urlChecker)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunReportsURLCheckError(t *testing.T) {
	expectedErr := errors.New("connection refused")
	fake := fakeURLChecker{
		results: []checker.CheckResult{
			{
				URL: "http://example.com",
				Err: expectedErr,
			},
		},
	}
	output := captureOutput(func() {
		run([]string{"http://example.com"}, &fake)

	})

	if !strings.Contains(output, "ERROR") {
		t.Fatalf(
			"expected output to contain %q, got %q",
			"ERROR",
			output,
		)
	}
}

func TestRunReportsDuration(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()
	urlChecker := checkerService{}
	output := captureOutput(func() {
		err := run([]string{server.URL}, urlChecker)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	if !strings.Contains(output, "ms") {
		t.Fatalf(
			"expected output to contain duration in milliseconds, got %q",
			output,
		)
	}
}

func TestRunReturnsErrorWhenURLCheckFails(t *testing.T) {
	expectedErr := errors.New("connection refused")
	fake := fakeURLChecker{
		results: []checker.CheckResult{
			{
				URL: "http://example.com",
				Err: expectedErr,
			},
		},
	}
	err := run([]string{
		"http://example.com",
	}, &fake)

	if err != expectedErr {
		t.Fatalf("expected error %q, got %q",
			expectedErr,
			err,
		)
	}
}

func TestRunSucceedsWhenURLReturnsSuccessStatus(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()
	urlChecker := checkerService{}

	err := run([]string{server.URL}, urlChecker)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

}

func TestRunReturnsErrorWhenURLReturnsServerError(t *testing.T) {

	fake := fakeURLChecker{
		results: []checker.CheckResult{
			{
				URL:        "http://example.com",
				StatusCode: 500,
			},
		},
	}
	err := run([]string{"http://example.com"}, &fake)

	if err == nil {
		t.Fatalf("expected  error for HTTP 500, got nil")
	}

}

func TestRunChecksAllURLsBeforeReturningError(t *testing.T) {
	successURL := "http://success.example.com"
	failureURL := "http://failure.example.com"

	fake := fakeURLChecker{
		results: []checker.CheckResult{
			{
				URL:        successURL,
				StatusCode: 200,
			},
			{
				URL:        failureURL,
				StatusCode: 500,
			},
		},
	}

	var err error
	output := captureOutput(func() {
		err = run([]string{
			successURL,
			failureURL,
		}, &fake)
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(output, successURL) {
		t.Fatalf(
			"expected output to contain %q, got %q",
			successURL,
			output,
		)
	}

	if !strings.Contains(output, failureURL) {
		t.Fatalf(
			"expected output to contain %q, got %q",
			failureURL,
			output,
		)
	}
}

func TestRunReportsURLCheckErrorDetail(t *testing.T) {
	expectedErr := errors.New("connection refused")

	fake := fakeURLChecker{
		results: []checker.CheckResult{
			{
				URL: "http://example.com",
				Err: expectedErr,
			},
		},
	}

	output := captureOutput(func() {
		run([]string{"http://example.com"}, &fake)
	})

	if !strings.Contains(output, "ERROR") {
		t.Fatalf(
			"expected output to contain status ERROR,got %q",
			output,
		)
	}

	if !strings.Contains(output, expectedErr.Error()) {
		t.Fatalf(
			"expected output to contain %q, got %q",
			expectedErr.Error(),
			output,
		)
	}
}

func TestMonitorChecksURLsRepeatedly(t *testing.T) {
	fake := &fakeURLChecker{}
	urlMonitor := NewURLMonitor()

	monitor(
		[]string{"http://example.com"},
		10*time.Millisecond,
		3,
		fake,
		urlMonitor,
	)

	if fake.calls != 3 {
		t.Fatalf(
			"expected 3 checks, got %d",
			fake.calls,
		)
	}
}

func TestURLMonitorDetectsStateChange(t *testing.T) {
	monitor := NewURLMonitor()

	monitor.recordState("http://example.com", URLUp)

	changed := monitor.recordState("http://example.com", URLDown)

	if !changed {
		t.Fatal("expected state change, got none")
	}

}

func TestURLMonitorDoesNotDetectChangeWhenStateIsSame(t *testing.T) {
	monitor := NewURLMonitor()
	monitor.recordState("http://example.com", URLUp)

	changed := monitor.recordState("http://example.com", URLUp)

	if changed {
		t.Fatalf("expected no change, got %v", changed)

	}

}

func TestURLMonitorDoesNotDetectChangeForFirstObservation(t *testing.T) {
	monitor := NewURLMonitor()
	changed := monitor.recordState("http://example.com", URLUp)

	if changed {
		t.Fatal("expected no state change, got a change")
	}

}

func TestURLMonitorRecordsURLAsUp(t *testing.T) {
	monitor := NewURLMonitor()
	monitor.recordState("http://example.com", URLUp)
	state := monitor.states["http://example.com"]
	if state != URLUp {
		t.Fatalf("expected the state to be UP, got %v", state)
	}
}

func TestDetermineStateReturnsUpWhenCheckSucceeds(t *testing.T) {
	result := checker.CheckResult{
		StatusCode: 200,
	}

	state := determineState(result)

	if state != URLUp {
		t.Fatalf("expected URL to be UP, got %v", state)
	}
}

func TestMonitorDetectsStateChange(t *testing.T) {
	fakeChecker := &fakeURLChecker{
		sequence: []checker.CheckResult{
			{
				URL:        "http://example.com",
				StatusCode: 200,
			},
			{
				URL:        "http://example.com",
				StatusCode: 500,
			},
		},
	}

	urlMonitor := NewURLMonitor()

	monitor(
		[]string{"http://example.com"},
		1*time.Millisecond,
		2,
		fakeChecker,
		urlMonitor,
	)

	state := urlMonitor.states["http://example.com"]

	if state != URLDown {
		t.Fatalf("expected final state to be DOWN, got %v", state)
	}
}
