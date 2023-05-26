package main

import (
	"net/http"
	"time"

	"github.com/wollomatic/httpcheck/pkg/httpcheck"
)

// checkService checks if a http/https service is available
// with ch as channel for service responses
// s is the service to check
// d is the delay to wait before checking
func checkService(ch chan<- httpcheck.Response, s httpcheck.Definition) {
	var o httpcheck.Response
	o.Service = s

	for o.Retries = 0; o.Retries <= s.Retries; o.Retries++ {

		if o.Retries != 0 {
			time.Sleep(time.Duration(s.ErrDelay) * time.Millisecond)
		}
		var response *http.Response
		var err error
		start := time.Now()
		response, err = httpcheck.Do(s.Url, s.Method, s.RequestContentType, s.RequestBody, s.Status, s.SearchText, time.Duration(s.Timeout)*time.Millisecond)
		o.RequestDuration = time.Since(start)
		o.Response = *response
		o.Err = err
		// if no error return immediately
		if o.Err == nil {
			ch <- o
			return
		}
	}
	// if we get here, we have reached the maximum number of retries
	ch <- o
}
