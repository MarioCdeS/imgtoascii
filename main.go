package main

import (
	"fmt"
	"os"

	"github.com/MarioCdeS/imgtoascii/converter"
	"github.com/MarioCdeS/imgtoascii/converter/config"
)

func main() {
	cfg := config.FromArgs()

	if err := converter.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)

		if err.Cause != nil {
			fmt.Fprintln(os.Stderr, "Cause:", err.Cause)
		}

		os.Exit(3)
	}
}
