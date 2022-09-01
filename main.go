package main

import "fmt"

// Version is set at build time
var Version string

func main() {
	fmt.Printf("Running version %s\n", Version)
}
