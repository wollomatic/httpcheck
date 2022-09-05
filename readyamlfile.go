package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func readYamlFile(filename string) ServiceCatalog {
	// try to open and read given file
	file, err := os.Open(filename)
	if err != nil {
		exitWithError(err)
	}
	bs, err := io.ReadAll(file)
	if err != nil {
		exitWithError(err)
	}
	err = file.Close()
	if err != nil {
		exitWithError(err)
	}

	// unmarshal yaml file
	sc := ServiceCatalog{}
	err = yaml.Unmarshal(bs, &sc)
	if err != nil {
		exitWithError(err)
	}

	// check input data for needed values and fill missing optional values with defaults
	for i, s := range sc.Services {
		if s.Name == "" {
			exitWithError(fmt.Errorf("Service #%v: name may not be empty", i+1))
		}
		if s.Url == "" {
			exitWithError(fmt.Errorf("URL of Service %s may not be empty", s.Name))
		}
		if _, err := url.ParseRequestURI(s.Url); err != nil {
			exitWithError(fmt.Errorf("URL of Service %s is invalid: %v", s.Name, err))
		}
		if s.Test == "" {
			sc.Services[i].Test = serviceDefaults.Test
		}
		sc.Services[i].Test = strings.ToUpper(sc.Services[i].Test)
		// check if String is in allowedMethods
		if !strings.Contains(allowedTests, s.Test) {
			exitWithError(fmt.Errorf("test \"%s\" of Service %s is not allowed. Allowed tests are: %v", s.Test, s.Name, allowedTests))
		}
		if s.Test == "HEAD" && s.Text != "" {
			exitWithError(fmt.Errorf("text \"%s\" of Service %s (HEAD) is not allowed. Text is only allowed for GET test", s.Text, s.Name))
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
