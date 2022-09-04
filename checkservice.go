package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// checkService checks if a http/https service is available
// with ch as channel for service responses
// s is the service to check
// d is the delay to wait before checking
func checkService(ch chan serviceResponse, s Service, delay time.Duration) {
	var o serviceResponse
	o.service = s

	time.Sleep(delay)

	for o.retries = 0; o.retries <= s.Retries; o.retries++ {

		if o.retries != 0 {
			time.Sleep(time.Duration(s.ErrDelay) * time.Millisecond)
		}
		start := time.Now()
		response, err := httpcheck(s.Url, s.Status, s.Text, time.Duration(s.Timeout)*time.Millisecond)
		o.requestDuration = time.Since(start)
		o.response = *response
		o.err = err
		if o.err == nil {
			ch <- o
			return
		}
	}
	ch <- o
}

func httpcheck(url string, stat int, str string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return &http.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != stat {
		return resp, fmt.Errorf("status code does not match: got %d, want %d", resp.StatusCode, stat)
	}

	if str != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		if !strings.Contains(string(body), str) {
			return resp, fmt.Errorf("search string \"%s\" not found", str)
		}
	}

	return resp, nil
}
