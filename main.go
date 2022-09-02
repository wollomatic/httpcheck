package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Status   int    `yaml:"status"`
	Text     string `yaml:"text"`
	Timeout  int    `yaml:"timeout"`
	Retries  int    `yaml:"retries"`
	ErrDelay int    `yaml:"err_delay"`
}

type ServiceCatalog struct {
	Delay          int       `yaml:"delay"`
	ServiceCatalog []Service `yaml:"service_catalog"`
}

type serviceResponse struct {
	service         Service
	requestDuration time.Duration
	gottenStatus    int
	retries         int
	err             error
}

func main() {

	// if no filename argument is given show help and exit
	if len(os.Args) != 2 {
		showHelp()
	}

	// try to open and read given file
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	bs, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()

	// unmarshal yaml file
	sc := ServiceCatalog{}
	err = yaml.Unmarshal(bs, &sc)
	if err != nil {
		log.Fatal(err)
	}

	// create channel for service responses
	ch := make(chan serviceResponse)

	// start concurrent service checks
	delay := 0
	for _, s := range sc.ServiceCatalog {
		go checkService(ch, s, time.Duration(delay*int(time.Millisecond)))
		delay += sc.Delay
	}

	var unhealthyServiceCount int

	// wait for all service checks to finish
	// print results to stdout and count errors
	for i := 0; i < len(sc.ServiceCatalog); i++ {
		o := <-ch
		if o.err != nil {
			unhealthyServiceCount++
			fmt.Fprintf(os.Stderr, "- %s\t%v\n", o.service.Name, o.err)
		} else {
			fmt.Fprintf(os.Stdout, "+ %s\tOK: http %v, %v,\tretries: %v,\t%s\n", o.service.Name, o.gottenStatus, o.requestDuration.Round(time.Millisecond), o.retries, o.service.Url)
		}
	}
	fmt.Println("---")
	if unhealthyServiceCount > 0 {
		fmt.Printf("Unhealthy services: %v\n", unhealthyServiceCount)
	} else {
		fmt.Println("No problems detected.")
	}
	os.Exit(unhealthyServiceCount)
}
