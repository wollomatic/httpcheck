package main

import (
	"fmt"
	"io"
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
	for i, s := range sc.Service {
		if s.Name == "" {
			exitWithError(fmt.Errorf("Service #%v: name may not be empty", i+1))
		}
		if s.Url == "" {
			exitWithError(fmt.Errorf("URL of Service %s may not be empty", s.Name))
		}
		if s.Method == "" {
			sc.Service[i].Method = serviceDefaults.Method
		}
		sc.Service[i].Method = strings.ToUpper(sc.Service[i].Method)
		// check if String is in allowedMethods
		if !strings.Contains(allowedMethods, s.Method) {
			exitWithError(fmt.Errorf("method \"%s\" of Service %s is not allowed. Allowed methods are: %v", s.Method, s.Name, allowedMethods))
		}
		if s.Status == 0 {
			sc.Service[i].Status = serviceDefaults.Status
		}
		if s.Timeout == 0 {
			sc.Service[i].Timeout = serviceDefaults.Timeout
		}
		if s.Retries == 0 {
			sc.Service[i].Retries = serviceDefaults.Retries
		}
		if s.ErrDelay == 0 {
			sc.Service[i].ErrDelay = serviceDefaults.ErrDelay
		}
	}
	return sc
}
