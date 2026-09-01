package main

import (
	"bytes"
	"context"
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
	results []checker.CheckResult
	calls   int
}

func (f *fakeURLChecker) CheckMultipleURLs(urls []string) []checker.CheckResult {
	f.calls++
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fake := &fakeURLChecker{}

	go monitor(
		ctx,
		[]string{"http://example.com"},
		10*time.Millisecond,
		fake,
	)

	time.Sleep(25 * time.Millisecond)
	cancel()

	if fake.calls < 2 {
		t.Fatalf(
			"expected 2 checks, got %d",
			fake.calls,
		)
	}
}

func TestMonitorStopsWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fake := &fakeURLChecker{}

	done := make(chan struct{})

	go func() {
		monitor(
			ctx,
			[]string{"http://example.com"},
			10*time.Millisecond,
			fake,
		)

		close(done)

	}()

	time.Sleep(25 * time.Millisecond)

	cancel()

	select {
	case <-done:
		//monitor stopped
	case <-time.After(100 * time.Millisecond):
		t.Fatal("monitor did not stop after context cancellation")
	}

}

func TestURLMonitorRecordsURLAsUp(t *testing.T) {
	monitor := NewURLMonitor()

	monitor.recordState("http://example.com", URLUp)
	if monitor.states["http://example.com"] != URLUp {
		t.Fatalf(
			"expected URL state to be %q, got %q",
			URLUp,
			monitor.states["http://example.com"],
		)
	}
}

func TestURLMonitorUpdatesURLState(t *testing.T) {
	monitor := NewURLMonitor()

	monitor.recordState("http://example.com", URLUp)
	monitor.recordState("http://example.com", URLDown)

	if monitor.states["http://example.com"] != URLDown {
		t.Fatalf(
			"expected URL state to be %q, got %q",
			URLDown,
			monitor.states["http://example.com"],
		)
	}

}

func TestURLMonitorDetectsStateChange(t *testing.T) {
	monitor := NewURLMonitor()

	monitor.recordState("http://example.com", URLUp)

	changed := monitor.recordState("http://example.com", URLDown)

	if !changed{
		t.Fatal("expected state change, got none")
	}
	
	
}
