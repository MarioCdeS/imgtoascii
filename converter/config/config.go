package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ImagePath string
	OutPath   string
	Ramp
}

func init() {
	output := flag.CommandLine.Output()

	flag.Usage = func() {
		fmt.Fprintf(output, "Usage: %s [flags] <image>\n", os.Args[0])
		fmt.Fprintln(
			output,
			"Image: path to image to convert (GIF, JPG, or PNG)",
		)
		fmt.Fprintln(output, "Flags:")
		flag.PrintDefaults()
	}
}

func FromArgs() *Config {
	outPath := flag.String(
		"o", "out.txt", "path to output text file (defaults to out.txt)",
	)

	ramp := Ramp10
	flag.Var(&ramp, "g", "greyscale ramp to use (10 or 70, defaults to 10)")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "no input image specified")
		flag.Usage()
		os.Exit(1)
	}

	imagePath := flag.Arg(0)

	return &Config{
		imagePath,
		*outPath,
		ramp,
	}
}
