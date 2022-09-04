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
		
		switch s.Test {
		case "GET":
			response, err = httpGETcheck(s.Url, s.Status, s.Text, time.Duration(s.Timeout)*time.Millisecond)
		case "HEAD":
			response, err = httpHEADcheck(s.Url, s.Status, time.Duration(s.Timeout)*time.Millisecond)
		default:
			response=&http.Response{}
			err = fmt.Errorf("unknown test type: %s", s.Test)
		}
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

func httpGETcheck(url string, stat int, str string, timeout time.Duration) (*http.Response, error) {
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

func httpHEADcheck(url string, stat int, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Head(url)
	if err != nil {
		return &http.Response{}, err
	}

	if resp.StatusCode != stat {
		return resp, fmt.Errorf("status code does not match: got %d, want %d", resp.StatusCode, stat)
	}

	return resp, nil
}
