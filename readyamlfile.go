package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func readYamlFile(filename string) ServiceCatalog {

	bs, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	// unmarshal yaml file
	sc := ServiceCatalog{}
	err = yaml.Unmarshal(bs, &sc)
	if err != nil {
		log.Fatalln(err)
	}

	// check input data for needed values and fill missing optional values with defaults
	for i, s := range sc.Services {
		if s.Name == "" {
			log.Fatalf("Service #%v: name may not be empty\n", i+1)
		}
		if s.Url == "" {
			log.Fatalf("URL of Service %s may not be empty\n", s.Name)
		}
		if _, err := url.ParseRequestURI(s.Url); err != nil {
			log.Fatalf("URL of Service %s is invalid: %v\n", s.Name, err)
		}
		if s.Method == "" {
			sc.Services[i].Method = serviceDefaults.Method
		}
		sc.Services[i].Method = strings.ToUpper(sc.Services[i].Method)
		// check if String is in allowedMethods
		if !strings.Contains(allowedTests, s.Method) {
			log.Fatalf("test \"%s\" of Service %s is not allowed. Allowed tests are: %v\n", s.Method, s.Name, allowedTests)
		}
		if s.Method == "HEAD" && s.SearchText != "" {
			log.Fatalf("text \"%s\" of Service %s (HEAD) is not allowed. Text is only allowed for GET test\n", s.SearchText, s.Name)
		}
		if s.Status == 0 {
			sc.Services[i].Status = serviceDefaults.Status
		}
		if s.Timeout == 0 {
			sc.Services[i].Timeout = serviceDefaults.Timeout
		}
		if s.Retries == 0 {
			sc.Services[i].Retries = serviceDefaults.Retries
		}
		if s.ErrDelay == 0 {
			sc.Services[i].ErrDelay = serviceDefaults.ErrDelay
		}
	}
	return sc
}
