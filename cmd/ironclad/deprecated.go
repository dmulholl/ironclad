package main

import (
	"fmt"
	"os"
)

// Exit with an error message and non-zero error code.
func exit(message any) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(1)
}
