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
	flag.Usage = func() {
		output := flag.CommandLine.Output()

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
	outPath := flag.String("o", "out.txt", "path to output text file")

	ramp := Ramp10
	flag.Var(&ramp, "g", "greyscale ramp to use (10 or 70, default 10)")

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
