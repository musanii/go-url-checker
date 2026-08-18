package checker

import (
	"net/http"
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

func CheckMultipleURLs(urls []string)[]CheckResult{
	results := make([]CheckResult, 0, len(urls))

	for _,url := range urls{
		results = append(results,CheckURL(url))
	}

	return results

}
