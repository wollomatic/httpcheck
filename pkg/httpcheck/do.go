package httpcheck

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Do checks a service and returns a http.Response object
func Do(url string, method string, status int, text string, timeout time.Duration) (*http.Response, error) {
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
