package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const VERSION = "0.3.0"

func main() {

	// if no filename argument is given show help and exit
	if len(os.Args) != 2 {
		showHelp()
	}

	// read yaml file
	sc := readYamlFile(os.Args[1])

	log.Println("httpcheck starting service checks")

	// create channel for service responses
	ch := make(chan serviceResponse)

	// start concurrent service checks
	delay := 0
	for _, s := range sc.Services {
		go checkService(ch, s, time.Duration(delay*int(time.Millisecond)))
		delay += sc.Delay
	}

	var unhealthyServiceCount int

	// wait for all service checks to finish
	// print results to stdout and count errors
	fmt.Println("Result   Service name                     Method Protocol   Response                    Duration     # retries   Server            Search text")
	for i := 0; i < len(sc.Services); i++ {
		o := <-ch
		if o.err != nil {
			unhealthyServiceCount++
			fmt.Printf("Problem: %-30s   %v\n", o.service.Name, o.err)
		} else {
			fmt.Printf("OK       %-30s   %-4s   %-10s %-25s %10v   %3v retries   %-15s   %s\n", o.service.Name, o.service.Method, o.response.Proto, o.response.Status, o.requestDuration.Round(time.Millisecond), o.retries, o.response.Header.Get("Server"), o.service.SearchText)
		}
	}
	fmt.Println("---")
	s := "s"
	if len(sc.Services) == 1 {
		s = ""
	}
	fmt.Printf("%v service%s checked. ", len(sc.Services), s)
	if unhealthyServiceCount > 0 {
		if unhealthyServiceCount == 1 {
			s = ""
		}
		fmt.Printf("Unhealthy service%s: %v\n", s, unhealthyServiceCount)
	} else {
		fmt.Println("No problems detected.")
	}
	os.Exit(unhealthyServiceCount)
}
