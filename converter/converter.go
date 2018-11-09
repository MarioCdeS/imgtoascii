package converter

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/MarioCdeS/imgtoascii/converter/config"
)

const (
	ramp10 = "@%#*+=-:. "
	ramp70 = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
)

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
	colWidth  int
	rowHeight int
	outRows   int
	outRamp   []rune
}

func Run(cfg *config.Config) ([]string, *Error) {
	img, errLoad := loadImage(cfg.ImagePath)

	if errLoad != nil {
		return nil, &Error{"unable to load image", errLoad}
	}

	intCfg, err := calculateInternalConfig(cfg, img.Bounds())

	if err != nil {
		return nil, err
	}

	return convertToASCII(img, intCfg), nil
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
	imgWidth := imgBounds.Max.X - imgBounds.Min.X

	if cfg.OutCols > imgWidth {
		return nil, &Error{
			fmt.Sprintf("image size is too small for the specified number of output columns (%d)", cfg.OutCols),
			nil,
		}
	}

	imgHeight := imgBounds.Max.Y - imgBounds.Min.Y
	colWidth := imgWidth / cfg.OutCols
	rowHeight := int(float64(colWidth) * cfg.ColRowRatio)
	outRows := imgHeight / rowHeight

	if outRows > imgHeight {
		return nil, &Error{
			fmt.Sprintf("image size is too small for the calculated number of output rows (%d)", outRows),
			nil,
		}
	}

	var outRampStr string

	if cfg.Ramp == config.Ramp10 {
		outRampStr = ramp10
	} else {
		outRampStr = ramp70
	}

	return &internalConfig{
		cfg,
		colWidth,
		rowHeight,
		outRows,
		[]rune(outRampStr),
	}, nil
}

func convertToASCII(img image.Image, cfg *internalConfig) []string {
	numRowsPerStrip := cfg.outRows / cfg.NumCPU
	chs := make([]chan string, cfg.NumCPU)

	for i := 0; i < cfg.NumCPU; i++ {
		var numRows int

		if i == cfg.NumCPU-1 {
			numRows = cfg.outRows - (cfg.NumCPU-1)*numRowsPerStrip
		} else {
			numRows = numRowsPerStrip
		}

		chs[i] = make(chan string, numRows)
		go convertImgStripToASCII(img, i*numRowsPerStrip*cfg.rowHeight, numRows, cfg, chs[i])
	}

	ascii := make([]string, cfg.outRows)
	counts := make([]int, cfg.NumCPU)
	busy := true

	for busy {
		busy = false

		for i := 0; i < cfg.NumCPU; i++ {
			if line, ok := <-chs[i]; ok {
				ascii[i*numRowsPerStrip+counts[i]] = line
				counts[i]++
				busy = true
			}
		}
	}

	return ascii
}

func convertImgStripToASCII(img image.Image, minY int, numRows int, cfg *internalConfig, ch chan<- string) {
	defer close(ch)

	line := make([]rune, cfg.OutCols)
	minX := img.Bounds().Min.X

	for j := 0; j < numRows; j++ {
		for i := 0; i < cfg.OutCols; i++ {
			charMinX := minX + i*cfg.colWidth
			charMinY := minY + j*cfg.rowHeight
			charRect := image.Rect(charMinX, charMinY, charMinX+cfg.colWidth, charMinY+cfg.rowHeight)
			idx := int(math.Floor((rectGrayAverage(img, &charRect) / math.MaxUint8) * float64(len(cfg.outRamp)-1)))
			line[i] = cfg.outRamp[idx]
		}

		ch <- string(line)
	}
}

func rectGrayAverage(img image.Image, rect *image.Rectangle) float64 {
	var total float64

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			total += float64(color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y)
		}
	}

	return total / float64((rect.Max.X-rect.Min.X)*(rect.Max.Y-rect.Min.Y))
}
