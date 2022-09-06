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
		var response *http.Response
		var err error
		start := time.Now()
		response, err = httpcheck(s.Url, s.Method, s.Status, s.SearchText, time.Duration(s.Timeout)*time.Millisecond)
		o.requestDuration = time.Since(start)
		o.response = *response
		o.err = err
		// if no error return immediately
		if o.err == nil {
			ch <- o
			return
		}
	}
	// if we get here, we have reached the maximum number of retries
	ch <- o
}

func httpcheck(url string, method string, status int, text string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var resp *http.Response
	var err error
	switch method {
	case "GET":
		resp, err = client.Get(url)
	case "HEAD":
		resp, err = client.Head(url)
	default:
		resp = &http.Response{}
		err = fmt.Errorf("invalid method: %s", method)
	}
	if err != nil {
		return &http.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != status {
		return resp, fmt.Errorf("status code does not match: got %d, want %d", resp.StatusCode, status)
	}

	if text != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		if !strings.Contains(string(body), text) {
			return resp, fmt.Errorf("search string \"%s\" not found", text)
		}
	}

	return resp, nil
}
