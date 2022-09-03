package main

import (
	"fmt"
	"os"
)

func showHelp() {
	fmt.Println(`Usage: httpcheck file.yaml
Commands:
check filename.yaml   check services
sample                show sample yaml file

Return codes:
0 - all services are ok
1 - 1 service is not ok or no filename is given or problem in yaml file
n - n services are not ok`)
	os.Exit(0)
}

