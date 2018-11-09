package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	ImagePath   string
	OutCols     int
	ColRowRatio float64
	NumCPU      int
	Ramp
}

func init() {
	flag.Usage = func() {
		output := flag.CommandLine.Output()

		fmt.Fprintf(output, "Usage: %s [<flag>...] <image>\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(output, "Image: path to the image to convert (GIF, JPG, or PNG)")
		fmt.Fprintln(output, "Flags:")
		flag.PrintDefaults()
	}
}

func FromArgs() *Config {
	outCols := flag.Uint("c", 80, "number of output columns")

	maxNumCPU := uint(runtime.NumCPU())
	numCPU := flag.Uint("n", maxNumCPU, "number of CPU cores to use when performing conversion")

	colRowRatio := flag.Float64("r", 2.33, "column-to-row ratio")

	ramp := Ramp10
	flag.Var(&ramp, "g", "grayscale ramp to use (10 or 70, default 10)")

	flag.Parse()

	if *numCPU > maxNumCPU {
		*numCPU = maxNumCPU
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "no input image specified")
		flag.Usage()
		os.Exit(1)
	}

	imagePath := flag.Arg(0)

	return &Config{
		imagePath,
		int(*outCols),
		*colRowRatio,
		int(*numCPU),
		ramp,
	}
}
