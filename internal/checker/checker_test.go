package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckURL(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	result := CheckURL(server.URL)

	if result.URL != server.URL {
		t.Fatalf("expected URL %q, got %q", server.URL, result.URL)
	}

	if result.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status code %d, got %d",
			http.StatusOK,
			result.StatusCode,
		)
	}

	if result.Err != nil {
		t.Fatalf("expected no error, got %v", result.Err)
	}
}

func TestCheckURLMeasuresDuration(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	result := CheckURL(server.URL)

	if result.Duration <= 0 {
		t.Fatal("expected duration to be greater than zero")
	}
}

func TestCheckURLReturnsErrorForUnreachableURL(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	url := server.URL

	server.Close()

	result := CheckURL(url)

	if result.Err == nil {
		t.Fatal("expected an error, got nil")
	}
	if result.URL != url {
		t.Fatalf(
			"expected URL %q, got %q",
			url,
			result.URL,
		)
	}
}

func TestCheckURLReturnsHTTPErrorStatus(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)

	defer server.Close()

	result := CheckURL(server.URL)

	if result.Err != nil {
		t.Fatalf("expected no request error, got %v", result.Err)
	}

	if result.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"expected status code %d, got %d",
			http.StatusNotFound,
			result.StatusCode,
		)
	}

}

func TestCheckURLTimesOut(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	result := CheckURLWithTimeout(server.URL, 50*time.Millisecond)

	if result.Err == nil {
		t.Fatal("expected timeout error, got nil")
	}

}

func TestCheckURLUsesDefaultTimeout(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	result := CheckURL(server.URL)

	if result.Err != nil {
		t.Fatalf("expected no error, got %v", result.Err)
	}

	if result.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status code %d, got %d",
			http.StatusOK,
			result.StatusCode,
		)
	}
}

func TestCheckMultipleURLs(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	urls := []string{
		server.URL,
		server.URL,
	}

	results := CheckMultipleURLs(urls)

	if len(results) != len(urls) {

		t.Fatalf(
			"expected %d results, got %d",
			len(urls),
			len(results),
		)
	}

	for i, result := range results {
		if result.URL != urls[i] {

			t.Fatalf(
				"expected URL %q at index %d, got %q",
				urls[i],
				i,
				result.URL,
			)
		}

		if result.Err != nil {
			t.Fatalf("expected no error, got %v", result.Err)
		}

		if result.StatusCode != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				result.StatusCode,
			)
		}
	}

}
