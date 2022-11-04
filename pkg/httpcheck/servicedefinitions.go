package httpcheck

import (
	"net/http"
	"time"
)

type Definition struct {
	Name       string `yaml:"name"`
	Method     string `yaml:"method"`
	Url        string `yaml:"url"`
	Status     int    `yaml:"status"`
	SearchText string `yaml:"searchtext"`
	Timeout    int    `yaml:"timeout"`
	Retries    int    `yaml:"retries"`
	ErrDelay   int    `yaml:"err_delay"`
}

type Catalog struct {
	Delay    int          `yaml:"delay"`
	Services []Definition `yaml:"services"`
}

type Response struct {
	Service         Definition
	RequestDuration time.Duration
	Response        http.Response
	Retries         int
	Err             error
}

var serviceDefaults = Definition{
	Method:   "GET",
	Status:   200,
	Timeout:  1000,
	Retries:  0,
	ErrDelay: 100,
}

const allowedTests = "GET, HEAD"
