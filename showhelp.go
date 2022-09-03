package main

import (
	"fmt"
	"os"
)

func showHelp() {
	fmt.Printf("httpcheck version %s (github.com/wollomatic/httpcheck)\n", VERSION)
	fmt.Println(`
Tests web service availability as specified in given yaml file.

Usage: httpcheck file.yaml

Return codes:
0 - all services are ok
1 - 1 service is not ok or no filename is given or problem in yaml file
n - n services are not ok`)
	os.Exit(0)
}
