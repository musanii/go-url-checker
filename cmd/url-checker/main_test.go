package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRunRequiresURLs(t *testing.T) {
	err := run([]string{})

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

	output := captureOutput(func() {
		err := run([]string{server.URL})

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

	err := run([]string{server.URL})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunReportsURLCheckError(t *testing.T){
	output := captureOutput(func ()  {
		err := run([]string{"http://127.0.0.1:59999"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
	})

	if !strings.Contains(output, "error"){
		t.Fatalf(
			"expected output to contain %q, got %q",
			"error",
			output,
		)
	}
}
