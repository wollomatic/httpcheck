package main

import (
	"fmt"
	"os"
	"time"
)

const VERSION = "0.1.0"

func main() {

	// if no filename argument is given show help and exit
	if len(os.Args) != 2 {
		showHelp()
	}

	// read yaml file
	sc := readYamlFile(os.Args[1])

	fmt.Printf("[%s] starting service checks\n", time.Now().Format("2006-01-02 15:04:05"))

	// create channel for service responses
	ch := make(chan serviceResponse)

	// start concurrent service checks
	delay := 0
	for _, s := range sc.Service {
		go checkService(ch, s, time.Duration(delay*int(time.Millisecond)))
		delay += sc.Delay
	}

	var unhealthyServiceCount int

	// wait for all service checks to finish
	// print results to stdout and count errors
	for i := 0; i < len(sc.Service); i++ {
		o := <-ch
		if o.err != nil {
			unhealthyServiceCount++
			fmt.Printf("- %-30s   %v\n", o.service.Name, o.err)
		} else {
			fmt.Printf("+ %-30s   %-10s %-30v %10v   %3v retries   %s\n", o.service.Name, o.response.Proto, o.response.Status, o.requestDuration.Round(time.Millisecond), o.retries, o.response.Header.Get("Server"))
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
