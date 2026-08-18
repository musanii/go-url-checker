package checker

import (
	"net/http"
	"sync"
	"time"
)

type CheckResult struct {
	URL        string
	StatusCode int
	Err        error
	Duration   time.Duration
}

func CheckURL(url string) CheckResult {
	return CheckURLWithTimeout(url, 7*time.Second)
}

func CheckURLWithTimeout(url string, timeout time.Duration) CheckResult {

	start := time.Now()

	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Get(url)

	duration := time.Since(start)

	if err != nil {
		return CheckResult{
			URL:      url,
			Err:      err,
			Duration: duration,
		}
	}

	defer response.Body.Close()

	return CheckResult{
		URL:        url,
		StatusCode: response.StatusCode,
		Duration:   duration,
	}

}

func CheckMultipleURLs(urls []string) []CheckResult {
	results := make([]CheckResult, len(urls))

	var wg sync.WaitGroup

	wg.Add(len(urls))

	for i, url := range urls {
		go func(i int, url string) {
			defer wg.Done()
			results[i] = CheckURL(url)

		}(i, url)

	}
	wg.Wait()

	return results

}
