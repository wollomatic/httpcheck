package main

import (
	"fmt"
	"os"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}
