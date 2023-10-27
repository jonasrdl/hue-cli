package main

import (
	"fmt"
	"os"

	"github.com/jonasrdl/hue-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n\n", err)
		os.Exit(1)
	}
}
