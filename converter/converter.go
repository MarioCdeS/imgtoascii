package converter

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/MarioCdeS/imgtoascii/converter/config"
)

const (
	ramp10      = " .:-=+*#%@"
	ramp70      = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	redWeight   = 0.299
	greenWeight = 0.587
	blueWeight  = 0.114
)

var ramp10Runes []rune
var ramp70Runes []rune

type Error struct {
	Msg   string
	Cause error
}

func (e *Error) Error() string {
	if e == nil {
		return "unknown error"
	}

	return e.Msg
}

type internalConfig struct {
	*config.Config
	imgWidth       int
	imgHeight      int
	colWidth       int
	rowHeight      int
	outRows        int
	numCPU         int
	numColsPerCPU  int
	numColsLastCPU int
}

func init() {
	ramp10Runes = []rune(ramp10)
	ramp70Runes = []rune(ramp70)
}

func Run(cfg *config.Config) *Error {
	img, errLoad := loadImage(cfg.ImagePath)

	if errLoad != nil {
		return &Error{"unable to load image", errLoad}
	}

	intCfg, err := calculateInternalConfig(cfg, img.Bounds())

	if err != nil {
		return err
	}

	convertToASCII(img, intCfg)
	return nil
}

func loadImage(path string) (image.Image, error) {
	reader, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer reader.Close()

	img, _, err := image.Decode(reader)
	return img, err
}

func calculateInternalConfig(cfg *config.Config, imgBounds image.Rectangle) (*internalConfig, *Error) {
	imgWidth := imgBounds.Max.X

	if cfg.OutCols > imgWidth {
		return nil, &Error{
			fmt.Sprintf("image size is too small for the specified number of output columns (%d)", cfg.OutCols),
			nil,
		}
	}

	imgHeight := imgBounds.Max.Y
	colWidth := imgWidth / cfg.OutCols
	rowHeight := int(float64(colWidth) * cfg.ColRowRatio)
	outRows := imgHeight / rowHeight

	if outRows > imgHeight {
		return nil, &Error{
			fmt.Sprintf("image size is too small for the calculated number of output rows (%d)", outRows),
			nil,
		}
	}

	numCPU := 1 // runtime.NumCPU()
	numColsPerCPU := cfg.OutCols / numCPU
	numColsLastCPU := cfg.OutCols - (numCPU-1)*numColsPerCPU

	return &internalConfig{
		cfg,
		imgWidth,
		imgHeight,
		colWidth,
		rowHeight,
		outRows,
		numCPU,
		numColsPerCPU,
		numColsLastCPU,
	}, nil
}

func convertToASCII(img image.Image, cfg *internalConfig) {
	var ramp []rune

	if cfg.Ramp == config.Ramp10 {
		ramp = ramp10Runes
	} else {
		ramp = ramp70Runes
	}

	segmentWidth := cfg.numColsPerCPU * cfg.colWidth
	outChannels := make([]chan rune, cfg.numCPU)

	var bounds image.Rectangle

	for i := 0; i < cfg.numCPU-1; i++ {
		bounds = image.Rect(i*segmentWidth, 0, (i+1)*segmentWidth-1, cfg.imgHeight)
		outChannels[i] = make(chan rune, cfg.numColsPerCPU)
		go convertSegmentToASCII(img, bounds, cfg.colWidth, cfg.rowHeight, ramp, outChannels[i])
	}

	bounds = image.Rect((cfg.numCPU-1)*segmentWidth, 0, cfg.imgWidth, cfg.imgHeight)
	outChannels[cfg.numCPU-1] = make(chan rune, cfg.numColsLastCPU)
	go convertSegmentToASCII(img, bounds, cfg.colWidth, cfg.rowHeight, ramp, outChannels[cfg.numCPU-1])

	outputLine := make([]rune, cfg.OutCols)

	for i := 0; i < cfg.outRows; i++ {
		for j, ch := range outChannels {
			var stop int

			if j == cfg.numCPU-1 {
				stop = cfg.numColsLastCPU
			} else {
				stop = cfg.numColsPerCPU
			}

			for k := 0; k < stop; k++ {
				outputLine[j*cfg.numColsPerCPU+k] = <-ch
			}
		}

		fmt.Println(string(outputLine))
	}
}

func convertSegmentToASCII(img image.Image, bounds image.Rectangle, colWidth, rowHeight int, ramp []rune, ch chan rune) {
	numPixPerSeg := uint32(colWidth * rowHeight)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += rowHeight {
		for x := bounds.Min.X; x < bounds.Max.X; x += colWidth {
			var totR, totG, totB uint32 = 0, 0, 0

			for i := 0; i < colWidth; i++ {
				for j := 0; j < rowHeight; j++ {
					r, g, b, _ := img.At(x+i, y+j).RGBA()
					totR += r
					totG += g
					totB += b
				}
			}

			lum := (redWeight*float64(totR/numPixPerSeg) + greenWeight*float64(totG/numPixPerSeg) +
				blueWeight*float64(totB/numPixPerSeg)) / 65535.0
			i := int(math.Round(lum * float64(len(ramp))))

			ch <- ramp[i]
		}
	}
}
