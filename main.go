package main

import (
	"fmt"
	"os"

	"github.com/MarioCdeS/imgtoascii/converter"
	"github.com/MarioCdeS/imgtoascii/converter/config"
)

func main() {
	cfg := config.FromArgs()
	ascii, err := converter.Run(cfg)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)

		if err.Cause != nil {
			fmt.Fprintln(os.Stderr, "Cause:", err.Cause)
		}

		os.Exit(3)
	}

	for _, line := range ascii {
		fmt.Println(line)
	}
}
