package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ImagePath   string
	OutPath     string
	OutCols     int
	ColRowRatio float64
	Ramp
}

func init() {
	flag.Usage = func() {
		output := flag.CommandLine.Output()

		fmt.Fprintf(output, "Usage: %s [flags] <image>\n", os.Args[0])
		fmt.Fprintln(output, "Image: path to image to convert (GIF, JPG, or PNG)")
		fmt.Fprintln(output, "Flags:")
		flag.PrintDefaults()
	}
}

func FromArgs() *Config {
	outPath := flag.String("o", "out.txt", "path to output text file")

	outCols := flag.Int("c", 80, "number of output columns")

	colRowRatio := flag.Float64("r", 2.33, "column-to-row ratio")

	ramp := Ramp10
	flag.Var(&ramp, "g", "grayscale ramp to use (10 or 70, default 10)")

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
		*outCols,
		*colRowRatio,
		ramp,
	}
}
