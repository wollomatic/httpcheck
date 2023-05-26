package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wollomatic/httpcheck/pkg/httpcheck"
)

const VERSION = "0.4.0"

func main() {

	// if no filename argument is given show help and exit
	if len(os.Args) != 2 {
		showHelp()
	}

	// read yaml file
	sc, err := httpcheck.ReadYamlFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	log.Println("httpcheck starting service checks")

	// create channel for service responses
	ch := make(chan httpcheck.Response)

	// start concurrent service checks
	go func() {
		for _, s := range sc.Services {
			go checkService(ch, s)
			time.Sleep(time.Duration(sc.Delay) * time.Millisecond)
		}
	}()

	var unhealthyServiceCount int

	// wait for all service checks to finish
	// print results to stdout and count errors
	fmt.Println("Result   Service name                     Method Protocol   Response                    Duration     # retries   Server            Search text")
	for i := 0; i < len(sc.Services); i++ {
		o := <-ch
		if o.Err != nil {
			unhealthyServiceCount++
			fmt.Printf("Problem: %-30s   %v\n", o.Service.Name, o.Err)
		} else {
			fmt.Printf("OK       %-30s   %-4s   %-10s %-25s %10v   %3v retries   %-15s   %s\n", o.Service.Name, o.Service.Method, o.Response.Proto, o.Response.Status, o.RequestDuration.Round(time.Millisecond), o.Retries, o.Response.Header.Get("Server"), o.Service.SearchText)
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
