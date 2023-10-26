package main

import (
	"fmt"
	"github.com/jonasrdl/hue-cli/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n\n", err)
		os.Exit(1)
	}
}
