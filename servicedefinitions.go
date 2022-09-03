package main

import (
	"net/http"
	"time"
)

type Service struct {
	Name     string `yaml:"name"`
	Method   string `yaml:"method"`
	Url      string `yaml:"url"`
	Status   int    `yaml:"status"`
	Text     string `yaml:"text"`
	Timeout  int    `yaml:"timeout"`
	Retries  int    `yaml:"retries"`
	ErrDelay int    `yaml:"err_delay"`
}

type ServiceCatalog struct {
	Delay   int       `yaml:"delay"`
	Service []Service `yaml:"service"`
}

type serviceResponse struct {
	service         Service
	requestDuration time.Duration
	response        http.Response
	retries         int
	err             error
}

var serviceDefaults = Service{
	Method:   "GET",
	Status:   200,
	Timeout:  1000,
	Retries:  0,
	ErrDelay: 100,
}

const allowedMethods = "GET"
